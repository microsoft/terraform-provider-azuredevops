// +build all helper converter

package converter

import (
	"testing"
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
