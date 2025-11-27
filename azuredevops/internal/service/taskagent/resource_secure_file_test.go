package taskagent

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCalculateContentHashes(t *testing.T) {
	content := "test-content"
	sha1Expected := "60b62e43b6a5e292b8fdbd41e57de248605d2c27"
	sha256Expected := "0a3666a0710c08aa6d0de92ce72beeb5b93124cce1bf3701c9d6cdeb543cb73e"
	sha1, sha256 := calculateContentHashes(content)
	if sha1 == "" || sha256 == "" {
		t.Errorf("Hashes should not be empty")
	}
	if sha1 != sha1Expected {
		t.Errorf("SHA1 mismatch")
	}
	if sha256 != sha256Expected {
		t.Errorf("SHA256 mismatch")
	}
}

func TestBuildPropertiesMap(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceSecureFile().Schema, map[string]interface{}{
		"properties":       map[string]interface{}{"foo": "bar"},
		"file_hash_sha1":   "abc",
		"file_hash_sha256": "def",
	})
	props := buildPropertiesMap(d)
	expected := map[string]string{"foo": "bar", "file_hash_sha1": "abc", "file_hash_sha256": "def"}
	if !reflect.DeepEqual(props, expected) {
		t.Errorf("Expected %v, got %v", expected, props)
	}
}
