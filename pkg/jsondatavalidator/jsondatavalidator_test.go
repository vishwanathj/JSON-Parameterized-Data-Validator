// +build unit

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package jsondatavalidator_test

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/vishwanathj/JSON-Parameterized-Data-Validator/pkg/jsondatavalidator"
)

var testJSONParamNonParamSchema = []byte(`
{
  "vmDeviceDefine": {
    "vm": {
      "additionalProperties": false, 
      "type": "object", 
      "required": [
        "vcpus"
      ], 
      "optional": [
        "memory"
      ], 
      "properties": {
        "vcpus": {
          "oneOf": [
            {
              "pattern": "^\\$[A-Za-z][-A-Za-z0-9_]*$", 
              "type": "string"
            }, 
            {
              "minimum": 2, 
              "type": "integer", 
              "maximum": 16, 
              "multipleOf": 2.0
            }
          ]
        },
        "memory": {
          "oneOf": [
            {
              "pattern": "^\\$[A-Za-z][-A-Za-z0-9_]*$", 
              "type": "string"
            }, 
            {
              "minimum": 512, 
              "type": "integer", 
              "maximum": 16384, 
              "multipleOf": 512
            }
          ]
        }
      }
    }
  }
}
`)

var testJSONNonParamSchema = []byte(`
{
  "vmDeviceDefine": {
    "vm": {
      "additionalProperties": false, 
      "type": "object", 
      "required": [
        "vcpus"
      ], 
      "optional": [
        "memory"
      ], 
      "properties": {
        "vcpus": {
          "minimum": 2, 
          "type": "integer", 
          "maximum": 16, 
          "multipleOf": 2
        }, 
        "name": {
          "pattern": "^[A-Za-z][-A-Za-z0-9_]*$", 
          "type": "string"
        }, 
        "memory": {
          "minimum": 512, 
          "type": "integer", 
          "maximum": 16384, 
          "multipleOf": 512
        }
      }
    }
  }
}
`)

var testInputParamJSONSchema = []byte(
	`
{
  "inputParam": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "pattern": "^[A-Za-z][-A-Za-z0-9_]*$"
      },
      "vm_id": {
        "type": "string",
        "pattern": "^VM-[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
      }
    },
    "required": [
      "name",
      "vm_id"
    ],
    "additionalProperties": false
  }
}
`)

func TestGetRegexMatchingListFromJSONBuff(t *testing.T) {
	var testJSONData = []byte(`{"devices": [{"path": "/sys/disk","name": "$parameterized_name"}]}`)

	testTable := []struct {
		description    string
		regExInput     string
		expectedOutput []string
		testData       []byte
	}{
		{"look for $ anywhere in the line and returns substring from point of match", `\$[A-Za-z][-A-Za-z0-9_]*`, []string{"$parameterized_name"}, testJSONData},
		{"look for $ anywhere in the line and returns entire line", `.*\$.*`, []string{`{"devices": [{"path": "/sys/disk","name": "$parameterized_name"}]}`}, testJSONData},
		{"look for $ beginning of line and returns substring from point of match", `^\$[A-Za-z][-A-Za-z0-9_]*`, nil, testJSONData},
	}
	for i, tc := range testTable {
		t.Run(fmt.Sprintf("%d:%s", i, tc.description), func(t *testing.T) {
			lst := jsondatavalidator.GetRegexMatchingListFromJSONBuff(tc.testData, tc.regExInput)
			t.Log(lst)
			if reflect.DeepEqual(tc.expectedOutput, lst) {
			} else {
				t.Error()
			}
		})
	}
}

func TestNewSearchResults(t *testing.T) {
	var testJSONData = []byte(`{"devices": [{"path": "/sys/disk","name": "$parameterized_name"}]}`)

	testTable := []struct {
		description string
		testData    []byte
		paramKeyVal int
		paramInput  string
		expectedOut []string
	}{
		{"Obtain Values Given a Key", testJSONParamNonParamSchema, jsondatavalidator.MatchKey, "vcpus", []string{"pattern"}},
		{"Obtain Key Given a Value", testJSONData, jsondatavalidator.MatchValue, `\$.*`, []string{"name"}},
	}
	for i, tdr := range testTable {
		t.Run(fmt.Sprintf("%d:%s", i, tdr.description), func(t *testing.T) {
			var m map[string]interface{}
			err := yaml.Unmarshal(tdr.testData, &m)
			if err != nil {
				t.Fatal(err)
			}
			pvm := jsondatavalidator.NewSearchResults(jsondatavalidator.MatchValue, `\$.*`)
			pvm.ParseMap(m)

			if tdr.expectedOut[0] == fmt.Sprint(pvm.Results[0]) {
				t.Log(pvm.Results)
			} else {
				t.Error()
			}
		})
	}
}

