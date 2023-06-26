#!/bin/bash
#
# Common setup for azure iot edge binaries

set -euxo pipefail
sudo swapoff -a

sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg \
    python3-pip vim make \
    --no-install-recommends
sudo pip3 install pip --upgrade

# Install golang
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz -O go1.20.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
rm -rf go1.20.5.linux-amd64.tar.gz
sudo echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
sudo echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile

# Install nerdctl and containerd
curl -fsSL https://github.com/containerd/nerdctl/releases/download/v1.4.0/nerdctl-full-1.4.0-linux-amd64.tar.gz \
    -o nerdctl.tar.gz
mkdir nerdctl
tar -xvf nerdctl.tar.gz -C ./nerdctl
sudo mkdir -p /opt/cni/bin
sudo cp ./nerdctl/libexec/cni/* /opt/cni/bin
sudo cp ./nerdctl/bin/* /usr/local/bin
sudo cp ./nerdctl/lib/systemd/system/* /lib/systemd/system
sudo systemctl enable containerd.service
sudo systemctl start containerd.service
sudo systemctl enable buildkit.service
sudo systemctl start buildkit.service

sudo rm -rf nerdctl*

# Create nerdctld group
sudo addgroup docker
sudo usermod -aG docker vagrant
sudo newgrp docker

# Grant all users in the group "nerdctl" to access without typing sudo
sudo mkdir -p /etc/systemd/system/nerdctl.socket.d
sudo cp /vagrant/systemd/10-group.conf /etc/systemd/system/nerdctl.socket.d/10-group.conf

# # Copy nerdctl daemon socket and service
# sudo cp /vagrant/systemd/nerdctl.service /usr/lib/systemd/system/nerdctl.service
# sudo cp /vagrant/systemd/nerdctl.socket /usr/lib/systemd/system/nerdctl.socket
# sudo systemctl enable nerdctl.socket
# sudo systemctl enable nerdctl.service

# Copy to another folder