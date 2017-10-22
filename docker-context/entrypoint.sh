#!/bin/sh

ansible-playbook /root/startup.yml

cd /go/src/server
go get ./...
go install
cd ~/
server
