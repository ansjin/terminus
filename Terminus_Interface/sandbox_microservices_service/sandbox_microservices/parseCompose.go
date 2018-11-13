package sandbox_microservices

import (
	"reflect"
	"strings"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/ghodss/yaml"
	"encoding/json"
)

func setField(v *ServicesStruct, field string, val interface{}) {
	r := reflect.ValueOf(v)
	field = strings.Title(field)
	log.Info(field)
	s := r.Elem()
	f := s.FieldByName(field)
	log.Info(f)

	if f.IsValid() {
		// A Value can be changed only if it is
		// addressable and was not obtained by
		// the use of unexported struct fields.
		if f.CanSet() {
			// change value of N
			if f.Kind() == reflect.String {

				x := string(val.(string))
				log.Info(x)
				f.SetString(x)
			}
		}
	}
}
func setFieldSlice(v *ServicesStruct, field string, val []interface{})  {
	r := reflect.ValueOf(v)
	field = strings.Title(field)
	log.Info(field)
	s := r.Elem()
	f := s.FieldByName(field)
	log.Info(f)
	if f.IsValid() {
		// A Value can be changed only if it is
		// addressable and was not obtained by
		// the use of unexported struct fields.
		if f.CanSet() {
			// change value of N
			log.Info(f.Type())
			if f.Kind() == reflect.Slice {
				f.Set(reflect.MakeSlice(f.Type(), len(val), len(val)))
				for i := 0; i < len(val); i++ {
					gg := reflect.ValueOf(val)
					gotVal := gg.Index(i)
					stringVal := gotVal.Interface().(string)
					f.Index(i).SetString(stringVal)
				}
			}
		}
	}
}
func setFieldMap(v *ServicesStruct, field string, val map[string]interface {})  {
	r := reflect.ValueOf(v)
	field = strings.Title(field)
	log.Info(field)
	s := r.Elem()
	f := s.FieldByName(field)
	log.Info(f)
	if f.IsValid() {
		// A Value can be changed only if it is
		// addressable and was not obtained by
		// the use of unexported struct fields.
		if f.CanSet() {
			// change value of N
			log.Info(f.Type())
			if f.Kind() == reflect.Map {
				f.Set(reflect.MakeMap(f.Type()))
				for key, value :=range val{

					ggKey := reflect.ValueOf(key)
					ggVal := reflect.ValueOf(value)
					f.SetMapIndex(ggKey,ggVal)
				}
			}
		}
	}
}
func setFieldNil(v *ServicesStruct, field string)  {
	r := reflect.ValueOf(v)
	field = strings.Title(field)
	log.Info(field)
	s := r.Elem()
	f := s.FieldByName(field)
	log.Info(f)
	if f.IsValid() {
		// A Value can be changed only if it is
		// addressable and was not obtained by
		// the use of unexported struct fields.
		if f.CanSet() {
			// change value of N
			log.Info(f.Type())
			if f.Kind() == reflect.Slice {
				f.Set(reflect.MakeSlice(f.Type(), 1, 1))
				f.Index(0).SetString("")
			}
		}
	}
}

func parseyamlToJson(filePath string)  DockerComposeParsedStruct{
	dat, err := ioutil.ReadFile(filePath)
	if(!checkError(err)){
		log.Info(string(dat))
		dockerComposeJson, err := yaml.YAMLToJSON(dat)
		if(!checkError(err)){

			var parsedObj DockerComposeParsedStructTemp
			var finalParsedObj DockerComposeParsedStruct
			var RelevantServices 		= make( map[string]ServicesStruct)
			if err := json.Unmarshal([]byte(dockerComposeJson), &parsedObj); err != nil {
				panic(err)
			}

			if err := json.Unmarshal([]byte(dockerComposeJson), &parsedObj.Services); err != nil {
				panic(err)
			}
			delete(parsedObj.Services, "version")

			type YamlData struct {
				Key string;
				Value interface {}
			}
			var yamlData []YamlData
			//var servicesInfo []ServicesStruct
			for _, value := range parsedObj.Services {

				for serviceName, serviceInfoAll := range value.(map[string]interface {}) {

					yamlData  = append(yamlData, YamlData{serviceName, serviceInfoAll})
					var serviceInfo ServicesStruct
					for serviceInfoKey, serviceInfoValue := range serviceInfoAll.(map[string]interface {}) {
						field :=strings.Title(serviceInfoKey)
						log.Info("FieldName:", field)
						switch v := serviceInfoValue.(type) {
						case int:
							// v is an int here, so e.g. v + 1 is possible.
							log.Info("Integer: %v", v)
						case float64:
							// v is a float64 here, so e.g. v + 1.0 is possible.
							log.Info("Float64: %v", v)
						case string:
							// v is a string here, so e.g. v + " Yeah!" is possible.
							log.Info("String:", v)
							setField(&serviceInfo, serviceInfoKey, serviceInfoValue)
						case []interface{}:
							log.Info("[]interface:",v)
							if(serviceInfoValue==nil){
								//ignore it
							}else {
								setFieldSlice(&serviceInfo, serviceInfoKey, serviceInfoValue.([]interface {}))
							}
						case map[string]interface {}:
							log.Info("map[string]interface {}",v )
							if(serviceInfoValue==nil){
								//ignore it
							}else{
								log.Info("serviceInfoValue:", serviceInfoValue)
								setFieldMap(&serviceInfo, serviceInfoKey, serviceInfoValue.(map[string]interface {}))
							}
						default:
							// And here I'm feeling dumb. ;)
							log.Info("I don't know interface type, ask stackoverflow.")
						}
					}
					//servicesInfo  = append(servicesInfo, serviceInfo)
					RelevantServices[serviceName] = serviceInfo
				}
			}
			finalParsedObj.Version  = parsedObj.Version
			finalParsedObj.Services  = RelevantServices

			return finalParsedObj
		}
	}
	return DockerComposeParsedStruct{}
}