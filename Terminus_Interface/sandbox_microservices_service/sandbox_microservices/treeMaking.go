package sandbox_microservices

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

type Node struct {
	ServiceName     string
	ServiceInfo ServicesStruct
	Links []*Node
}

func add(serviceName, parentName string, serviceInfo ServicesStruct,  ServiceTable map[string]*Node, root *Node) {
	log.Info("add: serviceName=%v parentName=%v\n", serviceName, parentName)

	node := &Node{ServiceName: serviceName, ServiceInfo:serviceInfo, Links: []*Node{}}

	if parentName == "0" {
		*root=*node
		node = root
	} else {

		parent, ok := ServiceTable[parentName]
		if !ok {
			log.Info("add: parentName=%v: not found\n", parentName)
			return
		}
		parent.Links = append(parent.Links, node)
	}
	ServiceTable[serviceName] = node
}

func addAllServicesToTree(serviceName, parent string, parseddockerComposeObj DockerComposeParsedStruct, ServiceTable map[string]*Node, root *Node) {

	mainServiceObj:= parseddockerComposeObj.Services[serviceName]

	add(serviceName, parent,mainServiceObj, ServiceTable, root)


	if(len(mainServiceObj.Links)>0){

		for i:=0; i< len(mainServiceObj.Links); i++{
			addAllServicesToTree(mainServiceObj.Links[i], serviceName, parseddockerComposeObj, ServiceTable, root)
		}

	}else if(len(mainServiceObj.Depends_on)>0){
		for i:=0; i< len(mainServiceObj.Links); i++{
			addAllServicesToTree(mainServiceObj.Depends_on[i], serviceName, parseddockerComposeObj, ServiceTable, root)
		}
	}else {

		return
	}
}

func showNode(node *Node, prefix string) {
	if prefix == "" {
		log.Info( node.ServiceName)
	} else {
		log.Info(prefix, node.ServiceName)
	}
	for _, n := range node.Links {
		showNode(n, prefix+"--")
	}
}

func show(root *Node) {
	if root == nil {
		log.Info("show: root node not found\n")
		return
	}
	log.Info("RESULT:")
	showNode(root, "")
}

func searchForToRelevantDepndentServices(serviceName string) bool{

	relevantStrings:=[]string{"mongo", "influx", "setup", "log", "elastic", "kibana", "grafana", "db", "kafka", "flink", "maria", "redis"}

	for _, b := range relevantStrings {
		if strings.Contains(serviceName, b) 	{
			return true
		}
	}
	return false

}

func searchAndRemoveDependentServices(serviceName string, parseddockerComposeObj *DockerComposeParsedStruct){

	for i:=0; i<len(parseddockerComposeObj.Services[serviceName].Links);i++{

		if(searchForToRelevantDepndentServices(parseddockerComposeObj.Services[serviceName].Links[i])==false){
			if(len(parseddockerComposeObj.Services[serviceName].Links) == 1){
				var x = parseddockerComposeObj.Services[serviceName]
				x.Links = []string{}
				parseddockerComposeObj.Services[serviceName] = x
			}else{
				var x = parseddockerComposeObj.Services[serviceName]
				x.Links[i] = x.Links[len(x.Links)-1]
				x.Links = x.Links [:len(x.Links )-1]
				parseddockerComposeObj.Services[serviceName] = x
				i--
			}
		}
	}
	for i:=0; i<len(parseddockerComposeObj.Services[serviceName].Depends_on);i++{
		if(searchForToRelevantDepndentServices(parseddockerComposeObj.Services[serviceName].Depends_on[i])==false){
			if(len(parseddockerComposeObj.Services[serviceName].Depends_on) == 1){
				var x = parseddockerComposeObj.Services[serviceName]
				x.Depends_on = []string{}
				parseddockerComposeObj.Services[serviceName] = x
			}else{
				var x = parseddockerComposeObj.Services[serviceName]
				x.Depends_on[i] = x.Depends_on[len(x.Depends_on)-1]
				x.Depends_on = x.Depends_on [:len(x.Depends_on )-1]
				parseddockerComposeObj.Services[serviceName] = x
				i--

			}
		}
	}
}

func searchAndReplaceDependentServices(serviceName string, parseddockerComposeObj *DockerComposeParsedStruct){

	for i:=0; i<len(parseddockerComposeObj.Services[serviceName].Links);i++{

		if(searchForToRelevantDepndentServices(parseddockerComposeObj.Services[serviceName].Links[i])==false){
			parseddockerComposeObj.Services[serviceName].Links[i] = "dummy_" + parseddockerComposeObj.Services[serviceName].Links[i]
		}
	}
	for i:=0; i<len(parseddockerComposeObj.Services[serviceName].Depends_on);i++{
		if(searchForToRelevantDepndentServices(parseddockerComposeObj.Services[serviceName].Depends_on[i])==false){
			parseddockerComposeObj.Services[serviceName].Depends_on[i] = "dummy_" + parseddockerComposeObj.Services[serviceName].Depends_on[i]
		}
	}
}
func formYamlNode(node *Node, prefix string, parseddockerComposeObj *DockerComposeParsedStruct) {
	if prefix == "" {
		log.Info( node.ServiceName)
		searchAndReplaceDependentServices(node.ServiceName, parseddockerComposeObj)
		log.Info( parseddockerComposeObj.Services[node.ServiceName])

	} else {
		log.Info( prefix, node.ServiceName)
		searchAndReplaceDependentServices(node.ServiceName, parseddockerComposeObj)
		log.Info(parseddockerComposeObj.Services[node.ServiceName])
	}
	for _, n := range node.Links {
		formYamlNode(n, prefix+"--", parseddockerComposeObj)
	}
}

func formYaml(root *Node, parseddockerComposeObj DockerComposeParsedStruct) DockerComposeParsedStruct{
	if root == nil {
		log.Info("formYaml: root node not found\n")
		return DockerComposeParsedStruct{}
	}
	log.Info("RESULT:\n")
	formYamlNode(root, "",&parseddockerComposeObj)
	return parseddockerComposeObj
}
