package TERMINUS

import (
	"fmt"
	"strings"
	"os/exec"
	"bytes"
	log "github.com/sirupsen/logrus"
)
func exe_cmd_output(cmd string) string {
	log.Info("Command : ",cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	fmt.Println(parts)
	cmdc :=exec.Command(head, parts...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdc.Stdout = &out
	cmdc.Stderr = &stderr
	err := cmdc.Run()
	if err != nil {
		log.Error(stderr.String())
	}
	substrParts := strings.Split(out.String(), "\n")
	for i:=0;i<len(substrParts);i++{
		log.Info(substrParts[i])
	}
	return out.String()
}
func exe_cmd(cmd string) {
	exe_cmd_output(cmd)
	return
}

func MakeDirectory(dirName string )  {

	// create a directoy here
	command := "mkdir /data/" + dirName
	exe_cmd(command)
	return
}
func RemoveDirectory(dirName string )  {
	// create a directoy here
	command := "rm -r /data/" + dirName
	exe_cmd(command)

	return
}

func runScriptWithoutReturningOutput(dirName string, scriptName string)  {
	// run a script
	command := "sh /data/"+ dirName+"/"+scriptName

	exe_cmd(command)

	return

}
func runScriptWithOutput(dirName string, scriptName string) string {
	// run a script
	command := "sh /data/"+ dirName+"/"+scriptName

	output:=exe_cmd_output(command)

	return output

}

func runKopsHelpCommand()  {

	command := "kops --help"

	exe_cmd(command)

	return
}
func runKubectlApplyCommand(yamlScriptPath string)  {

	command := "kubectl apply -f "+yamlScriptPath
	exe_cmd(command)
	return
}
func runKubectlDeleteCommand(yamlScriptPath string)  {

	command := "kubectl delete -f "+yamlScriptPath
	exe_cmd(command)
	return
}
func runKubectlClusterInfoCommand() string {

	command := "kubectl cluster-info"
	output:=exe_cmd_output(command)
	return output
}
func runKubectlGetPodsCommand() string {

	command := "kubectl get pods --all-namespaces"
	output:=exe_cmd_output(command)
	return output
}
func runKubectlGetServicesCommand() string {

	command := "kubectl get services --all-namespaces"
	output:=exe_cmd_output(command)
	return output
}

func runKubectlGetExternalIpSvcCommand(appName string) []string {

	command := "kubectl get svc " + appName + "  -o wide"
	output:=exe_cmd_output(command)

	s := strings.Split(output, "\n")
	s_sub := strings.Split(s[1], " ")

	return s_sub
}

func CopyDataToDirectory(dirName string )  {

	// create a directoy here
	command := "mkdir /data/" + dirName
	exe_cmd(command)

	command_mkdir_Mongo := "mkdir /data/" + dirName+"/mongo"
	exe_cmd(command_mkdir_Mongo)

	command_mkdir_influx:= "mkdir /data/" + dirName+"/influx"
	exe_cmd(command_mkdir_influx)

	copy_command_mongo := "cp -r /data/db /data/" + dirName+"/mongo"
	exe_cmd(copy_command_mongo)

	copy_command_influx := "cp -r /srv/docker/influxdb/data /data/" + dirName+"/influx"
	exe_cmd(copy_command_influx)
	return
}
func RemoveDirectoryMongoInflux( )  {

	// create a directoy here
	command := "rm -r /data/db"
	exe_cmd(command)

	commandRmInflux := "rm -r /srv/docker/influxdb/data"
	exe_cmd(commandRmInflux)

	return
}