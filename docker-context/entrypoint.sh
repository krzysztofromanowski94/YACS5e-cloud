#!/bin/sh

ansible-playbook /root/startup.yml

cd /go/src/server
go install
cd ~/
server
