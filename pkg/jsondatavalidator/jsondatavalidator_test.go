// +build unit

//https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go

package jsondatavalidator_test

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/vishwanathj/JSON-Parameterized-Data-Validator/pkg/jsondatavalidator"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// SchemaDir points to the relative path of where the schema files are located
var SchemaDir = "../../test/testdata/schema/"

// SchemaFileInputParam is name of schema file for input param files
var SchemaFileInputParam = "inputParam.json"

// SchemaFileDefineNonParam is name of schema file for non-parameterized templates.
// This is needed by the GenerateJSONSchemaFromParameterizedTemplate function
var SchemaFileDefineNonParam = "vnfdDefineNonParam.json"

var BASE_DIR = "../../test/testdata/yamlfiles/"
//var BASE_DIR_VALID_Parameterized_Input = "../../test/testdata/yamlfiles/valid/parameterizedInput/"
var BASE_DIR_VALID_Parameterized_Input = BASE_DIR + "valid/parameterizedInput/"
var BASE_DIR_VALID_Parameterized_Instance = BASE_DIR + "valid/parameterizedInstance/"
var BASE_DIR_INVALID_Parameterized_Input = BASE_DIR + "invalid/parameterizedInput/"
var BASE_DIR_INVALID_Parameterized_Instance = BASE_DIR + "invalid/parameterizedInstance/"
var BASE_DIR_VALID_NonParameterized_Input = BASE_DIR + "valid/nonParameterizedInput/"
var BASE_DIR_INVALID_NonParameterized_Input = BASE_DIR + "invalid/nonParameterizedInput/"
var BASE_DIR_VALID_Input_Param = BASE_DIR + "valid/inputParam/"
var BASE_DIR_INVALID_Input_Param = BASE_DIR + "invalid/inputParam/"
var BASE_DIR_VALID_Paginated = BASE_DIR + "valid/parameterizedPaginatedInstances/"
var BASE_DIR_INVALID_Paginated = BASE_DIR + "invalid/parameterizedPaginatedInstances/"

/*
func TestValidatePaginatedVnfdsInstancesBody_Positive(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Paginated)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidatePaginatedVnfdsInstancesBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidatePaginatedVnfdsInstancesBody_Negative(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Paginated)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidatePaginatedVnfdsInstancesBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}

func TestValidateVnfdPostBody_PositiveCases(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidateVnfdPostBody_NegativeCases(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidateVnfdPostBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}

func TestValidateVnfdInstanceBody_PositiveCases(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Log("Success")
			} else {
				t.Error(err)
			}

		}
	}
}

func TestValidateVnfdInstanceBody_NegativeCases(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			err = jsondatavalidator.ValidateVnfdInstanceBody(yamlText)

			if err == nil {
				t.Error("FAIL")
			} else {
				t.Log(err)
			}

		}
	}
}

func TestValidNonParameterizedInputYaml(t *testing.T) {

	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_NonParameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextNonParameterizedInput := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextNonParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				fmt.Println(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}
		}
	}
}

func TestInValidNonParameterizedInputYaml(t *testing.T) {

	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_NonParameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextNonParameterizedInput := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextNonParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}

			if err := schema.ValidateInterface(m); err != nil {
				t.Log(err)
			} else {
				t.Fail()
			}
		}
	}
}

func TestValidParameterizedInputYaml(t *testing.T) {

	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			schemaTextParameterizedInput := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				fmt.Println(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}

		}
	}
}

func TestInValidParameterizedInputYaml(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_INVALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInput())); err != nil {
			schemaTextParameterizedInput := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaInputPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInput)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				//panic(err)
				//fmt.Println(err)
				t.Logf("err")
			} else {
				t.Fail()
			}

		}
	}
}

func TestValidParameterizedInstanceYaml(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_VALID_Parameterized_Instance + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInstance())); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
			schemaTextParameterizedInstance := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaParameterizedInstanceRelPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				panic(err)
				t.Fail()
			} else {
				t.Log("Passed")
			}

		}
	}
}

func TestInValidParameterizedInstanceYaml(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Parameterized_Instance)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath + BASE_DIR_INVALID_Parameterized_Instance + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m interface{}
			err = yaml.Unmarshal(yamlText, &m)
			if err != nil {
				panic(err)
				t.Fail()
			}

			compiler := jsonschema.NewCompiler()
			//compiler.Draft = jsonschema.Draft4
			//if err := compiler.AddResource("schema.json", strings.NewReader(utilGetSchemaTextParameterizedInstance())); err != nil {
			//if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
			schemaTextParameterizedInstance := jsondatavalidator.GetSchemaStringWhenGivenFilePath(jsondatavalidator.SchemaParameterizedInstanceRelPath)
			if err := compiler.AddResource("schema.json", strings.NewReader(schemaTextParameterizedInstance)); err != nil {
				panic(err)
				t.Errorf("panic: AddResource ERROR")
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				panic(err)
				t.Errorf("panic: Compile ERROR")
			}
			if err := schema.ValidateInterface(m); err != nil {
				//panic(err)
				//fmt.Println(err)
				t.Logf("err")
			} else {
				t.Errorf("FAILED")
			}

		}
	}
}*/

