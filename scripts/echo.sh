#!/bin/bash
for i in {1..100};do echo test;sleep 2; done
#cd ./dumps
#rm dump*
#airodump-ng  wlan0 --essid-regex 'FRITZ!Box 7' --band abg -w dump --write-interval 15 -o csv