#!/bin/sh

ansible-playbook /root/startup.yml -vv
exec server
