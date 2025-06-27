package core

import (
	"fmt"
	"log"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"administrators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTeamCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamData := core.WebApiTeam{
		Name: converter.ToPtr(d.Get("name").(string)),
	}
	if description, ok := d.GetOk("description"); ok {
		teamData.Description = converter.String(description.(string))
	}

	team, err := clients.CoreClient.CreateTeam(clients.Ctx, core.CreateTeamArgs{
		ProjectId: &projectID,
		Team:      &teamData,
	})
	if err != nil {
		return fmt.Errorf("Creating Team: %+v", err)
	}

	teamID := team.Id.String()
	var administratorSet *schema.Set
	if v, ok := d.GetOk("administrators"); ok {
		administratorSet = v.(*schema.Set)
		administrators := tfhelper.ExpandStringSet(administratorSet)
		if err = updateTeamAdministrators(d, clients, team, &administrators); err != nil {
			if delErr := clients.CoreClient.DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
				ProjectId: converter.String(team.ProjectId.String()),
				TeamId:    converter.String(team.Id.String()),
			}); delErr != nil {
				log.Printf("[ERROR] Failed to delete project after update of administrators %+v", delErr)
			}
			return err
		}
	}

	var memberSet *schema.Set
	if v, ok := d.GetOk("members"); ok {
		memberSet = v.(*schema.Set)
		members := tfhelper.ExpandStringSet(memberSet)
		if err = setTeamMembers(clients, team, &members); err != nil {
			if delErr := clients.CoreClient.DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
				ProjectId: converter.String(team.ProjectId.String()),
				TeamId:    converter.String(team.Id.String()),
			}); delErr != nil {
				log.Printf("[ERROR] Failed to delete project after update of members %+v", delErr)
			}
			return err
		}
	}

	if err = waitForTeamStateChange(d, clients, projectID, teamID, teamData.Name, teamData.Description, memberSet, administratorSet); err != nil {
		return err
	}

	d.SetId(team.Id.String())
	return resourceTeamRead(d, m)
}

func resourceTeamRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Id()
	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(false),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if team == nil {
		d.SetId("")
		log.Printf(" team not found. Project ID : %s, Team ID: %s", projectID, teamID)
		return nil
	}

	members, err := getTeamMembers(clients, team)
	if err != nil {
		return err
	}

	administrators, err := getTeamAdministrators(d, clients, team)
	if err != nil {
		return err
	}

	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("administrators", administrators)
	d.Set("members", members)

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
		StorageKey: team.Id,
	})
	if err != nil {
		return fmt.Errorf("get team descriptor. Error: %+v", err)
	}
	d.Set("descriptor", descriptor.Value)

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var team *core.WebApiTeam
	var err error

	projectID := d.Get("project_id").(string)
	teamID := d.Id()

	var newTeamName *string
	var newDescription *string

	if d.HasChange("name") || d.HasChange("description") {
		teamData := core.WebApiTeam{}

		if d.HasChange("name") {
			teamName := d.Get("name").(string)
			newTeamName = &teamName
			teamData.Name = &teamName
		}

		if d.HasChange("description") {
			description := d.Get("description").(string)
			newDescription = &description
			teamData.Description = &description
		}

		team, err = clients.CoreClient.UpdateTeam(clients.Ctx, core.UpdateTeamArgs{
			ProjectId: &projectID,
			TeamId:    &teamID,
			TeamData:  &teamData,
		})
		if err != nil {
			return err
		}
	} else {
		team, err = clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId:      converter.String(projectID),
			TeamId:         converter.String(teamID),
			ExpandIdentity: converter.Bool(false),
		})
		if err != nil {
			return err
		}
	}

	var administratorSet *schema.Set
	if d.HasChange("administrators") {
		log.Printf("Updating list of administrators for team %s", *team.Name)

		administratorSet = d.Get("administrators").(*schema.Set)
		administrators := tfhelper.ExpandStringSet(administratorSet)
		err = updateTeamAdministrators(d, clients, team, &administrators)
		if err != nil {
			return err
		}
	}

	var memberSet *schema.Set
	if d.HasChange("members") {
		log.Printf("Updating list of members for team %s", *team.Name)

		memberSet = d.Get("members").(*schema.Set)
		members := tfhelper.ExpandStringSet(memberSet)
		err = setTeamMembers(clients, team, &members)
		if err != nil {
			return err
		}
	}

	if err := waitForTeamStateChange(d, clients, projectID, teamID, newTeamName, newDescription, memberSet, administratorSet); err != nil {
		return err
	}

	return resourceTeamRead(d, m)
}

func resourceTeamDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Id()

	err := clients.CoreClient.DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
		ProjectId: &projectID,
		TeamId:    &teamID,
	})
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func waitForTeamStateChange(d *schema.ResourceData, clients *client.AggregatedClient, projectID string, teamID string, name *string, description *string, memberSet *schema.Set, administratorSet *schema.Set) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"

			team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
				ProjectId:      converter.String(projectID),
				TeamId:         converter.String(teamID),
				ExpandIdentity: converter.Bool(false),
			})
			if err != nil {
				return nil, "", fmt.Errorf("Reading team data: %+v", err)
			}

			bDescriptionUpdated := nil == description || *team.Description == *description
			bNameUpdated := nil == name || *team.Name == *name

			bAdministratorsUpdated := true
			if administratorSet != nil {
				actualAdministrators, err := getTeamAdministrators(d, clients, team)
				if err != nil {
					return nil, "", fmt.Errorf("Reading team administrators: %+v", err)
				}
				bAdministratorsUpdated = actualAdministrators.Len() == administratorSet.Len()
			}

			dashboards, err := clients.DashboardClient.GetDashboardsByProject(clients.Ctx, dashboard.GetDashboardsByProjectArgs{
				Project: converter.String(projectID),
				Team:    converter.String(teamID),
			})
			if err != nil {
				return nil, "", fmt.Errorf("Reading Team dashboard: %+v", err)
			}
			dashboardUpdate := true
			if dashboards == nil && len(*dashboards) == 0 {
				dashboardUpdate = false
			}

			bMembersUpdated := true
			if memberSet != nil {
				actualMemberships, err := getTeamMembers(clients, team)
				if err != nil {
					return nil, "", fmt.Errorf("Reading team memberships: %+v", err)
				}
				bMembersUpdated = actualMemberships.Len() == memberSet.Len()
			}

			if bNameUpdated && bDescriptionUpdated && bAdministratorsUpdated && bMembersUpdated && dashboardUpdate {
				state = "Synched"
			}
			return state, state, nil
		},
		Timeout:                   30 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	if _, err := stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("waiting for state change for team %s in project %s. %v ", teamID, projectID, err)
	}

	return nil
}

func getTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam) (*schema.Set, error) {
	members, err := clients.IdentityClient.ReadMembers(clients.Ctx, identity.ReadMembersArgs{
		ContainerId: converter.String(team.Id.String()),
	})
	if err != nil {
		return nil, err
	}

	return getSubjectDescriptors(clients, members)
}

func setTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam, subjectDescriptors *[]string) error {
	var err error

	currentMemberSet, err := getTeamMembers(clients, team)
	if err != nil {
		return err
	}
	if (subjectDescriptors == nil || len(*subjectDescriptors) == 0) && currentMemberSet.Len() == 0 {
		return nil
	}
	if subjectDescriptors == nil {
		subjectDescriptors = &[]string{}
	}

	currentMembers := currentMemberSet.List()

	// determine the list of all removed members
	err = removeTeamMembers(clients, team, linq.From(currentMembers).Except(linq.From(*subjectDescriptors)))
	if err != nil {
		return err
	}

	// determine the list of all added members
	err = addTeamMembers(clients, team, linq.From(*subjectDescriptors).Except(linq.From(currentMembers)), false)
	if err != nil {
		return err
	}

	return nil
}

