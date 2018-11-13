package TERMINUS

import (
	"strings"
	"fmt"
	"os/exec"
	"bytes"
	log "github.com/sirupsen/logrus"
	"time"
	"strconv"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"github.com/hpcloud/tail"
	"math/rand"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