func TestPositiveGetRegexMatchingListFromJSONBuff(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	// The regexp looks for the $ anywhere in the line
	validRegexList := [...]string{`\$[A-Za-z][-A-Za-z0-9_]*`, `\$`}

	if err != nil {
		t.Fatal(err)
	}
	for _, rxp := range validRegexList {
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" || filepath.Ext(f.Name()) == ".yaml" {
				fmt.Println(f.Name())
				//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
				yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
				if err != nil {
					t.Error("Error while reading VNFD File # ", err)
					// if file read fail the continue to next file.
					panic(err)
					t.Fail()
				}

				// The regexp looks for the $ anywhere in the line

				lst := jsondatavalidator.GetRegexMatchingListFromJSONBuff(yamlText, rxp)
				if len(lst) != 0 {
					t.Log(lst, len(lst))
				} else {
					t.Errorf("FAILED")
				}
			}
		}
	}

}

func TestPositive_GetEntireLinesFromYAMLFile_WhenRegexMatched_FromJSONBuff(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	// The regexp looks for the $ anywhere in the line and returns the entire line
	validRegexList := [...]string{`.*\$.*`}

	if err != nil {
		t.Fatal(err)
	}
	for _, rxp := range validRegexList {
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" || filepath.Ext(f.Name()) == ".yaml" {
				fmt.Println(f.Name())
				yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
				if err != nil {
					t.Error("Error while reading VNFD File # ", err)
					// if file read fail the continue to next file.
					panic(err)
					t.Fail()
				}

				// The regexp looks for the $ anywhere in the line

				lst := jsondatavalidator.GetRegexMatchingListFromJSONBuff(yamlText, rxp)
				if len(lst) != 0 {
					t.Log(lst, len(lst))
				} else {
					t.Errorf("FAILED")
				}
			}
		}
	}

}

func TestNegative_GetRegexMatching_List_FromJSONBuff(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" || filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			//yamlText, err := ioutil.ReadFile(gopath+BASE_DIR_VALID_Parameterized_Input + f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			// The regexp looks for the $ in the beginning of the line and should fail
			lst := jsondatavalidator.GetRegexMatchingListFromJSONBuff(yamlText, `^\$[A-Za-z][-A-Za-z0-9_]*`)
			if len(lst) != 0 {
				t.Error(lst)
			} else {
				t.Log(len(lst))
			}
		}
	}
}

