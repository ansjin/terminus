package TERMINUS

import "io/ioutil"

func CreateScriptFile(scriptPath string, script [] byte ) {
	err := ioutil.WriteFile(scriptPath, script, 0644)
	if err != nil {
		panic(err)
	}

}
func makeStartScript(accessKey, secret,region, clusterName,s3BucketName string  ) string {
	var startScript =	"#!/bin/sh\n"+
		"rm /root/.ssh/id_rsa\n"+
		"ssh-keygen -t rsa -f /root/.ssh/id_rsa -P \"\"\n"+
		"echo \"export AWS_ACCESS_KEY_ID="+ accessKey+"\">>~/.bashrc\n"+
		"echo \"export AWS_SECRET_ACCESS_KEY="+ secret+"\">>~/.bashrc\n"+
		"echo \"export AWS_REGION="+ region+"\">>~/.bashrc\n"+
		"echo \"export KOPS_CLUSTER_NAME="+ clusterName+"\">>~/.bashrc\n"+
		"echo \"export KOPS_STATE_STORE=s3://"+ s3BucketName+"\">>~/.bashrc\n" +

		"echo \"AWS_ACCESS_KEY_ID="+ accessKey+"\">>~/.profile\n"+
		"echo \"AWS_SECRET_ACCESS_KEY="+ secret+"\">>~/.profile\n"+
		"echo \"AWS_REGION="+ region+"\">>~/.profile\n"+
		"echo \"KOPS_CLUSTER_NAME="+ clusterName+"\">>~/.profile\n"+
		"echo \"KOPS_STATE_STORE=s3://"+ s3BucketName+"\">>~/.profile\n" +

		"echo \"AWS_ACCESS_KEY_ID="+ accessKey+"\">>/etc/environment\n"+
		"echo \"AWS_SECRET_ACCESS_KEY="+ secret+"\">>/etc/environmentn"+
		"echo \"AWS_REGION="+ region+"\">>/etc/environment\n"+
		"echo \"KOPS_CLUSTER_NAME="+ clusterName+"\">>/etc/environment\n"+
		"echo \"KOPS_STATE_STORE=s3://"+ s3BucketName+"\">>/etc/environment\n" +
		"source ~/.bashrc\n"+
		"source ~/.profile\n"+
		"aws s3api create-bucket --bucket "+s3BucketName+" --region "+region+" --create-bucket-configuration LocationConstraint="+region+"\n"+
		"aws s3api put-bucket-versioning --bucket "+s3BucketName+" --versioning-configuration Status=Enabled\n"
		return startScript
}
func makeCreateClusterScript(nodeCount, nodeSize,masterSize, zone string  ) string {
	var createClusterScript = "#!/bin/sh\n" +
		"source ~/.bashrc\n" +
		"kops create cluster --node-count=" + nodeCount + " --node-size=" + nodeSize + " --master-size=" + masterSize + " --zones=" + zone + " --name=${KOPS_CLUSTER_NAME}"
	return createClusterScript
}

func makeUpdateClusterScript() string {
	var updateClusterScript =	"#!/bin/sh\n"+
		"source ~/.bashrc\n"+
		"kops update cluster --name ${KOPS_CLUSTER_NAME} --yes"
	return updateClusterScript
}

func makeValidateClusterScript() string {
	var validateClusterScript =	"#!/bin/sh\n"+
		"source ~/.bashrc\n"+
		"kops validate cluster --name ${KOPS_CLUSTER_NAME}"
	return validateClusterScript
}

func makeGetPasswdScriptt() string {
	var getPasswdScript =	"#!/bin/sh\n"+
		"source ~/.bashrc\n"+
		"kops get secrets kube --type secret -oplaintext"
	return getPasswdScript
}

func makeGetTokenScript() string {
	var getTokenScript =	"#!/bin/sh\n"+
		"source ~/.bashrc\n"+
		"kops get secrets admin --type secret -oplaintext"
	return getTokenScript
}

func makeDeleteClusterScript() string {
	var deleteClusterScript =	"#!/bin/sh\n"+
		"source ~/.bashrc\n"+
		"kops delete cluster --name ${KOPS_CLUSTER_NAME} --yes"
	return deleteClusterScript
}

