package jsondatavalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/peterbourgon/mergemap"
	log "github.com/sirupsen/logrus"

	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema"
)

// SchemaDir points to the relative path of where the schema files are located
var SchemaDir string

// SchemaInputPath is path to schema file for
// Parameterized templates
var SchemaInputPath string

// SchemaParameterizedInstanceRelPath is path to schema file for
// instantiated Parameterized templates
var SchemaParameterizedInstanceRelPath string

// SchemaPaginatedInstancesRelPath is path to schema file
// for paginated output structure
var SchemaPaginatedInstancesRelPath string

// SchemaFileInputParam is name of schema file for input param files
var SchemaFileInputParam string

// SchemaFileDefineNonParam is name of schema file for non-parameterized templates.
// This is needed by the GenerateJSONSchemaFromParameterizedTemplate function
var SchemaFileDefineNonParam string

func init() {
	log.Debug()
	localUnitTest := os.Getenv("TEST")
	log.Debug(localUnitTest)

	SchemaDir = "../schema/"
	SchemaInputPath = "../schema/vnfdInputSchema.json#/vnfdInput"
	SchemaParameterizedInstanceRelPath = "../schema/vnfdInstanceSchema.json#/vnfdInstance"
	SchemaPaginatedInstancesRelPath = "../schema/vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances"
	SchemaFileInputParam = "inputParam.json"
	SchemaFileDefineNonParam = "vnfdDefineNonParam.json"
	/*
	if localUnitTest == "true" {
		SchemaDir = "../schema/"
		SchemaInputPath = "../schema/vnfdInputSchema.json#/vnfdInput"
		SchemaParameterizedInstanceRelPath = "../schema/vnfdInstanceSchema.json#/vnfdInstance"
		SchemaPaginatedInstancesRelPath = "../schema/vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances"
		SchemaFileInputParam = "inputParam.json"
		SchemaFileDefineNonParam = "vnfdDefineNonParam.json"
	} else {
		SchemaDir = "/usr/share/vnfdservice/schema/"
		SchemaInputPath = "/usr/share/vnfdservice/schema/vnfdInputSchema.json#/vnfdInput"
		SchemaParameterizedInstanceRelPath = "/usr/share/vnfdservice/schema/vnfdInstanceSchema.json#/vnfdInstance"
		SchemaPaginatedInstancesRelPath = "/usr/share/vnfdservice/schema/vnfdPaginatedInstanceSchema.json#/vnfdsPaginatedInstances"
		SchemaFileInputParam = "inputParam.json"
		SchemaFileDefineNonParam = "vnfdDefineNonParam.json"
	}*/
}

const (
	// MatchKey is set when keys in a map have to be matched
	MatchKey = 1
	// MatchValue is set when values in a map have to be matched
	MatchValue = 2
	// KeyInputParam holds name of a key in a map
	KeyInputParam = "inputParam"
	// KeyRequired holds name of a key in a map
	KeyRequired = "required"
	// KeyVnfdID holds name of a key in a map
	KeyVnfdID = "vnfd_id"
	// KeyName holds name of a key in a map
	KeyName = "name"
	// KeyProperties holds name of a key in a map
	KeyProperties = "properties"
)

// SearchResults stores the results when parsing a map structure for
// a certain pattern
type SearchResults struct {
	SearchType          int
	SearchPatternString string
	Results             []interface{}
	re                  regexp.Regexp
}

// contains function checks if a given type "d" exists in the list "s"
func contains(s []interface{}, d interface{}) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, d) {
			return true
		}
	}
	return false
}

// NewSearchResults creates a new struct type "SearchResults" and
// initializes it with default values for the "SearchType" and
// "SearchPatternString" attributes
func NewSearchResults(stype int, spattern string) *SearchResults {
	return &SearchResults{
		stype,
		spattern,
		make([]interface{}, 0),
		*regexp.MustCompile(spattern),
	}
}

// UpdateSearchResults appends a new value to the search results
// list if NOT ALREADY present in the "Results" list
func (resmap *SearchResults) UpdateSearchResults(val interface{}) {
	if !contains(resmap.Results, val) {
		resmap.Results = append(resmap.Results, val)
	}
}

