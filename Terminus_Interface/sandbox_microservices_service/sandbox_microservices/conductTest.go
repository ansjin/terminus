package sandbox_microservices

import (
	b64 "encoding/base64"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
	log "github.com/sirupsen/logrus"
)

func GetVMStartScript(dockerComposeFileContentsJson string)string{
	var VMStartScript = "#!bin/sh \n"+
		"echo \"setup\"  \n"+
		"sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual  \n"+
		"sudo apt-get update  \n"+
		"sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common  \n"+
		"#Microservice Run  \n"+
		"apt-get remove -y docker docker-engine docker.io  \n"+
		"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -  \n"+
		"add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"  \n"+
		"apt-get update  \n"+
		"apt-get install -y docker-ce  \n"+
		"service docker start  \n"+
		"curl -L https://github.com/docker/compose/releases/download/1.13.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose  \n"+
		"chmod +x /usr/local/bin/docker-compose  \n"+
		"FILE_DOCKER_COMPOSE=\"./docker-compose.json\"  \n"+
		"/bin/cat <<EOM >$FILE_DOCKER_COMPOSE  \n"+
		dockerComposeFileContentsJson+"  \n"+
		"EOM\n"+
		"sudo docker-compose -f docker-compose.json up&  \n"
	encodedString:=b64.StdEncoding.EncodeToString([]byte(VMStartScript))
	return encodedString
}
func startTestVM( dockerComposeFileContentsJson string)  string {

	sessionAWS := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc := ec2.New(sessionAWS)
	var allInstancesStarted []Ec2Instances

	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(20),
				},
			},
		},
		ImageId:      aws.String(GetImageFromRegion(AWSConfig.Region)),
		InstanceType: aws.String("t2.large"),
		KeyName:      aws.String(AWSConfig.KeyPairName),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String(AWSConfig.SecurityGroup),
		},
		SubnetId: aws.String(AWSConfig.SubnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Purpose"),
						Value: aws.String("sand_boxing"),
					},
				},
			},
		},
		UserData: aws.String(GetVMStartScript(dockerComposeFileContentsJson)),
	}
	result, err := svc.RunInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error(aerr.Error())
		}
		return ""
	}

	for _, instance := range result.Instances {

		oneInstance := Ec2Instances{InstanceId: ValueAssignString(instance.InstanceId, ""),
			ImageId: ValueAssignString(instance.ImageId, ""),
			InstanceType: ValueAssignString(instance.InstanceType, ""),
			LaunchTime: *instance.LaunchTime,
			InstanceState: ValueAssignString(instance.State.Name, ""),
			AvailabilityZone: ValueAssignString(instance.Placement.AvailabilityZone, ""),
			CoreCount: ValueAssignInt64(instance.CpuOptions.CoreCount, 0),
			PublicIpAddress: ValueAssignString(instance.PublicIpAddress, "")}

		allInstancesStarted = append(allInstancesStarted, oneInstance)
	}
	log.Info(allInstancesStarted)
	return allInstancesStarted[0].InstanceId
}

func getVMPublicIP(startedInstanceId string)  string{
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc2 := ec2.New(session)

	var allInstances []Ec2Instances

	input2 := ec2.DescribeInstancesInput{InstanceIds: []*string{
		aws.String(startedInstanceId),
	},}
	result2, er2r := svc2.DescribeInstances(&input2)
	if er2r != nil {
		log.Error(er2r)
		return ""
	}
	for _, reservation := range result2.Reservations {
		for _, instance := range reservation.Instances {

			oneInstance := Ec2Instances{InstanceId: ValueAssignString(instance.InstanceId, ""),
				ImageId: ValueAssignString(instance.ImageId, ""),
				InstanceType: ValueAssignString(instance.InstanceType, ""),
				LaunchTime: *instance.LaunchTime,
				InstanceState: ValueAssignString(instance.State.Name, ""),
				AvailabilityZone: ValueAssignString(instance.Placement.AvailabilityZone, ""),
				CoreCount: ValueAssignInt64(instance.CpuOptions.CoreCount, 0),
				PublicIpAddress: ValueAssignString(instance.PublicIpAddress, "")}

			allInstances = append(allInstances, oneInstance)
		}
	}
	log.Info(allInstances[0].PublicIpAddress)
	return allInstances[0].PublicIpAddress
}

func terminateTestVM(instanceId string) {

	session := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AWSConfig.AwsAccessKeyId, AWSConfig.AwsSecretAccessKey, ""),
		Region:      aws.String(AWSConfig.Region),
	}))

	svc2 := ec2.New(session)

	var allInstances []string

	input2 := ec2.TerminateInstancesInput{InstanceIds: []*string{
		aws.String(instanceId),
	},}
	result2, er2r := svc2.TerminateInstances(&input2)
	if er2r != nil {
		log.Error(er2r)
		return
	}
	for _, instance := range result2.TerminatingInstances {
		allInstances = append(allInstances, ValueAssignString(instance.InstanceId, ""))
	}

	log.Info("Terminate Instances with id: ", allInstances)

}