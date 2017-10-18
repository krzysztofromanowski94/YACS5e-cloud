#!/usr/bin/python
import time
import os

print "\n#########################################################\n"
print time.strftime("%Y-%m-%d %H:%M")

# shut down docker-compose backend
os.system("su -c 'docker-compose -f /home/backend/.backend/docker-compose.yml down' -s /bin/bash backend")

# create backup of backend services
tar_name="aws-backup-backend." + time.strftime("%Y-%m-%d-%H-%M") + ".tar.gz"
os.system("tar -zcvf /tmp/aws-backup-backend.tar.gz -C /home/backend/ .backend")
os.system("mv /tmp/aws-backup-backend.tar.gz /tmp/" + tar_name)

# turn on backend services
os.system("su -c 'docker-compose -f /home/backend/.backend/docker-compose.yml up -d' -s /bin/bash backend")

# get backups from server
backend_backup_list = os.popen("gdrive list | grep aws-backup-backend | awk '{print $2 \" \" $1}' | sort").read().split("\n")

# remove last empty element
backend_backup_list.pop()

if len(backend_backup_list) >= 5:
  os.system("gdrive delete " + backend_backup_list[0].split(" ")[1])

os.system("gdrive upload /tmp/" + tar_name)

# remove locally created backup
os.system("rm -rf /tmp/" + tar_name)



##########################################################################

# shut down docker-compose services
os.system("su -c 'docker-compose -f /home/services/.services/docker-compose.yml down' -s /bin/bash services")

# create backup of services
tar_name="aws-backup-services." + time.strftime("%Y-%m-%d-%H-%M") + ".tar.gz"
os.system("tar -zcvf /tmp/aws-backup-services.tar.gz -C /home/services/ .services")
os.system("mv /tmp/aws-backup-services.tar.gz /tmp/" + tar_name)

# turn on services
os.system("su -c 'docker-compose -f /home/services/.services/docker-compose.yml up -d' -s /bin/bash services")

# get backups from server
services_backup_list = os.popen("gdrive list | grep aws-backup-services | awk '{print $2 \" \" $1}' | sort").read().split("\n")

# remove last empty element
services_backup_list.pop()

if len(services_backup_list) >= 5:
  os.system("gdrive delete " + services_backup_list[0].split(" ")[1])

os.system("gdrive upload /tmp/" + tar_name)

# remove locally created backup
os.system("rm -rf /tmp/" + tar_name)


#########################################################################

# create backup of .admin dir
tar_name="aws-backup-admin." + time.strftime("%Y-%m-%d-%H-%M") + ".tar.gz"
os.system("tar -zcvf /tmp/aws-backup-admin.tar.gz -C /home/ec2-user/ .admin")
os.system("mv /tmp/aws-backup-admin.tar.gz /tmp/" + tar_name)

# get backups from server
admin_backup_list = os.popen("gdrive list | grep aws-backup-admin | awk '{print $2 \" \" $1}' | sort").read().split("\n")

# remove last empty element
admin_backup_list.pop()

if len(admin_backup_list) >= 5:
  os.system("gdrive delete " + admin_backup_list[0].split(" ")[1])

os.system("gdrive upload /tmp/" + tar_name)

# remove locally created backup
os.system("rm -rf /tmp/" + tar_name)

