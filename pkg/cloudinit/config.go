package cloudinit

import (
	"github.com/Fred78290/kubernetes-desktop-autoscaler/desktop"
)

// BuildCloudInit build map for guestinfo
func BuildCloudInit(hostName string, userName, authKey string, cloudInit interface{}, network *desktop.Network, nodeIndex int, allowUpgrade bool) (desktop.GuestInfos, error) {
	return desktop.BuildCloudInit(hostName, userName, authKey, cloudInit, network, nodeIndex, allowUpgrade)
}
