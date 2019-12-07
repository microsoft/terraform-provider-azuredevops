// +build all helper converter

package converter

import (
	"fmt"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	value := "Hello World"
	valuePtr := String(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestInt(t *testing.T) {
	value := 123456
	valuePtr := Int(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestBoolTrue(t *testing.T) {
	value := true
	valuePtr := Bool(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestBoolFalse(t *testing.T) {
	value := false
	valuePtr := Bool(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestLicenseTypeAccount(t *testing.T) {
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.None)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.EarlyAdopter)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Advanced)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Professional)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Express)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Professional)

	_, err := AccountLicenseType("foo")
	assert.Equal(t, err.Error(), "Error unable to match given AccountLicenseType:foo")
}

func assertAccountLicenseType(t *testing.T, accountLicenseType licensing.AccountLicenseType) {
	actualAccountLicenseType, err := AccountLicenseType(string(accountLicenseType))
	assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", string(accountLicenseType)))
	assert.Equal(t, &accountLicenseType, actualAccountLicenseType, fmt.Sprintf("%s should be able to convert into the AccountLicenseType", string(accountLicenseType)))
}
