#!/bin/bash
cd ./dumps
rm dump*
airodump-ng  wlan0 --essid-regex 'FRITZ!Box 7' --band abg -w dump --write-interval 15 -o csv