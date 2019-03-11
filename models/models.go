package models

type ExecuteData struct {
	CommandId string
	Data      string
	Time      int64
}

type NetworkInfo struct {
	Clients      []NetworkClient
	AccessPoints []AccessPoint
}

type NetworkClient struct {
	APMac       string
	Mac         string
	ProbedSSIDs []string
}
type AccessPoint struct {
	Mac     string
	Channel int
	Name    string
}