/*
func TestParse_NestedMap_ToObtain_LeafKeyValues(t *testing.T) {
	bpath := GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)

	// The regexp looks for the $ anywhere in the line
	validRegexList := [...]string{ `.*\$.*`}

	if err != nil {
		t.Fatal(err)
	}
	for _, rxp := range validRegexList{
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" || filepath.Ext(f.Name()) == ".yaml"{
				fmt.Println(f.Name())
				yamlText, err := ioutil.ReadFile(bpath +"/" + f.Name())
				//fmt.Println(string(yamlText))

				if err != nil {
					t.Error("Error while reading VNFD File # ", err)
					// if file read fail the continue to next file.
					panic(err)
					t.Fail()
				}

				var m map[string]interface{}
				yaml.Unmarshal(yamlText, &m)
				//fmt.Println(m)

				pvm := NewSearchResults(MatchKey, `\$.*`)
				pvm.ParseMap(m)
				fmt.Println(pvm.Results)

				// The regexp looks for the $ anywhere in the line

				lst := GetRegexMatchingListFromJSONBuff(yamlText, rxp)

				if len(lst) != 0 {
					t.Log(pvm)
				} else {
					t.Errorf("FAILED")
				}
			}
		}
	}
}*/

func TestPositive_Parse_NestedJSONSchema_ToObtain_Values_GivenAKey(t *testing.T) {
	//bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(jsondatavalidator.SchemaDir)
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}

	yamlText, err := ioutil.ReadFile(bpath + "/" + "vnfdDefine.json")
	//yamlText, err := ioutil.ReadFile(bpath + "/" + "nonParameterizedVnfdDefine.json")
	//yamlText, err := ioutil.ReadFile(bpath +"/" + "parameterizedVnfdDefine.json")
	if err != nil {
		t.Error("Error while reading VNFD File # ", err)
		// if file read fail the continue to next file.
		panic(err)
		t.Fail()
	}

	//patternList := []string{"connection_point", "constraints", "scale_in_out", "vdu", "vnfc", "virtual_link"}
	patternList := []string{"vcpus", "ip_address", "dedicated", "vim_id", "high_availability", "memory", "disk_size", "image", "default", "maximum", "minimum"}

	var m map[string]interface{}
	yaml.Unmarshal(yamlText, &m)

	for _, elem := range patternList {
		pvm := jsondatavalidator.NewSearchResults(jsondatavalidator.MatchKey, elem)
		pvm.ParseMap(m)

		if len(pvm.Results) == 0 {
			t.Fail()
		} else {
			t.Log(elem)
			t.Log(len(pvm.Results))
			t.Log(pvm.Results)
		}
	}
}

func TestPositive_Parse_NestedJSONSchema_ToObtain_Key_GivenValue(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" || filepath.Ext(f.Name()) == ".yaml" {
			fmt.Println(f.Name())
			yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
			if err != nil {
				t.Error("Error while reading VNFD File # ", err)
				// if file read fail the continue to next file.
				panic(err)
				t.Fail()
			}

			var m map[string]interface{}
			yaml.Unmarshal(yamlText, &m)
			pvm := jsondatavalidator.NewSearchResults(jsondatavalidator.MatchValue, `\$.*`)
			pvm.ParseMap(m)

			if len(pvm.Results) == 0 {
				t.Fail()
			} else {
				t.Log(len(pvm.Results))
				t.Log(pvm.Results)
			}
		}
	}
}

func TestGenerateJSONSchemaFromParameterizedTemplate_Positive(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}

	yamlText, err := ioutil.ReadFile(bpath + "/" + "validParameterizedVNFDInputWithOptionalPropConstraintsMissing.yaml")
	if err != nil {
		t.Error("Error while reading VNFD File # ", err)
		// if file read fail the continue to next file.
		panic(err)
		t.Fail()
	}

	abspath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileDefineNonParam
	nonParamDefineJSONBuf, err := jsondatavalidator.GetSchemaDefinitionFileAsJSONBuf(abspath)
	if err != nil {
		t.Fatal(err)
	}

	abspath = jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileInputParam
	inputParamSchemaJSONBuf, err := jsondatavalidator.GetSchemaDefinitionFileAsJSONBuf(abspath)
	//inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(SchemaFileInputParam)
	if err != nil {
		t.Fatal(err)
	}

	//r, e := jsondatavalidator.GenerateJSONSchemaFromParameterizedTemplate(yamlText)
	r, e := jsondatavalidator.GenerateJSONSchemaFromParameterizedTemplate(yamlText, nonParamDefineJSONBuf, inputParamSchemaJSONBuf)

	if e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(r))
	}
}

