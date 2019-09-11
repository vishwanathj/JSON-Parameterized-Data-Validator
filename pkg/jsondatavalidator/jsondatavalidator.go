package jsondatavalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/peterbourgon/mergemap"
	log "github.com/sirupsen/logrus"

	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema"
)

func init() {
	log.Debug()
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
		return fmt.Errorf("UnMarshallError")
	}
	compiler := jsonschema.NewCompiler()
	//compiler.Draft = jsonschema.Draft4
	if err := compiler.AddResource(url, schemaDefAsReaderObj); err != nil {
		log.WithFields(log.Fields{"AddResourceError": err}).Error()
		return fmt.Errorf("AddResourceError")
	}
	schema, err := compiler.Compile(url)
	if err != nil {
		log.WithFields(log.Fields{"CompileError": err}).Error()
		return fmt.Errorf("CompilerError")
	}

	if zerr := schema.ValidateInterface(m); zerr != nil {
		l := len(strings.Split(zerr.Error(), "\n"))
		log.WithFields(log.Fields{"SchemaValidateInterfaceError": zerr}).Error()
		return errors.New(strings.Split(zerr.Error(), "\n")[l-1])
	}
	return nil
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
func GenerateJSONSchemaFromParameterizedTemplate(parameterizedJSON []byte,
	nonParamDefineJSONBuf []byte, inputParamSchemaJSONBuf []byte,
	keysToAddToRequiredSection []string, regExpStr string) ([]byte, error) {
	// The regexp looks for the $ anywhere in the line and returns the entire line
	log.Debug()
	//validRegexList := `.*\$.*`

	slist := GetRegexMatchingListFromJSONBuff(parameterizedJSON, regExpStr)
	log.WithFields(log.Fields{"RegexMatchingList": slist}).Debug()

	mapParameterizedParamAndDefinition := CreateRevMapStructFromGivenStringListWithSpecifiedSeparator(slist, ":", "-")
	log.WithFields(log.Fields{"mapParameterizedParamAndDefinition": mapParameterizedParamAndDefinition}).Debug()

	propjson := createSchemaForInputParamsFromParameterizedProperties(
		mapParameterizedParamAndDefinition,
		nonParamDefineJSONBuf)

	var src map[string]interface{}
	_ = json.Unmarshal(propjson, &src)

	var inputParamSchemaMap map[string]interface{}
	_ = json.Unmarshal(inputParamSchemaJSONBuf, &inputParamSchemaMap)

	inter := mergemap.Merge(inputParamSchemaMap, src)

	reqjson := createSchemaForInputParamsWithRequiredSection(len(src),
		mapParameterizedParamAndDefinition, keysToAddToRequiredSection)
	var req map[string]interface{}
	_ = json.Unmarshal(reqjson, &req)

	final := mergemap.Merge(inter, req)

	r, e := json.Marshal(final["inputParam"])
	log.Debug(string(r), e)
	return r, e
}

// createSchemaForInputParamsWithRequiredSection takes as argument:
// i) reqCnt : number of keys to be added to the "required" section of the inputParams
// ii) a map that contains as its
//		key: the parameterized param from the parameterized template
//		value: the definition key that can be looked up in the json schema for allowable format and values
func createSchemaForInputParamsWithRequiredSection(reqCnt int,
	m map[string]interface{}, keysToAddToRequiredSection []string) []byte {
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

	keys = append(keys, keysToAddToRequiredSection...)
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
