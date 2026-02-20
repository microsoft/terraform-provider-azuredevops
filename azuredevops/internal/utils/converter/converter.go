package converter

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/licensing"
)

// String Get a pointer to a string
func String(value string) *string {
	return &value
}

// StringFromInterface get a string pointer from an interface
func StringFromInterface(value interface{}) *string {
	return String(value.(string))
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// Int Get a pointer to an integer value
func Int(value int) *int {
	return &value
}

func ToPtr[E any](e E) *E {
	return &e
}

// ASCIIToIntPtr Convert a string to an Int Pointer
func ASCIIToIntPtr(value string) (*int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return Int(i), nil
}

// UInt64 Get a pointer to an uint64 value
func UInt64(value uint64) *uint64 {
	return &value
}

// ToString Given a pointer return its value, or a default value of the pointer is nil
func ToString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}

// ToBool Given a pointer return its value, or a default value of the pointer is nil
func ToBool(value *bool, defaultValue bool) bool {
	if value != nil {
		return *value
	}

	return defaultValue
}

// ToInt Given a pointer return its value, or a default value if the pointer is nil
func ToInt(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}

	return defaultValue
}

// AccountLicenseType Get a pointer to an AccountLicenseType
func AccountLicenseType(accountLicenseTypeValue string) (*licensing.AccountLicenseType, error) {
	var accountLicenseType licensing.AccountLicenseType
	switch strings.ToLower(accountLicenseTypeValue) {
	case "none":
		accountLicenseType = licensing.AccountLicenseTypeValues.None
	case "earlyadopter":
		accountLicenseType = licensing.AccountLicenseTypeValues.EarlyAdopter
	case "basic":
		fallthrough
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

// AccountLicensingSource convert a string value to a licensing.AccountLicenseType pointer
func AccountLicensingSource(licensingSourceValue string) (*licensing.LicensingSource, error) {
	var licensingSource licensing.LicensingSource
	switch strings.ToLower(licensingSourceValue) {
	case "none":
		licensingSource = licensing.LicensingSourceValues.None
	case "account":
		licensingSource = licensing.LicensingSourceValues.Account
	case "msdn":
		licensingSource = licensing.LicensingSourceValues.Msdn
	case "profile":
		licensingSource = licensing.LicensingSourceValues.Profile
	case "auto":
		licensingSource = licensing.LicensingSourceValues.Auto
	case "trial":
		licensingSource = licensing.LicensingSourceValues.Trial
	default:
		return nil, fmt.Errorf("Error unable to match given LicensingSource :%s", licensingSourceValue)
	}
	return &licensingSource, nil
}

// UUID converts a string to a pointer to a UUID, will panic if the string can't be parsed to a UUID
func UUID(szuuid string) *uuid.UUID {
	uuid := uuid.MustParse(szuuid)
	return &uuid
}

// DecodeUtf16HexString decodes a binary representation of an UTF16 string
func DecodeUtf16HexString(message string) (string, error) {
	b, err := hex.DecodeString(message)
	if err != nil {
		return "", err
	}
	ints := make([]uint16, len(b)/2)
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &ints); err != nil {
		return "", err
	}
	return string(utf16.Decode(ints)), nil
}

// EncodeUtf16HexString encodes a string into an binary representation with UTF16 enoding
func EncodeUtf16HexString(message string) (string, error) {
	runeByte := []rune(message)
	encodedByte := utf16.Encode(runeByte)
	var sb strings.Builder
	for i := 0; i < len(encodedByte); i++ {
		fmt.Fprintf(&sb, "%02x%02x", byte(encodedByte[i]), byte(encodedByte[i]>>8))
	}
	return sb.String(), nil
}
