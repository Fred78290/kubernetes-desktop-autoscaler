package desktop

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Fred78290/kubernetes-desktop-autoscaler/api"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/constantes"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/uuid"
)

// VirtualMachine virtual machine wrapper
type VirtualMachine struct {
	Name   string
	Uuid   string
	Vmx    string
	Vcpus  int32
	Memory int64
}

// GuestInfos the guest infos
// Must not start with `guestinfo.`
type GuestInfos map[string]string

func encodeCloudInit(name string, object interface{}) (string, error) {
	var result string
	var out bytes.Buffer
	var err error

	fmt.Fprintln(&out, "#cloud-init")

	wr := yaml.NewEncoder(&out)
	err = wr.Encode(object)

	wr.Close()

	if err == nil {
		var stdout bytes.Buffer
		var zw = gzip.NewWriter(&stdout)

		zw.Name = name
		zw.ModTime = time.Now()

		if _, err = zw.Write(out.Bytes()); err == nil {
			if err = zw.Close(); err == nil {
				result = base64.StdEncoding.EncodeToString(stdout.Bytes())
			}
		}
	}

	return result, err
}

func encodeObject(name string, object interface{}) (string, error) {
	var result string
	out, err := yaml.Marshal(object)

	if err == nil {
		var stdout bytes.Buffer
		var zw = gzip.NewWriter(&stdout)

		zw.Name = name
		zw.ModTime = time.Now()

		if _, err = zw.Write(out); err == nil {
			if err = zw.Close(); err == nil {
				result = base64.StdEncoding.EncodeToString(stdout.Bytes())
			}
		}
	}

	return result, err
}

func buildVendorData(userName, authKey string) interface{} {
	tz, _ := time.Now().Zone()

	return map[string]interface{}{
		"package_update":  true,
		"package_upgrade": true,
		"timezone":        tz,
		"users": []string{
			"default",
		},
		"ssh_authorized_keys": []string{
			authKey,
		},
		"system_info": map[string]interface{}{
			"default_user": map[string]string{
				"name": userName,
			},
		},
	}
}

func buildCloudInit(hostName string, userName, authKey string, cloudInit interface{}, network *Network, nodeIndex int) (GuestInfos, error) {
	var metadata, userdata, vendordata string
	var err error
	var guestInfos GuestInfos
	var fqdn string

	if len(network.Domain) > 0 {
		fqdn = fmt.Sprintf("%s.%s", hostName, network.Domain)
	}

	netconfig := &NetworkConfig{
		InstanceID:    string(uuid.NewUUID()),
		LocalHostname: hostName,
		Hostname:      fqdn,
	}

	if network != nil && len(network.Interfaces) > 0 {
		netconfig.Network = network.GetCloudInitNetwork(nodeIndex)
	}

	if metadata, err = encodeObject("metadata", netconfig); err != nil {
		err = fmt.Errorf(constantes.ErrUnableToEncodeGuestInfo, "metadata", err)
	}

	if err == nil {
		if cloudInit != nil {
			if userdata, err = encodeCloudInit("userdata", cloudInit); err != nil {
				return nil, fmt.Errorf(constantes.ErrUnableToEncodeGuestInfo, "userdata", err)
			}
		} else if userdata, err = encodeCloudInit("userdata", map[string]string{}); err != nil {
			return nil, fmt.Errorf(constantes.ErrUnableToEncodeGuestInfo, "userdata", err)
		}

		if len(userName) > 0 && len(authKey) > 0 {
			if vendordata, err = encodeCloudInit("vendordata", buildVendorData(userName, authKey)); err != nil {
				return nil, fmt.Errorf(constantes.ErrUnableToEncodeGuestInfo, "vendordata", err)
			}
		} else if vendordata, err = encodeCloudInit("vendordata", map[string]string{}); err != nil {
			return nil, fmt.Errorf(constantes.ErrUnableToEncodeGuestInfo, "vendordata", err)
		}

		guestInfos = GuestInfos{
			"metadata":            metadata,
			"metadata.encoding":   "gzip+base64",
			"userdata":            userdata,
			"userdata.encoding":   "gzip+base64",
			"vendordata":          vendordata,
			"vendordata.encoding": "gzip+base64",
		}
	}

	return guestInfos, nil
}

func buildNetworkInterface(interfaces []*NetworkInterface, nodeIndex int) []*api.NetworkInterface {
	result := make([]*api.NetworkInterface, 0, len(interfaces))

	for _, inf := range interfaces {
		result = append(result, &api.NetworkInterface{
			Macaddress: inf.GetMacAddress(nodeIndex),
			Vnet:       inf.VNet,
			Type:       inf.ConnectionType,
			Device:     inf.VirtualDev,
		})
	}
	return result
}