func TestGenerateJSONSchemaFromParameterizedTemplate(t *testing.T) {
	var regExpStr1 = `.*\$.*`
	var regExpStr2 = `.*\{.*`

	var testValidAllPropsNonParameterizedData = []byte(`
vm:
  vcpus: 4
  memory: 1024
`)

	var testValidAllPropsParameterizedData1 = []byte(`
vm:
  vcpus: $vcpus
  memory: $memory
`)
	var testValidAllPropsParameterizedData2 = []byte(`
vm:
  vcpus: {vcpus
  memory: {memory
`)

	testTable := []struct {
		description             string
		testJSON                []byte
		nonParamDefineJSONBuf   []byte
		inputParamSchemaJSONBuf []byte
		keysToAddToRequiredSec 	[]string
		regExpStr 				string
	}{
		{"Parameterized (`$`) Template test", testValidAllPropsParameterizedData1, testJSONNonParamSchema, testInputParamJSONSchema, []string{"vm_id", "name"}, regExpStr1},
		{"Non Parameterized Template test", testValidAllPropsNonParameterizedData, testJSONNonParamSchema, testInputParamJSONSchema, []string{"vm_id", "name"}, regExpStr1},
		{"Parameterized (`{`) Template test", testValidAllPropsParameterizedData2, testJSONNonParamSchema, testInputParamJSONSchema, []string{"vm_id", "name"}, regExpStr2},
	}

	for i, tdr := range testTable {
		t.Run(fmt.Sprintf("%d:%s", i, tdr.description), func(t *testing.T) {
			r, e := jsondatavalidator.GenerateJSONSchemaFromParameterizedTemplate(tdr.testJSON, tdr.nonParamDefineJSONBuf, tdr.inputParamSchemaJSONBuf, tdr.keysToAddToRequiredSec, tdr.regExpStr)
			t.Log(string(r))
			if r == nil && e != nil {
				t.Fatal("ERROR: JSONSchema failed to be generated.")
			}
		})
	}
}

func TestValidateJSONBufAgainstSchema(t *testing.T) {
	testValidSchema := `{"type": "object", "properties": {"vm": {"additionalProperties": false, "type": "object", "required": ["vcpus"], "optional": ["memory"], "properties": {"vcpus": {"oneOf": [{"pattern": "^\\$[A-Za-z][-A-Za-z0-9_]*$", "type": "string"}, {"minimum": 2, "type": "integer", "maximum": 16, "multipleOf": 2.0}]},"memory": {"oneOf": [{"pattern": "^\\$[A-Za-z][-A-Za-z0-9_]*$", "type": "string"}, {"minimum": 512, "type": "integer", "maximum": 16384, "multipleOf": 512}]}}}}}`
	testValidJSONData := []byte(`{"vm": {"vcpus": "$vcpus","memory": "$memory"}}`)
	//testInValidJSONData := []byte(`{"vm": {"cpus": "$vcpus","memory": "$memory"}}`)
	testInValidJSONData := []byte(`{"vm": {"memory": "$memory"}}`)
	testInvalidAdditionalProperty := []byte(`{"vm": {"vcpus":"$vcpus", "proc":"$proc"}}`)

	testTable := []struct {
		description          string
		jsonval              []byte
		schemaDefAsReaderObj io.Reader
		url                  string
		expectedOutput       error
	}{
		{"Invalid URL", []byte(`{"key": "val"}`), strings.NewReader("dummy"), "d", fmt.Errorf("AddResourceError")},
		{"Malformed JSON", []byte(`{"key":`), strings.NewReader("dummy"), "d", fmt.Errorf("UnMarshallError")},
		{"Valid JSON", testValidJSONData, strings.NewReader(string(testValidSchema)), "sch.json", nil},
		{"Invalid: additional property", testInvalidAdditionalProperty, strings.NewReader(string(testValidSchema)), "sch.json", fmt.Errorf("I[#/vm] S[#/properties/vm/additionalProperties] additionalProperties \"proc\" not allowed")},
		{"Invalid: missing required property", testInValidJSONData, strings.NewReader(string(testValidSchema)), "sch.json", fmt.Errorf("I[#/vm] S[#/properties/vm/required] missing properties: \"vcpus\"")},
	}
	for i, tdr := range testTable {
		t.Run(fmt.Sprintf("%d:%s", i, tdr.description), func(t *testing.T) {
			err := jsondatavalidator.ValidateJSONBufAgainstSchema(tdr.jsonval, tdr.schemaDefAsReaderObj, tdr.url)
			if err != nil {
				//fmt.Println(err)
				t.Log(err.Error())
			}

			if err == tdr.expectedOutput {

			} else if strings.TrimSpace(err.Error()) != tdr.expectedOutput.Error() {
				t.Log(len(err.Error()))
				t.Log(len(tdr.expectedOutput.Error()))
				t.Errorf("%s", tdr.expectedOutput)
			}
		})

	}
}