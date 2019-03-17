package fileParser

import (
	"log"
	"strconv"
	"strings"

	models "github.com/smiech/go-web/models"
)

const _APHeader = "BSSID, First time seen, Last time seen, channel, Speed, Privacy, Cipher, Authentication, Power, # beacons, # IV, LAN IP, ID-length, ESSID, Key"
const _ClientHeader = "Station MAC, First time seen, Last time seen, Power, # packets, BSSID, Probed ESSIDs"

func Parse(dumpContent string) (models.NetworkInfo, error) {
	const clientInfoColumnCount = 6
	temp := strings.Split(dumpContent, "\n")
	isAPSection := false
	isClientSection := false
	networkInfo := models.NetworkInfo{}
	for _, element := range temp {
		elementTrimmed := strings.TrimSpace(element)
		if elementTrimmed == "" {
			continue
		}
		if elementTrimmed == _APHeader {
			isAPSection = true
			log.Println("Parsing AP section")
			continue
		} else if elementTrimmed == _ClientHeader {
			isClientSection = true
			isAPSection = false
			log.Println("Parsing Client section")
			continue
		}

		if isAPSection {
			var newAP = models.AccessPoint{}
			apColumns := strings.Split(elementTrimmed, ",")
			if len(apColumns) < 14 {
				log.Println("Couldn't parse the AP line: ", apColumns)
				continue
			}
			newAP.Mac = strings.TrimSpace(apColumns[0])
			channel, _ := strconv.ParseInt(strings.TrimSpace(apColumns[3]), 0, 32)
			newAP.Channel = int(channel)
			newAP.Name = strings.TrimSpace(apColumns[13])
			//log.Println(newAP)
			networkInfo.AccessPoints = append(networkInfo.AccessPoints, newAP)
		}
		if isClientSection {
			var newClient = models.NetworkClient{}
			clientColumn := strings.Split(elementTrimmed, ",")
			if len(clientColumn) < clientInfoColumnCount {
				log.Println("Couldn't parse the client line: ", clientColumn)
				continue
			}
			newClient.APMac = strings.TrimSpace(clientColumn[5])
			newClient.Mac = strings.TrimSpace(clientColumn[0])
			if len(clientColumn) > clientInfoColumnCount {
				var probedSsidsCsv = strings.Join(clientColumn[6:], ",")
				var probedSsidsSplitted = strings.Split(probedSsidsCsv, ",")
				if !(len(probedSsidsSplitted) == 1 && strings.TrimSpace(probedSsidsSplitted[0]) == "") {
					newClient.ProbedSSIDs = strings.Split(probedSsidsCsv, ",")
				}
			}
			//log.Println(newClient)
			networkInfo.Clients = append(networkInfo.Clients, newClient)
		}

		//log.Println("CLIENT: ", elementTrimmed)
	}
	//log.Printf("Element: %v at index %v", element, index)
	// index is the index where we are
	// element is the element from someSlice for where we are

	return networkInfo, nil
}
