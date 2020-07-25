package core

import (
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"administrators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Computed:   true,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
		},
	}
}

func resourceTeamCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamName := d.Get("name").(string)

	team, err := clients.CoreClient.CreateTeam(clients.Ctx, core.CreateTeamArgs{
		ProjectId: &projectID,
		Team: &core.WebApiTeam{
			Name: &teamName,
		},
	})

	if err != nil {
		return err
	}

	if v, ok := d.GetOk("administrators"); ok {
		administrators := tfhelper.ExpandStringSet(v.(*schema.Set))
		err := updateTeamAdministrators(clients, team, &administrators)
		if err != nil {
			ierr := clients.CoreClient.DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
				ProjectId: converter.String(team.ProjectId.String()),
				TeamId:    converter.String(team.Id.String()),
			})
			if ierr != nil {
				log.Printf("[ERROR] Failed to delete project after update of administrators failed %+v", ierr)
			}
			return err
		}
	}

	if v, ok := d.GetOk("members"); ok {
		members := tfhelper.ExpandStringSet(v.(*schema.Set))
		err := updateTeamMembers(clients, team, &members)
		if err != nil {
			ierr := clients.CoreClient.DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
				ProjectId: converter.String(team.ProjectId.String()),
				TeamId:    converter.String(team.Id.String()),
			})
			if ierr != nil {
				log.Printf("[ERROR] Failed to delete project after update of members failed %+v", ierr)
			}
			return err
		}
	}

	d.SetId(team.Id.String())
	return resourceTeamRead(d, m)
}

func resourceTeamRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamName := d.Get("name").(string)

	team, members, administrators, err := readTeam(clients, projectID, teamName)
	if err != nil {
		return err
	}

	flattenTeam(d, team, members, administrators)
	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
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

func readTeam(clients *client.AggregatedClient, projectID string, teamName string) (*core.WebApiTeam, *schema.Set, *schema.Set, error) {
	teamList, err := clients.CoreClient.GetTeams(clients.Ctx, core.GetTeamsArgs{
		ProjectId:      converter.String(projectID),
		Mine:           converter.Bool(false),
		ExpandIdentity: converter.Bool(true), // required for readTeamMembers
	})

	if err != nil {
		return nil, nil, nil, err
	}

	if teamList == nil || len(*teamList) <= 0 {
		return nil, nil, nil, fmt.Errorf("Project [%s] does not contain any teams", projectID)
	}

	iTeam := linq.From(*teamList).
		FirstWith(func(v interface{}) bool {
			item := v.(core.WebApiTeam)
			return strings.EqualFold(*item.Name, teamName)
		})
	if iTeam == nil {
		return nil, nil, nil, fmt.Errorf("Unable to find Team with name [%s] in project with ID [%s]", teamName, projectID)
	}

	team := iTeam.(core.WebApiTeam)
	members, err := readTeamMembers(clients, &team)
	if err != nil {
		return nil, nil, nil, err
	}
	administrators, err := readTeamAdministrators(clients, &team)
	if err != nil {
		return nil, nil, nil, err
	}

	return &team, members, administrators, nil
}

func flattenTeam(d *schema.ResourceData, team *core.WebApiTeam, members *schema.Set, administrators *schema.Set) {
	if team == nil {
		d.SetId("")
		return
	}
	d.SetId(team.Id.String())
	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("administrators", administrators)
	d.Set("members", members)
}

func readTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam) (*schema.Set, error) {
	members, err := clients.IdentityClient.ReadMembers(clients.Ctx, identity.ReadMembersArgs{
		ContainerId: converter.String(team.Id.String()),
	})
	if err != nil {
		return nil, err
	}

	return readIdentities(clients, members)
}

func updateTeamMembers(clients *client.AggregatedClient, team *core.WebApiTeam, subjectDescriptors *[]string) error {
	return fmt.Errorf("Not implemented")
}

func readTeamAdministrators(clients *client.AggregatedClient, team *core.WebApiTeam) (*schema.Set, error) {
	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Identity,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return nil, err
	}

	token := fmt.Sprintf("%s\\%s", team.ProjectId.String(), team.Id.String())
	acl, err := sn.GetAccessControlList(&token, nil)
	if err != nil {
		return nil, err
	}

	adminDescriptorList := []string{}
	for _, ace := range *acl.AcesDictionary {
		if *ace.Allow&15 > 0 {
			adminDescriptorList = append(adminDescriptorList, *ace.Descriptor)
		}
	}

	return readIdentities(clients, &adminDescriptorList)
}

func updateTeamAdministrators(clients *client.AggregatedClient, team *core.WebApiTeam, subjectDescriptors *[]string) error {
	return fmt.Errorf("Not implemented")
}

func readIdentities(clients *client.AggregatedClient, members *[]string) (*schema.Set, error) {
	set := schema.NewSet(schema.HashString, nil)

	if members == nil || len(*members) <= 0 {
		return set, nil
	}

	descriptors := linq.From(*members).
		Aggregate(func(r interface{}, i interface{}) interface{} {
			if r.(string) == "" {
				return i
			}
			return r.(string) + "," + i.(string)
		}).(string)

	identities, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		Descriptors: &descriptors,
	})

	if err != nil {
		return nil, err
	}

	if identities != nil && len(*identities) > 0 {
		for _, identity := range *identities {
			set.Add(*identity.SubjectDescriptor)
		}
	}

	return set, nil
}
