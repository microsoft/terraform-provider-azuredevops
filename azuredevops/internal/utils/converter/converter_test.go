//go:build all || helper || converter
// +build all helper converter

package converter

import (
	"fmt"
	"testing"

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

func TestASCIIToIntPtrErrorCase(t *testing.T) {
	type TestCase struct {
		testName  string
		input     string
		outputVal *int
		hasError  bool
	}
	cases := []TestCase{
		{
			testName:  "Positive Int",
			input:     "100",
			outputVal: Int(100),
			hasError:  false,
		}, {
			testName:  "Negative Int",
			input:     "-100",
			outputVal: Int(-100),
			hasError:  false,
		}, {
			testName:  "Zero",
			input:     "0",
			outputVal: Int(0),
			hasError:  false,
		}, {
			testName:  "Empty String",
			input:     "",
			outputVal: nil,
			hasError:  true,
		}, {
			testName:  "Not an Int",
			input:     "Hello World!",
			outputVal: nil,
			hasError:  true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			val, err := ASCIIToIntPtr(tc.input)

			hasError := err != nil
			if hasError != tc.hasError {
				t.Fatal("Expectation of tests error scenario was not met")
			}

			if hasError {
				return
			}

			if *val != *tc.outputVal {
				t.Fatalf("Expected output value to be %+v but was %+v", tc.outputVal, val)
			}
		})
	}
}

func TestStringFromInterface_StringValue(t *testing.T) {
	value := "Hello World"
	valuePtr := StringFromInterface(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestStringFromInterface_InterfaceValue(t *testing.T) {
	value := "Hello World"
	var interfaceValue interface{}

	interfaceValue = value
	valuePtr := StringFromInterface(interfaceValue)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

type encodeTestType struct {
	plainString   string
	encodedString string
}

var encodeTestCases = []encodeTestType{
	{
		plainString:   "branch_1_1",
		encodedString: "6200720061006e00630068005f0031005f003100",
	},
	{
		plainString:   "master",
		encodedString: "6d0061007300740065007200",
	},
	{
		plainString:   "refs/heads/main",
		encodedString: "72006500660073002f00680065006100640073002f006d00610069006e00",
	},
	{
		plainString:   "feature",
		encodedString: "6600650061007400750072006500",
	},
	{
		plainString:   "A",
		encodedString: "4100",
	},
	{
		plainString:   "test_branch",
		encodedString: "74006500730074005f006200720061006e0063006800",
	},
	{
		plainString:   "Hello",
		encodedString: "480065006c006c006f00",
	},
}

func TestDecodeUtf16HexString(t *testing.T) {
	for _, etest := range encodeTestCases {
		val, err := DecodeUtf16HexString(etest.encodedString)
		assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", etest.encodedString))
		assert.EqualValues(t, etest.plainString, val)
	}
}

func TestEncodeUtf16HexString(t *testing.T) {
	for _, etest := range encodeTestCases {
		val, err := EncodeUtf16HexString(etest.plainString)
		assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", etest.plainString))
		assert.EqualValues(t, etest.encodedString, val)
	}
}

// TestEncodeUtf16HexString_EdgeCases tests edge cases and special characters
func TestEncodeUtf16HexString_EdgeCases(t *testing.T) {
	testCases := []encodeTestType{
		{
			plainString:   "",
			encodedString: "",
		},
		{
			plainString:   "!@#$%^&*()",
			encodedString: "210040002300240025005e0026002a0028002900",
		},
		{
			plainString:   "123",
			encodedString: "310032003300",
		},
		{
			plainString:   " ", // space
			encodedString: "2000",
		},
		{
			plainString:   "a/b\\c",
			encodedString: "61002f0062005c006300",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Encode_%s", tc.plainString), func(t *testing.T) {
			val, err := EncodeUtf16HexString(tc.plainString)
			assert.Nil(t, err)
			assert.EqualValues(t, tc.encodedString, val, fmt.Sprintf("Expected %s but got %s", tc.encodedString, val))
		})
	}
}

// TestEncodeDecodeUtf16HexString_RoundTrip ensures encode/decode are inverse operations
func TestEncodeDecodeUtf16HexString_RoundTrip(t *testing.T) {
	testStrings := []string{
		"Hello World",
		"terraform-provider-azuredevops",
		"refs/heads/feature/test-branch",
		"aBcDeFgHiJkLmNoPqRsTuVwXyZ",
		"0123456789",
		"!@#$%^&*()",
		"",
		"A",
		"test_underscore_name",
	}

	for _, original := range testStrings {
		t.Run(fmt.Sprintf("RoundTrip_%s", original), func(t *testing.T) {
			// Encode
			encoded, err := EncodeUtf16HexString(original)
			assert.Nil(t, err, "Encoding should not error")

			// Decode
			decoded, err := DecodeUtf16HexString(encoded)
			assert.Nil(t, err, "Decoding should not error")

			// Verify round-trip
			assert.Equal(t, original, decoded, "Round-trip should preserve original string")
		})
	}
}
