#!/bin/bash

#shutdown RIC, CIB and Thumbor Instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/CIB/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'
cd ~/go/src/github.com/phzfi/RIC/thumbor



echo "Clearing thumbor temp files"
sudo rm -r /tmp/thumbor

echo "Clearing cache"
sudo -c 'sync && echo 3 >/proc/sys/vm/drop_caches'

#Starts thumbor in a new screen
sh thumbor_do_run.sh
