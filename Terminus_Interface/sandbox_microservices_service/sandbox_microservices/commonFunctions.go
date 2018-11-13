package sandbox_microservices

import (
	"strings"
	"fmt"
	"os/exec"
	"bytes"
	log "github.com/sirupsen/logrus"
	"time"
	"strconv"
	"github.com/hpcloud/tail"
	"math/rand"
	"os"
)
func InitUserConfigurationFirst(){

	AllInstanceTypes = []string{ "100","200","500","700", "t2.micro", "t2.nano", "t2.small", "t2.medium", "t2.large","t2.xlarge"}

	DefaultRegion = []string{ "us-east-2", "us-east-1", "us-west-1", "eu-central-1", "ap-south-1", "eu-west-1", "ap-southeast-2"}

	DefaultZone = []string{ "us-east-2a","us-east-1a","us-west-1a","eu-central-1a","ap-south-1a","eu-west-1a", "ap-southeast-2a"}

	DefaultAMI = []string{ "ami-5e8bb23b", "ami-759bc50a", "ami-4aa04129", "ami-de8fb135", "ami-188fba77","ami-2a7d75c0", "ami-47c21a25"}


	B = &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
	AWSConfig =  AWSConfigStruct{AwsAccessKeyId: os.Getenv("AWS_KEY"), AwsSecretAccessKey: os.Getenv("AWS_SECRET"),
		Region:  os.Getenv("AWS_DEFAULT_REGION"), KeyPairName:os.Getenv("AWS_KEY_PAIR_NAME"), SubnetId:os.Getenv("AWS_SUBNET_ID"),
		SecurityGroup:os.Getenv("AWS_SECURITY_GROUP")}

}
func checkError(e error) bool{
	if e != nil {
		log.Error(e)
		return true
	}else {
		return false
	}
}

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
func StringToFloat(stringVal string) float64 {
	// to convert a float number to a string
	if s, err := strconv.ParseFloat(stringVal, 64); err == nil {
		return s// 3.14159265
	}
	return 0
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
func init() {
	rand.Seed(time.Now().UnixNano())
}
//function to get the public ip address
func GetOutboundIP() string {

	cmd:="dig +short myip.opendns.com @resolver1.opendns.com"
	wanip:=exe_cmd_output(cmd)
	wanip = strings.TrimSuffix(wanip, "\n")
	return string(wanip)
}

func GetImageFromRegion(region string) string{
	for index, b := range DefaultRegion {
		if b == region {
			return DefaultAMI[index]
		}
	}
	return ""
}
func GetZoneFromRegion(region string) string{
	for index, b := range DefaultRegion {
		if b == region {
			return DefaultZone[index]
		}
	}
	return ""
}