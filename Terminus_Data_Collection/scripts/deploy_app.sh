#!/bin/sh
#sudo apt-get update
#iptables -A INPUT -p tcp -m tcp --dport 3000 -m state --state NEW -j ACCEPT
#iptables -A INPUT -p tcp -m tcp --dport 3001 -m state --state NEW -j ACCEPT
#iptables -A INPUT -p tcp -m tcp --dport 3306 -m state --state NEW -j ACCEPT
#iptables-restore < /etc/iptables/rules.v4
myip="$(dig +short myip.opendns.com @resolver1.opendns.com)"
echo "My Public IP address: ${myip}"
sed -i "s/SERVER_PUBLIC_IP_ADDRESS/$myip/g" ../data/etc/heapster.yaml
apt-get remove -y docker docker-engine docker.io
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce
service docker start
curl -L https://github.com/docker/compose/releases/download/1.13.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
cd ..
sudo docker-compose up --build &
read -p "Press enter to Exit"
