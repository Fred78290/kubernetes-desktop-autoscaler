package cloudinit

import (
	"github.com/Fred78290/kubernetes-desktop-autoscaler/desktop"
)

// BuildCloudInit build map for guestinfo
func BuildCloudInit(hostName, userName, authKey, tz string, cloudInit interface{}, network *desktop.Network, nodeIndex int, allowUpgrade bool) (desktop.GuestInfos, error) {
	return desktop.BuildCloudInit(hostName, userName, authKey, tz, cloudInit, network, nodeIndex, allowUpgrade)
}