/*
func TestPositive_ValidateInputParamAgainstParameterizedVnfd(t *testing.T) {
	tables := []struct {
		inputParamFileName    string
		parameterizedVnfdName string
	}{
		{"inputParamConstraintsMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropConstraintsMissing.yaml"},
		{"inputParamHAMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropHAMissing.yaml"},
		{"inputParamOptionalProps.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamScaleMissing.yaml",
			"validParameterizedVNFDInputWithOptionalPropScaleMissing.yaml"},
		{"inputParamsRequiredProps.yaml",
			"validParameterizedVNFDInputWithRequiredProps.yaml"},
		{"inputParamSubnetPools.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
	}

	vnfdpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	inputparampath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Input_Param)

	for _, table := range tables {
		inparam, ierr := ioutil.ReadFile(inputparampath + "/" + table.inputParamFileName)

		if ierr != nil {
			t.Fatal(ierr)
		}

		vnfd, verr := ioutil.ReadFile(vnfdpath + "/" + table.parameterizedVnfdName)
		if verr != nil {
			t.Fatal(verr)
		}

		err := jsondatavalidator.ValidateInputParamAgainstParameterizedVnfd(inparam, vnfd)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestNegative_ValidateInputParamAgainstParameterizedVnfd(t *testing.T) {
	tables := []struct {
		inputParamFileName    string
		parameterizedVnfdName string
	}{
		{"inputParamInvalidDedicatedConstraint.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidDiskSize.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidHA.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidImageName.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidIPAddress.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidMemory.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidMinScale.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidName.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVcpus.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVimConstraint.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
		{"inputParamInvalidVnfdIDFormat.yaml",
			"validParameterizedVNFDInputWithOptionalProps.yaml"},
	}

	vnfdpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_Parameterized_Input)
	inputparampath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_INVALID_Input_Param)

	for _, table := range tables {
		inparam, ierr := ioutil.ReadFile(inputparampath + "/" + table.inputParamFileName)

		if ierr != nil {
			t.Fatal(ierr)
		}

		vnfd, verr := ioutil.ReadFile(vnfdpath + "/" + table.parameterizedVnfdName)
		if verr != nil {
			t.Fatal(verr)
		}

		err := jsondatavalidator.ValidateInputParamAgainstParameterizedVnfd(inparam, vnfd)
		if err != nil {
			t.Log(err)
		} else {
			t.Fail()
		}
	}
}
*/

func TestGenerateJSONSchemaFromNonParameterizedVNFDTemplate_Positive(t *testing.T) {
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(BASE_DIR_VALID_NonParameterized_Input)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	if err != nil {
		t.Fatal(err)
	}

	yamlText, err := ioutil.ReadFile(bpath + "/" + "validNonParameterizedVNFDInputWithOptionalPropConstraintsMissing.yaml")
	if err != nil {
		t.Error("Error while reading VNFD File # ", err)
		// if file read fail the continue to next file.
		panic(err)
		t.Fail()
	}

	abspath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileDefineNonParam
	nonParamDefineJSONBuf, err := jsondatavalidator.GetSchemaDefinitionFileAsJSONBuf(abspath)
	if err != nil {
		t.Fatal(err)
	}

	abspath = jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir) + "/" + SchemaFileInputParam
	inputParamSchemaJSONBuf, err := jsondatavalidator.GetSchemaDefinitionFileAsJSONBuf(abspath)
	//inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(SchemaFileInputParam)
	if err != nil {
		t.Fatal(err)
	}

	//r, e := jsondatavalidator.GenerateJSONSchemaFromParameterizedTemplate(yamlText)
	r, e := jsondatavalidator.GenerateJSONSchemaFromParameterizedTemplate(yamlText, nonParamDefineJSONBuf, inputParamSchemaJSONBuf)

	if e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(r))
	}
}