func getIdentitiesFromSubjects(clients *client.AggregatedClient, query linq.Query) (*[]identity.Identity, error) {
	if !query.Any() {
		return &[]identity.Identity{}, nil
	}

	discriptors := query.
		Aggregate(func(r interface{}, i interface{}) interface{} {
			if r.(string) == "" {
				return i
			}
			return r.(string) + "," + i.(string)
		}).(string)

	idlist, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SubjectDescriptors: converter.String(discriptors),
	})
	if err != nil {
		return nil, err
	}

	return idlist, err
}

func removeTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam, query linq.Query) error {
	idList, err := getIdentitiesFromSubjects(clients, query)
	if err != nil {
		return err
	}

	for _, id := range *idList {
		log.Printf("[TRACE] Removing member %s from team %s", id.Id.String(), *team.Name)

		_, err := clients.IdentityClient.RemoveMember(clients.Ctx, identity.RemoveMemberArgs{
			ContainerId: converter.String(team.Id.String()),
			MemberId:    converter.String(id.Id.String()),
		})
		if err != nil {
			return fmt.Errorf("Error removing member %s from team %s: %+v", id.Id.String(), *team.Name, err)
		}
	}
	return nil
}

func addTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam, query linq.Query, isAddMode bool) error {
	idList, err := getIdentitiesFromSubjects(clients, query)
	if err != nil {
		return err
	}
	if idList == nil || len(*idList) != query.Count() {
		return fmt.Errorf("Failed to load identity data for subjects")
	}

	for _, id := range *idList {
		log.Printf("[TRACE] Adding member %s to team %s", id.Id.String(), *team.Name)

		ok, err := clients.IdentityClient.AddMember(clients.Ctx, identity.AddMemberArgs{
			ContainerId: converter.String(team.Id.String()),
			MemberId:    converter.String(id.Id.String()),
		})
		if err != nil {
			return fmt.Errorf("Error adding member %s to team %s: %+v", *id.SubjectDescriptor, *team.Name, err)
		}
		if ok != nil && !*ok {
			if !isAddMode {
				return fmt.Errorf("Failed adding member %s to team %s", *id.SubjectDescriptor, *team.Name)
			} else {
				log.Printf("[TRACE] Member %s is already a member of team %s", *id.SubjectDescriptor, *team.Name)
			}
		}
	}

	return nil
}

func getIdentitySecurityNamespace(d *schema.ResourceData, clients *client.AggregatedClient, team *core.WebApiTeam) (*securityhelper.SecurityNamespace, error) {
	return securityhelper.NewSecurityNamespace(d,
		clients,
		securityhelper.SecurityNamespaceIDValues.Identity,
		func(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
			return team.ProjectId.String() + "\\" + team.Id.String(), nil
		})
}

// getTeamAdministrators returns the current list of team administrators as a set of SubjectDescriptors
func getTeamAdministrators(d *schema.ResourceData, clients *client.AggregatedClient, team *core.WebApiTeam) (*schema.Set, error) {
	sn, err := getIdentitySecurityNamespace(d, clients, team)
	if err != nil {
		return nil, err
	}

	actionDefinitions, err := sn.GetActionDefinitions()
	if err != nil {
		return nil, err
	}

	acl, err := sn.GetAccessControlList(nil)
	if err != nil {
		return nil, err
	}

	adminDescriptorList := []string{}
	if acl != nil && acl.AcesDictionary != nil {
		bit := *(*actionDefinitions)["Read"].Bit | *(*actionDefinitions)["Write"].Bit | *(*actionDefinitions)["Delete"].Bit | *(*actionDefinitions)["ManageMembership"].Bit | *(*actionDefinitions)["CreateScope"].Bit
		for _, ace := range *acl.AcesDictionary {
			if *ace.Allow&bit == bit {
				adminDescriptorList = append(adminDescriptorList, *ace.Descriptor)
			}
		}
	}
	return getSubjectDescriptors(clients, &adminDescriptorList)
}

