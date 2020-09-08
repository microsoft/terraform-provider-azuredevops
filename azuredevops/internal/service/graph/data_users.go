package graph

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataUsers schema and implementation for users data source
func DataUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataUsersRead,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"origin", "origin_id"},
			},
			"subject_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"origin": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"principal_name"},
			},
			"origin_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{"principal_name"},
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getUserHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"descriptor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mail_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataUsersRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	users := make([]interface{}, 0)
	subjectTypes := []string{}

	linq.From(d.Get("subject_types").(*schema.Set).List()).
		SelectT(func(x interface{}) string {
			return x.(string)
		}).
		ToSlice(&subjectTypes)
	principalName := d.Get("principal_name").(string)
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)

	var currentToken string
	for hasMore := true; hasMore; {
		newUsers, latestToken, err := getUsersWithContinuationToken(clients, &subjectTypes, currentToken)
		currentToken = latestToken
		hasMore = currentToken != ""
		if err != nil {
			return err
		}

		linq.From(newUsers).
			WhereT(func(x interface{}) bool {
				usr := x.(graph.GraphUser)
				b := true
				if principalName != "" {
					b = usr.PrincipalName != nil && strings.EqualFold(*usr.PrincipalName, principalName)
				}
				if b && origin != "" {
					b = usr.Origin != nil && strings.EqualFold(*usr.Origin, origin)
				}
				if b && originID != "" {
					b = usr.OriginId != nil && strings.EqualFold(*usr.OriginId, originID)
				}
				return b
			}).
			ToSlice(&newUsers)
		fusers, err := flattenUsers(&newUsers)
		if err != nil {
			return err
		}
		users = append(users, fusers...)
	}

	var descriptors []string
	linq.From(users).
		SelectT(
			func(x interface{}) string {
				item := x.(map[string]interface{})
				return item["descriptor"].(string)
			},
		).
		ToSlice(&descriptors)

	h := sha1.New()
	if _, err := h.Write([]byte(strings.Join(descriptors, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for user descriptors: %v", err)
	}
	d.SetId("users#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	if err := d.Set("users", users); err != nil {
		return fmt.Errorf("Error setting `users`: %+v", err)
	}

	return nil
}

func getUserHash(v interface{}) int {
	return hashcode.String(v.(map[string]interface{})["descriptor"].(string))
}

func flattenUsers(input *[]graph.GraphUser) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	results := make([]interface{}, len(*input))
	for i, element := range *input {
		output, err := flattenUser(&element)
		if err != nil {
			return nil, err
		}
		results[i] = output
	}
	return results, nil
}

func flattenUser(user *graph.GraphUser) (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if v := user.Descriptor; v != nil {
		s["descriptor"] = *v
	}
	if v := user.PrincipalName; v != nil {
		s["principal_name"] = *v
	}
	if v := user.Origin; v != nil {
		s["origin"] = *v
	}
	if v := user.OriginId; v != nil {
		s["origin_id"] = *v
	}
	if v := user.DisplayName; v != nil {
		s["display_name"] = *v
	}
	if v := user.MailAddress; v != nil {
		s["mail_address"] = *v
	}

	return s, nil
}

func getUsersWithContinuationToken(clients *client.AggregatedClient, subjectTypes *[]string, continuationToken string) ([]graph.GraphUser, string, error) {
	args := graph.ListUsersArgs{
		SubjectTypes: subjectTypes,
	}
	if continuationToken != "" {
		args.ContinuationToken = &continuationToken
	}
	response, err := clients.GraphClient.ListUsers(clients.Ctx, args)
	if err != nil {
		return nil, "", fmt.Errorf("Error listing users: %q", err)
	}

	continuationToken = ""
	if response.ContinuationToken != nil && (*response.ContinuationToken)[0] != "" {
		continuationToken = (*response.ContinuationToken)[0]
	}

	return *response.GraphUsers, continuationToken, nil
}
