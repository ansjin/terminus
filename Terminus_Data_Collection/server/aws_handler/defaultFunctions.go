package TERMINUS

import (
	"strings"
	"time"
	"strconv"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/hpcloud/tail"
	"fmt"
	"os"
)

func InitUserConfigurationFirst(){

	AllInstanceTypes = []string{ "100","200","500","700", "t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large","t2.xlarge", "t2.2xlarge","m5.xlarge"}

	DefaultRegion = []string{ "us-east-2", "us-east-1", "us-west-1", "eu-central-1", "ap-south-1", "eu-west-1", "ap-southeast-2"}

	DefaultZone = []string{ "us-east-2a","us-east-1a","us-west-1a","eu-central-1a","ap-south-1a","eu-west-1a", "ap-southeast-2a"}

	DefaultAMI = []string{ "ami-5e8bb23b", "ami-759bc50a", "ami-4aa04129", "ami-de8fb135", "ami-188fba77","ami-2a7d75c0", "ami-47c21a25"}

	B = &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
	UserConfig =  KopsConfig{AwsAccessKeyId: os.Getenv("AWS_KEY"), AwsSecretAccessKey: os.Getenv("AWS_SECRET"),
		Region:  os.Getenv("AWS_DEFAULT_REGION"), ContainerName:os.Getenv("KOPS_CONTAINER_NAME"), KubeClusterName:os.Getenv("KUBE_CLUSTER_NAME"),
		S3BucketName:strings.ToLower(os.Getenv("KOPS_S3_BUCKET_NAME")), Zone:os.Getenv("AWS_DEFAULT_ZONE"),
		NodeSize:"t2.medium", NodeCount:"2", MasterSize:"t2.large"}

	AllScriptsNames = ScriptTypeNames{"initClusterScript.sh","createClusterScript.sh",
									  "updateClusterScript.sh","validateClusterScript.sh",
									  "getPasswdScript.sh", "getTokenScript.sh",
									  "deleteClusterScript.sh" }

	SupportedNumMicroserviceReplicas = []string{"1", "2", "3","4", "5", "10"}
	SupportedNumMicroservicesApp = []string{"1", "4"}
	DefaultAppNames = []string{"primeapp", "movieapp","webacapp", "mixalapp", "serveapp"}
	SupportedTypesOfMicroservices = []string{"compute", "dbaccess", "web", "mix", "sandbox"}
}

func Schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func ValueAssignString(value *string, fallback string) string{

	if value!=nil {
		return *value
	} else {
		return fallback
	}
}
func ValueAssignInt64(value *int64, fallback int64) int64{

	if value!=nil {
		return *value
	} else {
		return fallback
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
func getKeyFile(directoyPath string) (key ssh.Signer, err error) {
	file := directoyPath
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error(err)
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		log.Error(err)
	}
	return
}
func ReadLogsContinously(){
	t, err := tail.TailFile("/data/logs", tail.Config{Follow: true})
	if err != nil {
		fmt.Printf("%s", err)
	}
	for line := range t.Lines {
		B.messages<-line.Text
	}
}
