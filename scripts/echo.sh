#!/bin/bash
#for i in {1..100};do echo test;sleep 2; done
cd ./dumps
rm dump*
ifconfig wlan0 down && iwconfig wlan0 mode monitor && ifconfig wlan0 up 
airodump-ng  wlan0 --band abg -w dump --write-interval 15 -o csv