func updateTeamAdministrators(d *schema.ResourceData, clients *client.AggregatedClient, team *core.WebApiTeam, subjectDescriptors *[]string) error {
	currentAdministratorSet, err := getTeamAdministrators(d, clients, team)
	if err != nil {
		return err
	}
	if (subjectDescriptors == nil || len(*subjectDescriptors) == 0) && currentAdministratorSet.Len() == 0 {
		return nil
	}

	currentAdministrators := currentAdministratorSet.List()

	log.Print("[DEBUG] updateTeamAdministrators::removing deleted administrators from team")
	err = setTeamAdministratorsPermissions(d,
		clients,
		team,
		// determine the list of all removed administrators
		linq.From(currentAdministrators).Except(linq.From(*subjectDescriptors)),
		securityhelper.PermissionTypeValues.NotSet)
	if err != nil {
		return err
	}

	log.Print("[DEBUG] updateTeamAdministrators::adding missing administrators to team")
	err = setTeamAdministratorsPermissions(d,
		clients,
		team,
		// determine the list of all added administrators
		linq.From(*subjectDescriptors).Except(linq.From(currentAdministrators)),
		securityhelper.PermissionTypeValues.Allow)
	if err != nil {
		return err
	}

	return nil
}

func setTeamAdministratorsPermissions(d *schema.ResourceData, clients *client.AggregatedClient, team *core.WebApiTeam, subjectDescriptors linq.Query, permission securityhelper.PermissionType) error {
	if !subjectDescriptors.Any() {
		log.Print("[DEBUG] setTeamAdministratorsPermissions::list of subject descriptors is empty")
		return nil
	}

	sn, err := getIdentitySecurityNamespace(d, clients, team)
	if err != nil {
		return err
	}

	principalPermissionCreator := func(query linq.Query, permission securityhelper.PermissionType) *[]securityhelper.SetPrincipalPermission {
		var subjectList []securityhelper.SetPrincipalPermission

		query.Select(func(item interface{}) interface{} {
			// item: SubjectDescriptor (string)
			return securityhelper.SetPrincipalPermission{
				Replace: true,
				PrincipalPermission: securityhelper.PrincipalPermission{
					SubjectDescriptor: item.(string),
					Permissions: map[securityhelper.ActionName]securityhelper.PermissionType{
						"Read":             permission,
						"Write":            permission,
						"Delete":           permission,
						"ManageMembership": permission,
						"CreateScope":      permission,
					},
				},
			}
		}).ToSlice(&subjectList)
		return &subjectList
	}

	principalPermissions := principalPermissionCreator(subjectDescriptors, permission)
	err = sn.SetPrincipalPermissions(principalPermissions)
	if err != nil {
		return err
	}

	return nil
}

// readIdentities returns the SubjectDescriptor for every identity passed
func getSubjectDescriptors(clients *client.AggregatedClient, members *[]string) (*schema.Set, error) {
	set := schema.NewSet(schema.HashString, nil)

	if members == nil || len(*members) == 0 {
		return set, nil
	}

	start := 0
	step := 20 // 20 descriptors per request
	var subMembers []string
	for start*step <= len(*members) {
		if (start+1)*step < len(*members) {
			subMembers = (*members)[start*step : (start+1)*step]
		} else {
			subMembers = (*members)[start*step:]
		}
		start++
		if len(subMembers) > 0 {
			descriptors := linq.From(subMembers).
				Aggregate(func(r interface{}, i interface{}) interface{} {
					if r.(string) == "" {
						return i
					}
					return r.(string) + "," + i.(string)
				}).(string)

			memberIdentities, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
				Descriptors: &descriptors,
			})
			if err != nil {
				return nil, err
			}

			if memberIdentities != nil && len(*memberIdentities) > 0 {
				for _, memberIdentity := range *memberIdentities {
					set.Add(*memberIdentity.SubjectDescriptor)
				}
			}
		}
	}

	return set, nil
}