// GetAbsDIRPathGivenRelativePath returns the absolute path on the file system given the
// relative path from where this function resides
func GetAbsDIRPathGivenRelativePath(relpath string) string {
	log.Debug()
	_, fname, _, _ := runtime.Caller(0)
	var path string
	if strings.HasPrefix(relpath, "../") {
		path = filepath.Join(filepath.Dir(fname), relpath)
	} else {
		path = relpath
	}
	return path
}

// GetSchemaStringWhenGivenFilePath generates a string that needs to
// be passed to the schema validator method when compiling a json schema
func GetSchemaStringWhenGivenFilePath(relativePathOfJSONSchemaFile string) string {
	log.Debug()
	_, fname, _, _ := runtime.Caller(0)
	var path string
	if strings.HasPrefix(relativePathOfJSONSchemaFile, "../") {
		path = filepath.Join(filepath.Dir(fname), relativePathOfJSONSchemaFile)
	} else {
		path = relativePathOfJSONSchemaFile
	}

	var schemaText = `{"$ref": "` + path + `"}`
	log.Debug(schemaText)
	return schemaText
}

// GetSchemaDefinitionFileAsJSONBuf reads a Schema file and returns JSON buf
func GetSchemaDefinitionFileAsJSONBuf(schemaFileName string) ([]byte, error) {
	log.Debug()
	bpath := GetAbsDIRPathGivenRelativePath(SchemaDir)
	yamlText, err := ioutil.ReadFile(bpath + "/" + schemaFileName)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(yamlText, &m)
	log.Debug(string(yamlText), err)

	return yamlText, err
}

// ValidateVnfdPostBody validates the given JSON body against the parameterized
// VNFD Input JSON schema "parameterizedVnfdInputSchema.json" for compliance
func ValidateVnfdPostBody(body []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaInputPath)
	ioReaderObj := strings.NewReader(schemaText)
	return ValidateJSONBufAgainstSchema(body, ioReaderObj, "vnfdPostBody.json")
}

// ValidateVnfdInstanceBody validates the given JSON body against the parameterized
// VNFD Instance JSON schema "parameterizedVnfdInstanceSchema.json" for compliance
func ValidateVnfdInstanceBody(jsonval []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaParameterizedInstanceRelPath)
	ioReaderObj := strings.NewReader(schemaText)
	return ValidateJSONBufAgainstSchema(jsonval, ioReaderObj, "vnfdInstanceBody.json")
}

// ValidatePaginatedVnfdsInstancesBody validates that JSON body returning the
// Vnfds adhere to the pagination format
func ValidatePaginatedVnfdsInstancesBody(jsonval []byte) error {
	log.Debug()
	var schemaText = GetSchemaStringWhenGivenFilePath(SchemaPaginatedInstancesRelPath)
	ioReaderObj := strings.NewReader(schemaText)
	return ValidateJSONBufAgainstSchema(jsonval, ioReaderObj, "vnfdsPaginatedInstancesBody.json")
}

// ValidateJSONBufAgainstSchema takes as arguments:
// i) a json buffer that needs to be validated against a schema
// ii) a io.Reader object that contains the schema definition information
// iii) a string (one of `schema.json` or `sch.json`) that represents if
// schema definition is in a file or in memory
// The function returns an error if the json buffer does not validate against
// the defined schema
func ValidateJSONBufAgainstSchema(jsonval []byte,
	schemaDefAsReaderObj io.Reader, url string) error {
	log.Debug()
	var m interface{}
	err := yaml.Unmarshal(jsonval, &m)
	if err != nil {
		log.WithFields(log.Fields{"UnMarshallError": err}).Error()
		return err
	}
	compiler := jsonschema.NewCompiler()
	//compiler.Draft = jsonschema.Draft4
	if err := compiler.AddResource(url, schemaDefAsReaderObj); err != nil {
		log.WithFields(log.Fields{"AddResourceError": err}).Error()
		return err
	}
	schema, err := compiler.Compile(url)
	if err != nil {
		log.WithFields(log.Fields{"CompileError": err}).Error()
		return err
	}

	if zerr := schema.ValidateInterface(m); zerr != nil {
		l := len(strings.Split(zerr.Error(), "\n"))
		log.WithFields(log.Fields{"SchemaValidateInterfaceError": zerr}).Error()
		return errors.New(strings.Split(zerr.Error(), "\n")[l-1])
	}
	return nil
}

