package tfhelper

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/secretmemo"
	"log"
	"strconv"
	"strings"
)

func calcSecretHashKey(secretKey string) string {
	return secretKey + "_hash"
}

// DiffFuncSupressSecretChanged is used to supress unneeded `apply` updates to a resource.
//
// It returns `true` when `new` appears to be the same value
// as a previously stored and bcrypt'd value stored in state during a previous `apply`.
// Relies on flatten/expand logic to help store that hash. See FlattenSecret, below.*/
func DiffFuncSupressSecretChanged(k, old, new string, d *schema.ResourceData) bool {
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
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Printf(string(b))
	}
	return
}

// ParseImportedID parse the imported Id from the terraform import
func ParseImportedID(id string) (string, int, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", 0, fmt.Errorf("unexpected format of ID (%s), expected projectid/resourceId", id)
	}
	project := parts[0]
	resourceID, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("Error converting getting the resource id: %+v", err)
	}
	return project, resourceID, nil
}
