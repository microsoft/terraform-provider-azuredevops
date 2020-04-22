package tfhelper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/secretmemo"
)

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

//PrettyPrint json
//func PrettyPrint(v interface{}) (err error) {
//	b, err := json.MarshalIndent(v, "", "  ")
//	if err == nil {
//		log.Printf(string(b))
//	}
//	return
//}

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