// ValidateInputParamAgainstParameterizedVnfd validates the given "input_param"
// JSON file against the dynamically generated JSON Schema
func ValidateInputParamAgainstParameterizedVnfd(inputParamJSON []byte,
	parameterizedVnfdJSON []byte) error {
	log.Debug()
	inputParamDynSchema, e := GenerateJSONSchemaFromParameterizedTemplate(parameterizedVnfdJSON)
	if e != nil {
		return e
	}
	data, e := yaml.YAMLToJSON(inputParamDynSchema)
	if e != nil {
		return e
	}
	fmt.Println(string(data))

	return ValidateJSONBufAgainstSchema(inputParamJSON, strings.NewReader(string(data)), "inputParam.json")

}

// GetRegexMatchingListFromJSONBuff returns a list of strings that match
// a given pattern from the input json buffer
func GetRegexMatchingListFromJSONBuff(jsonval []byte, regexpattern string) []string {
	//re := regexp.MustCompile(`\$[A-Za-z][-A-Za-z0-9_]*`)
	re := regexp.MustCompile(regexpattern)
	res := re.FindAllString(string(jsonval), -1)

	return res
}

// CreateRevMapStructFromGivenStringListWithSpecifiedSeparator takes
// 		an array of strings
// 		a separator as argument, that separates string as key and value
//      a prefix trimmer, that trims the value of prefix such as '-'
// It creates a map with the RIGHT side of the separator as "KEY" and the LEFT side of the separator as "VALUE"
func CreateRevMapStructFromGivenStringListWithSpecifiedSeparator(
	slist []string,
	separator string,
	prefixToBeTrimmedFromVal string) map[string]interface{} {
	//var m map[string]interface{}
	var m = make(map[string]interface{})

	for _, elem := range slist {
		s := strings.Split(elem, separator)
		k := strings.TrimSpace(s[1])
		v := strings.TrimSpace(s[0])
		//m[k] = v
		m[k] = strings.TrimSpace(strings.TrimPrefix(v, prefixToBeTrimmedFromVal))
	}
	return m
}

// ParseMap iterates through a NESTED MAP and creates a MAP from the leaf KEY and VALUES
func (resmap *SearchResults) ParseMap(aMap map[string]interface{}) {
	for key, val := range aMap {
		var matchKeyFlag = false
		var matchValueFlag = false
		//if resmap.SearchKeyInput == key {
		if resmap.SearchType == MatchKey && resmap.re.MatchString(key) {
			matchKeyFlag = true
		} else if resmap.SearchType == MatchValue &&
			reflect.TypeOf(val).Kind() == reflect.String && resmap.re.MatchString(val.(string)) {
			matchValueFlag = true
		}
		//switch concreteVal := val.(type) {
		switch val.(type) {
		case map[string]interface{}:
			if matchKeyFlag {
				resmap.UpdateSearchResults(val.(map[string]interface{}))
			} else if matchValueFlag {
				resmap.UpdateSearchResults(key)
			}
			resmap.ParseMap(val.(map[string]interface{}))
		case []interface{}:
			if matchKeyFlag {
				resmap.UpdateSearchResults(val.([]interface{}))
			} else if matchValueFlag {
				resmap.UpdateSearchResults(key)
			}
			resmap.ParseArray(val.([]interface{}))
		case bool:
			if matchKeyFlag {
				resmap.UpdateSearchResults(val.(bool))
			} else if matchValueFlag {
				resmap.UpdateSearchResults(key)
			}
		case float64:
			if matchKeyFlag {
				resmap.UpdateSearchResults(val.(float64))
			} else if matchValueFlag {
				resmap.UpdateSearchResults(key)
			}
		default:
			if matchKeyFlag {
				resmap.UpdateSearchResults(val.(string))
			} else if matchValueFlag {
				resmap.UpdateSearchResults(key)
			}
		}
	}
}

// ParseArray iterates through an array
func (resmap *SearchResults) ParseArray(anArray []interface{}) {
	//for i, val := range anArray {
	for _, val := range anArray {
		//switch concreteVal := val.(type) {
		switch val.(type) {
		case map[string]interface{}:
			//fmt.Println("Index:", i)
			resmap.ParseMap(val.(map[string]interface{}))
		case []interface{}:
			//fmt.Println("Index:", i)
			resmap.ParseArray(val.([]interface{}))
			//default:
			//fmt.Println("Index", i, ":", concreteVal)
			//resmap.pvm[i] = val
		}
	}
}

