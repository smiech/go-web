#!/bin/bash
for i in "$@"
do case $i in
    -b=*|--bssid=*)
    BSSID="${i#*=}"
    shift # past argument=value
    ;;
    *)
          # unknown option
    ;;
esac
done

echo "BSSID monitored: ${BSSID}"
cd ./dumps
rm dump*
ps -aux | grep airodump-ng | grep -E "^root[ ]+[[:digit:]]+" | sed -r 's/root[ ]+([0-9]+)/\1/g' | cut -d " " -f 1 | xargs kill -9
ifconfig wlan0 down && iwconfig wlan0 mode monitor && ifconfig wlan0 up 
airodump-ng  wlan0 --essid-regex "${BSSID}" -c 36 --band abg -w dump --write-interval 15 -o csv