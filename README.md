# Terminus [![Docker Status](https://github.com/ansjin/terminus/blob/master/Docs/docker-hub.jpg)](https://hub.docker.com/u/terminusimages/dashboard/)
Terminus is a microservice’s performance modeling and sand-boxing tool.
Terminus is derived from the word "terminal" which means end point of something. Terminus consists of components implementing microservices architecture; these components can be scaled up and down on demand.
The tool automates the setup of a Kubernetes cluster and deploys the monitoring services and a load generator. It was developed using Golang and Python. Terminus comprises both the API and the user interface allowing the user to easily interact with the tool and conduct comprehensive performance modeling for an arbitrary microservice application. 

## Overall Architecture
The architecture of the tool and the communication paths between its components in a typical use-case are shown in below figure: 
<p align="center">
<img src="https://github.com/ansjin/terminus/blob/master/Docs/architecture.png"></img>
</p> 

## Setup
### Environment setup

 Currently, only AWS cloud is supported. So this tool can only be run on AWS cloud.
 So you can either use AWS CLI or AWS console for initial tasks.
#### Using AWS CLI 

##### 1. Installation of AWS CLI 
https://aws.amazon.com/cli/

 On MacOS, Windows and Linux OS:
 
 The officially supported way of installing the tool is with `pip`:
 
```bash
pip install awscli
```

###### MacOS

You can grab the tool with homebrew, although this is not officially supported by AWS.
```bash
brew update && brew install awscli
```

###### Windows

You can download the MSI installer from this page and follow the steps through the installer which requires no other dependencies: 
https://docs.aws.amazon.com/cli/latest/userguide/awscli-install-windows.html
 
