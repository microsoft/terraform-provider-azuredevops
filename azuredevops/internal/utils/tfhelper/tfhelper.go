package tfhelper

import (
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/secretmemo"
)

func HashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func calcSecretHashKey(secretKey string) string {
	return secretKey + "_hash"
}

// DiffFuncSuppressSecretChanged is used to suppress unneeded `apply` updates to a resource.
//
// It returns `true` when `new` appears to be the same value
// as a previously stored and bcrypt'd value stored in state during a previous `apply`.
// Relies on flatten/expand logic to help store that hash. See FlattenSecret, below.*/
func DiffFuncSuppressSecretChanged(k, old, new string, d *schema.ResourceData) bool {
	memoKey := calcSecretHashKey(k)
	memoValue := d.Get(memoKey).(string)

	isUpdating, _, err := secretmemo.IsUpdating(new, memoValue)
	isUnchanged := !isUpdating

	if nil != err {
		log.Printf("Change forced. Swallowing err while using secret hashing: %s", err)
		return false
	}

	log.Printf("\nk: %s, old: %s, new: %s, memoKey: %s, memoValue: %s, isUnchanged: %t\n",
		k, old, new, memoKey, memoValue, isUnchanged)
	return isUnchanged
}

// HelpFlattenSecretNested is used to store a hashed secret value into `tfstate`
func HelpFlattenSecretNested(d *schema.ResourceData, parentKey string, d2 map[string]interface{}, secretKey string) (string, string) {
	hashKey := calcSecretHashKey(secretKey)
	if len(d2) == 0 {
		return hashKey, ""
	}
	oldHash := d2[hashKey].(string)
	if !d.HasChange(parentKey) {
		log.Printf("key %s didn't get updated.", parentKey)
		return oldHash, hashKey
	}
	newSecret := d2[secretKey].(string)
	_, newHash, err := secretmemo.IsUpdating(newSecret, oldHash)
	if nil != err {
		log.Printf("Swallowing err while using secret hashing: %s", err)
	}
	log.Printf("Secret has changed. It's new hash value is %s.", newHash)
	return newHash, hashKey
}

// HelpFlattenSecret is used to store a hashed secret value into `tfstate`
func HelpFlattenSecret(d *schema.ResourceData, secretKey string) {
	if !d.HasChange(secretKey) {
		log.Printf("Secret key %s didn't get updated.", secretKey)
		return
	}
	hashKey := calcSecretHashKey(secretKey)
	newSecret := d.Get(secretKey).(string)
	oldHash := d.Get(hashKey).(string)
	_, newHash, err := secretmemo.IsUpdating(newSecret, oldHash)
	if nil != err {
		log.Printf("Swallowing err while using secret hashing: %s", err)
	}
	log.Printf("Secret key %s is updated. It's new hash key and value is %s and %s.", secretKey, hashKey, newHash)
	d.Set(hashKey, newHash)
}

// GenerateSecreteMemoSchema is used to create Schema defs to house the hashed secret in `tfstate`
func GenerateSecreteMemoSchema(secretKey string) (string, *schema.Schema) {
	out := schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Default:     nil,
		Description: fmt.Sprintf("A bcrypted hash of the attribute '%s'", secretKey),
		Sensitive:   true,
	}
	return calcSecretHashKey(secretKey), &out
}

// ParseProjectIDAndResourceID parses from the schema's resource data.
func ParseProjectIDAndResourceID(d *schema.ResourceData) (string, int, error) {
	projectID := d.Get("project_id").(string)
	resourceID, err := strconv.Atoi(d.Id())

	return projectID, resourceID, err
}

// ParseImportedID parse the imported int Id from the terraform import
func ParseImportedID(id string) (string, int, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || strings.EqualFold(parts[0], "") || strings.EqualFold(parts[1], "") {
		return "", 0, fmt.Errorf("unexpected format of ID (%s), expected projectid/resourceId", id)
	}
	project := parts[0]
	resourceID, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("error expected a number but got: %+v", err)
	}
	return project, resourceID, nil
}

// ParseImportedName parse the imported Id (Name) from the terraform import
func ParseImportedName(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || strings.EqualFold(parts[0], "") || strings.EqualFold(parts[1], "") {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected projectid/resourceName", id)
	}
	project := parts[0]
	resourceID := parts[1]

	return project, resourceID, nil
}

