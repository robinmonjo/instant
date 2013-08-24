#!/bin/bash

sudo apt-get update
sudo apt-get install -y lxc golang

echo "GOPATH=/go" >> ~/.bashrc
source ~/.bashrc