package converter

import (
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
)

// String Get a pointer to a string
func String(value string) *string {
	return &value
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// Int Get a pointer to an integer value
func Int(value int) *int {
	return &value
}

// ToString Given a pointer return its value, or a default value of the poitner is nil
func ToString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}

// AccountLicenseType Get a pointer to an AccountLicenseType
func AccountLicenseType(accountLicenseTypeValue string) (*licensing.AccountLicenseType, error) {
	var accountLicenseType licensing.AccountLicenseType
	switch accountLicenseTypeValue {
	case "none":
		accountLicenseType = licensing.AccountLicenseTypeValues.None
	case "earlyAdopter":
		accountLicenseType = licensing.AccountLicenseTypeValues.EarlyAdopter
	case "express":
		accountLicenseType = licensing.AccountLicenseTypeValues.Express
	case "professional":
		accountLicenseType = licensing.AccountLicenseTypeValues.Professional
	case "advanced":
		accountLicenseType = licensing.AccountLicenseTypeValues.Advanced
	case "stakeholder":
		accountLicenseType = licensing.AccountLicenseTypeValues.Stakeholder
	default:
		return nil, fmt.Errorf("Error unable to match given AccountLicenseType:%s", accountLicenseTypeValue)
	}
	return &accountLicenseType, nil
}