##### 2. Setup IAM user 
 1. In order to build clusters within AWS we'll create a dedicated IAM user for
    `kops`.  This user requires API credentials in order to use `kops`. 
 2. The `kops` user will require the following IAM permissions to function properly:
    
    ```
    AmazonEC2FullAccess
    AmazonRoute53FullAccess
    AmazonS3FullAccess
    IAMFullAccess
    AmazonVPCFullAccess
    ```
    
    You can create the kops IAM user from the command line using the following (if AWS CLI is installed):
    
    ```bash
    aws iam create-group --group-name kops
    
    aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops
    aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops
    aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops
    aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops
    aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops
    
    aws iam create-user --user-name kops
    
    aws iam add-user-to-group --user-name kops --group-name kops
    
    aws iam create-access-key --user-name kops
    ```
    
    You should record the SecretAccessKey and AccessKeyID in the returned JSON
    output,
 
 #### Using AWS Console
 ##### 1. Setup IAM user 
 Create the user, and credentials, using the [AWS console](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html#id_users_create_console).
 
 The `kops` user will require the following IAM permissions to function properly:
     
     ```
     AmazonEC2FullAccess
     AmazonRoute53FullAccess
     AmazonS3FullAccess
     IAMFullAccess
     AmazonVPCFullAccess
     ```
 Add the above permission to the user account and note down the <b> SecretAccessKey and AccessKeyID</b>.
 ##### 2. Create Security Key pair
 1. Go into the Services->EC2->Network&Security->KeyPairs
 2. Create a new key pair and note down its <b>name</b>.
 
 ##### 3. Create Security Group
 1. Go into the Services->EC2->Network&Security->Security Groups
 2. Create a new security Group with allow all traffics from any source (for the time being) for both inbound and outbound rules.
    Note down its <b>ID</b>.
    
##### 3. Subnet ID
  1. Go into Services->VPC->Subnets
  2. Note down the default <b>Subnet ID</b> 

### VM creation for Terminus Deployment
 1. Create a t2.large with OS as Ubuntu 18.04 VM on AWS.
 2. Add storage disk of minimum 90GB.
 3. Add the security group with allowed access to all traffic from any source to it.   
 4. SSH into it
### Required Installation on the VM
 1. Mongodb needs to installed on the machine from which this restoreing is done. <a href = "https://docs.mongodb.com/manual/installation/">Here</a> is the procedure to install or follwo these command for Ubuntu 18.04 : 
 ```sh
   $ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
   $ echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.0.list
   $ sudo apt-get update
   $ sudo apt-get install -y mongodb-org
   ```
 2. Install s3cmd 
 ```sh
    $ sudo apt-get install -y s3cmd
   ```
 3. Configure s3cmd
 ```sh
    $ s3cmd --configure
   ``` 
 4. Provide your Access key and Secret Key. Default Region as ap-southeast-2
    Leave the rest to default and save the configuration.
 5. Influxdb installation : 
 ```sh
    $ curl -sL https://repos.influxdata.com/influxdb.key | sudo apt-key add -
    $ source /etc/lsb-release
    $ echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/influxdb.list
    $ sudo apt-get update && sudo apt-get install -y influxdb  
   ```
## Running  
### Terminus Deployment
  After doing SSH into the VM follow below procedure.
  1. Clone the Repository
  2. Fill the following AWS values in the ```Terminus_Interface/.env``` file
  
      ```   
         AWS_KEY=
         AWS_SECRET=
         AWS_DEFAULT_REGION=us-east-2
         AWS_KEY_PAIR_NAME=
         AWS_SUBNET_ID=example subnet-3eb5a357
         AWS_SECURITY_GROUP= example sg-07ba209aa8a4db1bd
         KUBE_CLUSTER_NAME=<name>.k8s.local
         KOPS_CONTAINER_NAME=<containerName>
         KOPS_S3_BUCKET_NAME=<name>-kops-state-store
      ```
  2. cd into ```Terminus_Interface/scripts``` directory
  3. Run the script using the commands
     ```
        chmod +x deploy_app.sh
        sudo sh deploy_app.sh
     ```
  4. Then use the web browser to visit 
        ```
            http://VM_IP:8082
        ```
### Data Dump Restore/Download

#### Mongodb dump restore
1. All the mongodb dump is present inside ```Terminus_Interface\mongo\mongodump directory.```
2. It can be restored to the started mongodb instance by running this command :

```sh
  $ cd  Terminus_Interface\mongo\mongodump
  $ mongorestore -h <VM_PUBLIC_IP> -u <username> -p <password> ./dump
  ``` 
#### Influxdb Data download
It needs to be only downloaded if you want to perform traning again or view the data other than MSC. 
All the collected monitoring data is stored in the <a href= "https://s3.console.aws.amazon.com/s3/buckets/terminusinfluxddatamain/?region=ap-southeast-2&tab=overview"> S3 bucket</a>.
To use those dataset for the initial setup, follow the below procedure.
 1. Make a directory to contain all the data  
 ```sh 
    $ sudo mkdir /data
 ```
 2. Now, gGet all the compressed influxdb files by running 
 ```sh 
    $ sudo s3cmd get --recursive s3://terminusinfluxdatamain/ /data
   ```
 3. Now we will have all the data inside the directory, we need to restore back to influxdb by running the following command:
 ```sh 
     $ sudo influxd restore -portable -host <VM_public_IP:8088> /data/sept/
   ```
  
### Initial Training and Setup ```[not required if exported mongo dump]```
 1. From web browser visit 
         ```
             http://VM_IP:8082
         ```
 2. Do all types of training for each services by selecting from the menu. It will take some time to complete the training.
 3. After the training Click on <b>Store All MSCs.</b> This will also take some time to complete.
 4. After that, click on <b>Store MSC Prediction Data</b>.
 5. Lastly, Do the KubeEventAnalysis for all the services.
 
### Testing and Visulaization
 
 #### Performance Modeling of Microservices
 ##### To Conduct a new Test
 1. Go to MSC->ConductTest->MSC Determine.
 2. Fill in the values of 
 ```
 1. Service Name: Choose from the selected services
 2. Replicas: For the Number of Replicas for that service
 3. Limits/Non-Limits: Whether to apply any limits on the service pod or not.
 4. Main-Instance Type: Select the Slave node Instance Type
 5. Limiting-Instance Type: The limits to be applied on the slave nodes. keep in mind that, this should always be less than Main-Instance type. 
 6. Master Node Instance Type: Type of Master node.
 7. Deployment Instance Type: Type of the instance for the deployment of agent and generating load.
 8. MAX RPS: Maximum RPS to which the testing needs to be done.
 9. Services: Number of services present in the file.
 10. NodeCount: Number of slave nodes. This will be dependent on the number of replicas required. 
 
 ```
 3. Click Start Test.
 4. Wait for the test to get completed. Data and logs in real time data can be seen by going to MSC->Ongoing/Completed Tests ->Select the Test.
    After that choose Grafana for viewing the Kubernetes and Load Generation data, kibana and logs for the Logs.   
 
 ##### Results
 1. Once everything is done go to Analyzed Results->MSC Analysis ALL -> Select the service -> Click Submit.
    You would see graphs for the experimental and predicted MSC.
 2. Other results can be seen like the comparison of limits and non-limits data from the menu.   
 
 #### Sandboxing
 
 ##### To record API responses
 1. To test Sandboxing approach, docker-compose file need to uploaded in ```Terminus_Interface\sandbox_microservices_service\data\docker_compose_files```
 2. After uploading, go to Sandboxing-Microservice->Conduct Test->API intercept
 3. Select the file name, enter the Main service name and api endpoint.
 4. Then click submit to record API responses.
 ##### Get Docker-compose.yml services network
 1. After the response interception is done, go to Sandboxing-Microservice->Conduct Test->API intercept->Get Yaml for service
 2. Select the file name and enter the main service name.
 3. Then click submit. Graph structure for all the services will be formed.
 ##### Get Docker-compose.yml for sandboxed services 
 1. After the response interception is done, go to Sandboxing-Microservice->Conduct Test->API intercept->Sandboxed Tree
 2. Enter the service name and choose between to get the compose file from the new version or the old one.
 3. Then click submit. Docker-compose.yml file will be provided.
 4. Afterwards, this file can be converted to Kubernetes yaml file and testing can be done on it for building performance model.
 
 Tree formed from the docker-compose.yml file of this application can be seen in below figure as an Example.
  
 <p align="center">
 <img width="550" height="271" src="https://github.com/ansjin/terminus/blob/master/Docs/app_tree.png"></img>
 </p>
 
## REST API Documentation
 
 It's automatically created for all the new API endpoints added. 
 The header needs to be added above every API. Here is an example for one: 
   ```
 // ConductTestToCalculatePodBootTime godoc
 // @Summary Collect the data for pod botting time
 // @Description Collect the data for pod botting time /conductTestToCalculatePodBootTime/1/1/1/t2.xlarge/40/t2.large/compute/primeapp/_api_prime
 // @Tags START_TEST
 // @Accept text/html
 // @Produce json
 // @Param numsvs query string true "number of services in microservice"
 // @Param appname query string true "name of application[primeapp, movieapp, webacapp]"
 // @Param apptype query string true "application type[compute, dbaccess, web]"
 // @Param replicas query string true "number of replicas of service "
 // @Param instancetype query string true " host instance type "
 // @Param limitinginstancetype query string true " limitinginstancetype instance type on main instance"
 // @Param apiendpoint query string true " api end point specify like _api_test for /api/test"
 // @Param maxRPS query string true " max RPS"
 // @Success 200 {string} string "ok"
 // @Failure 400 {string} string "ok"
 // @Failure 404 {string} string "ok"
 // @Failure 500 {string} string "ok"
 // @Router /conductTestToCalculatePodBootTime/{numsvs}/{nodecount}/{replicas}/{instancetype}/{mastertype}/{testvmtype}/{maxRPS}/{limitinginstancetype}/{apptype}/{appname}/{mainservicename}/{apiendpoint}/ [get]

   ```
  1. Add new headers or function.
  2. Download swag by using:
  ```sh
  $ go get -u github.com/swaggo/swag/cmd/swag
  ```
  
  3. Run `swag init` in the `Terminus_Interface\server\` folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).
  ```sh
  $ swag init
  ```  
  4. Copy `Terminus_Interface\server\docs\swagger\swagger.yaml` to `Terminus_Interface\server\assets\swagger.yaml`
  5. Once the website is running, you can access the Swager.yaml file 
  ```http://VM_IP:8082/assets/swagger.yaml```
  6. This yaml file then can be used on ```https://editor.swagger.io/``` to generate API documentation.
     
## Python Notebooks and other Useful Scripts
Some of the testing python notebooks and other useful scripts can be found in the importantScripts folder.
These notebooks are the test version for the linear regression with different parameters and data.  
 
## Help and Contibution

Please add issues if you have a question or found a problem. 

Pull requests are welcome too!

Contributions are most welcome. Please message me if you like the idea and want to contribute. 