func TestPositive_GetEntireLinesFromJSONSchemaFile_WhenStringMatched_FromJSONBuff(t *testing.T) {
	//bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(jsondatavalidator.SchemaDir)
	bpath := jsondatavalidator.GetAbsDIRPathGivenRelativePath(SchemaDir)
	files, err := ioutil.ReadDir(bpath)
	fmt.Println("len:=", len(files))
	fmt.Println(bpath)

	// The regexp looks for the $ anywhere in the line and returns the entire line
	validRegexList := [...]string{`name`, `high_availability`, `image`, `memory`, `maximum`, `minimum`, `disk_size`, `vcpus`, `is_management`, `default`, `ip_address`}

	if err != nil {
		t.Fatal(err)
	}

	for _, rxp := range validRegexList {
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" {
				fmt.Println(f.Name())
				yamlText, err := ioutil.ReadFile(bpath + "/" + f.Name())
				if err != nil {
					t.Error("Error while reading VNFD File # ", err)
					// if file read fail the continue to next file.
					panic(err)
					t.Fail()
				}

				// The regexp looks for the $ anywhere in the line

				lst := jsondatavalidator.GetRegexMatchingListFromJSONBuff(yamlText, `.*`+rxp+`.*`)
				if len(lst) != 0 {
					fmt.Println(lst)
					t.Log(lst, len(lst))
				} /*else {
					t.Errorf("FAILED")
				}*/
			}
		}
	}
}

func TestValidateJSONBufAgainstSchema_Negative_InvalidJSON(t *testing.T) {
	malformedJson := []byte(`{"key":`)
	err := jsondatavalidator.ValidateJSONBufAgainstSchema(malformedJson, strings.NewReader("dummy"), "d")
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}

func TestValidateJSONBufAgainstSchema_Negative_InvalidURL(t *testing.T) {
	jsonval := []byte(`{"key": "val"}`)
	err := jsondatavalidator.ValidateJSONBufAgainstSchema(jsonval, strings.NewReader("dummy"), "d")
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}

func TestGetSchemaStringWhenGivenFilePath(t *testing.T) {

	var dir string
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		t.Errorf("Fatal error: %s", err)
	}
	parent := filepath.Dir(dir)

	testTable := [] struct{
		inputPath      string
		expectedOutput string
	} {
		{"schema/vnfdInstanceSchema.json#/vnfdInstance", `{"$ref": "` + "schema/vnfdInstanceSchema.json#/vnfdInstance" + `"}`},
		{"../schema/vnfdInstanceSchema.json#/vnfdInstance", `{"$ref": "` + parent + "/" + "schema/vnfdInstanceSchema.json#/vnfdInstance" + `"}`},
		{"", `{"$ref": ""}`},
	}
	for _, tdr := range testTable {
		res := jsondatavalidator.GetSchemaStringWhenGivenFilePath(tdr.inputPath)


		if (res != tdr.expectedOutput) {
			fmt.Println(res)
			t.Errorf("Output %s incorrect", tdr.expectedOutput)
		}
	}
}

/*
func TestValidateJSONBufAgainstSchema_Negative_FailCompile(t *testing.T) {
	jsonval := []byte(`{"type": "object","properties": {"name": {"type": "integer"}},"additionalProperties": false}`)
	err := ValidateJSONBufAgainstSchema(jsonval, strings.NewReader(`{"name": "hola"}`), "h.json")
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}*/