// GenerateJSONSchemaFromParameterizedTemplate generated a dynamic schema
// by parsing the template for parameterized variables and looking up
// allowable values for those parameterized variables.
func GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON []byte) ([]byte, error) {
	// The regexp looks for the $ anywhere in the line and returns the entire line
	log.Debug()
	validRegexList := `.*\$.*`

	slist := GetRegexMatchingListFromJSONBuff(parameterizedJSON, validRegexList)
	log.WithFields(log.Fields{"RegexMatchingList": slist}).Debug()

	mapParameterizedParamAndDefinition := CreateRevMapStructFromGivenStringListWithSpecifiedSeparator(slist, ":", "-")
	log.WithFields(log.Fields{"mapParameterizedParamAndDefinition": mapParameterizedParamAndDefinition}).Debug()

	nonParamDefineJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(SchemaFileDefineNonParam)
	if err != nil {
		return nil, err
	}

	propjson := createSchemaForInputParamsFromParameterizedProperties(mapParameterizedParamAndDefinition, nonParamDefineJSONBuf)

	var src map[string]interface{}
	//_ = yaml.Unmarshal(propjson, &src)
	_ = json.Unmarshal(propjson, &src)

	////
	inputParamSchemaJSONBuf, err := GetSchemaDefinitionFileAsJSONBuf(SchemaFileInputParam)
	if err != nil {
		return nil, err
	}
	var inputParamSchemaMap map[string]interface{}
	//_ = yaml.Unmarshal(inputParamSchemaJSONBuf, &inputParamSchemaMap)
	_ = json.Unmarshal(inputParamSchemaJSONBuf, &inputParamSchemaMap)

	inter := mergemap.Merge(inputParamSchemaMap, src)

	reqjson := createSchemaForInputParamsWithRequiredSection(len(src), mapParameterizedParamAndDefinition)
	var req map[string]interface{}
	//_ = yaml.Unmarshal(reqjson, &req)
	_ = json.Unmarshal(reqjson, &req)

	final := mergemap.Merge(inter, req)

	//r, e := yaml.Marshal(final["inputParam"])
	r, e := json.Marshal(final["inputParam"])
	log.Debug(string(r), e)
	return r, e
}

// createSchemaForInputParamsWithRequiredSection takes as argument:
// i) reqCnt : number of keys to be added to the "required" section of the inputParams
// ii) a map that contains as its
//		key: the parameterized param from the parameterized template
//		value: the definition key that can be looked up in the json schema for allowable format and values
func createSchemaForInputParamsWithRequiredSection(reqCnt int, m map[string]interface{}) []byte {
	log.Debug()
	reqmap := make(map[string]map[string]interface{})
	reqmap[KeyInputParam] = make(map[string]interface{})
	reqmap[KeyInputParam][KeyRequired] = make([]string, reqCnt)

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k[1:]
		i++
	}
	keys = append(keys, KeyVnfdID, KeyName)
	reqmap[KeyInputParam][KeyRequired] = keys

	//reqjson, e := yaml.Marshal(reqmap)
	reqjson, e := json.Marshal(reqmap)

	if e != nil {
		panic(e)
	}
	return reqjson
}

// createSchemaForInputParamsFromParameterizedProperties takes as argument:
// i) a map that contains as its
// 		key: the parameterized param from the parameterized template
//		value: the definition key that can be looked up in the json schema for allowable format and values
// ii) json schema that contains property definitions and formats
// The function returns dynamically created Schema for the "input_param" as JSON buffer
func createSchemaForInputParamsFromParameterizedProperties(m map[string]interface{}, schemaJSON []byte) []byte {
	log.Debug()
	var schema map[string]interface{}
	//_ = yaml.Unmarshal(schemaJSON, &schema)
	_ = json.Unmarshal(schemaJSON, &schema)

	propmap := make(map[string]map[string]map[string]interface{})
	propmap[KeyInputParam] = make(map[string]map[string]interface{})
	propmap[KeyInputParam][KeyProperties] = make(map[string]interface{})

	for k, v := range m {
		pvm := NewSearchResults(MatchKey, v.(string))
		pvm.ParseMap(schema)
		for _, elem := range pvm.Results {
			switch elem.(type) {
			case map[string]interface{}:
				//k[1:] removes the first character '$'
				propmap[KeyInputParam][KeyProperties][k[1:]] = elem
			}
		}
	}

	//propjson, _ := yaml.Marshal(propmap)
	propjson, _ := json.Marshal(propmap)

	return propjson
}


