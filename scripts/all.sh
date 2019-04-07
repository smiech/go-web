#!/bin/bash
#for i in {1..100};do echo test;sleep 2; done
cd ./dumps
rm dump*
ps -aux | grep airodump-ng | grep -E "^root[ ]+[[:digit:]]+" | sed -r 's/root[ ]+([0-9]+)/\1/g' | cut -d " " -f 1 | xargs kill -9
ifconfig wlan0 down && iwconfig wlan0 mode monitor && ifconfig wlan0 up 
airodump-ng  wlan0 --band abg -w dump --write-interval 15 -o csv