// ParseImportedUUID parse the imported uuid from the terraform import
func ParseImportedUUID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || strings.EqualFold(parts[0], "") || strings.EqualFold(parts[1], "") {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected projectid/resourceId", id)
	}
	project := parts[0]
	_, err := uuid.ParseUUID(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("%s isn't a valid UUID", parts[1])
	}
	return project, parts[1], nil
}

// ExpandStringList expand an array of interface into array of string
func ExpandStringList(d []interface{}) []string {
	vs := make([]string, 0, len(d))
	for _, v := range d {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// ExpandStringSet expand a set into array of string
func ExpandStringSet(d *schema.Set) []string {
	return ExpandStringList(d.List())
}

// ImportProjectQualifiedResource Import a resource by an ID that looks like one of the following:
//
//	<project ID>/<resource ID>
//	<project name>/<resource ID>
func ImportProjectQualifiedResource() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			projectNameOrID, resourceID, err := ParseImportedName(d.Id())

			if err != nil {
				return nil, fmt.Errorf("error parsing the resource ID from the Terraform resource data: %v", err)
			}

			if projectNameOrID, err = GetRealProjectId(projectNameOrID, meta); err == nil {
				d.Set("project_id", projectNameOrID)
				d.SetId(resourceID)
				return []*schema.ResourceData{d}, nil
			}
			return nil, err
		},
	}
}

// ImportProjectQualifiedResourceInteger Import a resource by an ID that looks like one of the following:
//
//	<project ID>/<resource ID as integer>
//	<project name>/<resource ID as integer>
func ImportProjectQualifiedResourceInteger() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			projectNameOrID, resourceID, err := ParseImportedName(d.Id())

			if err != nil {
				return nil, fmt.Errorf("error parsing the resource ID from the Terraform resource data: %v", err)
			}

			_, err = strconv.Atoi(resourceID)
			if err != nil {
				return nil, fmt.Errorf("resource ID was expected to be integer, but was not: %+v", err)
			}

			if projectNameOrID, err = GetRealProjectId(projectNameOrID, meta); err == nil {
				d.Set("project_id", projectNameOrID)
				d.SetId(resourceID)
				return []*schema.ResourceData{d}, nil
			}
			return nil, err
		},
	}
}

// ImportProjectQualifiedResourceUUID Import a resource by an ID that looks like one of the following:
//
//	<project ID>/<resource ID as uuid>
//	<project name>/<resource ID as uuid>
func ImportProjectQualifiedResourceUUID() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			projectNameOrID, resourceID, err := ParseImportedUUID(d.Id())

			if err != nil {
				return nil, fmt.Errorf("error parsing the resource ID from the Terraform resource data: %v", err)
			}

			if projectNameOrID, err = GetRealProjectId(projectNameOrID, meta); err == nil {
				d.Set("project_id", projectNameOrID)
				d.SetId(resourceID)
				return []*schema.ResourceData{d}, nil
			}
			return nil, err
		},
	}
}

// Get real project ID
func GetRealProjectId(projectNameOrID string, meta interface{}) (string, error) {
	//If request params is project name, try get the project ID
	if _, err := uuid.ParseUUID(projectNameOrID); err != nil {
		clients := meta.(*client.AggregatedClient)
		project, err := clients.CoreClient.GetProject(clients.Ctx, core.GetProjectArgs{
			ProjectId:           &projectNameOrID,
			IncludeCapabilities: converter.Bool(true),
			IncludeHistory:      converter.Bool(false),
		})
		if err != nil {
			return "", fmt.Errorf(" Failed to get the project with specified projectNameOrID: %s , %+v", projectNameOrID, err)
		}
		return (*project.Id).String(), nil
	}
	return projectNameOrID, nil
}

// FindMapInSetWithGivenKeyValue Pulls an element of `TypeSet` from the state. The values of this set are assumed to be
// `TypeMap`. The maps in the set are searched until a map is found with a value for `keyName` equal to `keyValue`.
//
// If no such map is found, `nil` is returned
func FindMapInSetWithGivenKeyValue(d *schema.ResourceData, setName string, keyName string, keyValue interface{}) map[string]interface{} {
	for _, value := range d.Get(setName).(*schema.Set).List() {
		valueAsMap := value.(map[string]interface{})
		// Note: casing matters here so we will use `==` over `strings.EqualFold`
		if valueAsMap[keyName] == keyValue {
			return valueAsMap
		}
	}
	return nil
}
