package desktop

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Fred78290/kubernetes-desktop-autoscaler/constantes"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/context"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/uuid"
)

// VirtualMachine virtual machine wrapper
type VirtualMachine struct {
	Name string
	Uuid string
}

// GuestInfos the guest infos
// Must not start with `guestinfo.`
type GuestInfos map[string]string

func encodeMetadata(object interface{}) (string, error) {
	var result string
	out, err := json.Marshal(object)

	if err == nil {
		var stdout bytes.Buffer
		var zw = gzip.NewWriter(&stdout)

		zw.Name = "metadata"
		zw.ModTime = time.Now()

		if _, err = zw.Write(out); err == nil {
			if err = zw.Close(); err == nil {
				result = base64.StdEncoding.EncodeToString(stdout.Bytes())
			}
		}
	}

	return result, err
}

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

func (g GuestInfos) isEmpty() bool {
	return len(g) == 0
}

func (vm *VirtualMachine) findHardDrive(ctx *context.Context) (interface{}, error) {
	/// TODO

	return nil, errors.New("the given disk values match multiple disks")
}

func (vm *VirtualMachine) addOrExpandHardDrive(ctx *context.Context, diskSize int, expandHardDrive bool) error {
	var err error

	/// TODO

	return err
}

type VirtualMachineConfigSpec struct {
	NumCPUs      int32
	MemoryMB     int64
	InstanceUuid string
	ExtraConfig  map[string]string
}

func (vm *VirtualMachine) collectNetworkInfos(ctx *context.Context, network *Network, nodeIndex int) error {
	return nil
}

func (vm *VirtualMachine) addNetwork(ctx *context.Context, network *Network, nodeIndex int) error {
	return nil
}

// Configure set characteristic of VM a virtual machine
func (vm *VirtualMachine) Configure(ctx *context.Context, userName, authKey string, cloudInit interface{}, network *Network, annotation string, expandHardDrive bool, memory, cpus, disk, nodeIndex int) error {
	var err error

	vmConfigSpec := VirtualMachineConfigSpec{
		NumCPUs:      int32(cpus),
		MemoryMB:     int64(memory),
		InstanceUuid: string(uuid.NewUUID()),
	}

	if err = vm.addOrExpandHardDrive(ctx, disk, expandHardDrive); err != nil {

		err = fmt.Errorf(constantes.ErrUnableToAddHardDrive, vm.Name, err)

	} else if err = vm.addNetwork(ctx, network, nodeIndex); err != nil {

		err = fmt.Errorf(constantes.ErrUnableToAddNetworkCard, vm.Name, err)

	} else if vmConfigSpec.ExtraConfig, err = vm.cloudInit(ctx, vm.Name, userName, authKey, cloudInit, network, nodeIndex); err != nil {

		err = fmt.Errorf(constantes.ErrCloudInitFailCreation, vm.Name, err)

	} else if err = vm.Reconfigure(ctx, vmConfigSpec); err != nil {

		err = fmt.Errorf(constantes.ErrUnableToReconfigureVM, vm.Name, err)

	}

	return err
}

func (vm *VirtualMachine) Reconfigure(ctx *context.Context, vmConfigSpec VirtualMachineConfigSpec) error {
	return nil
}

// IsToolsRunning returns true if VMware Tools is currently running in the guest OS, and false otherwise.
func (vm VirtualMachine) IsToolsRunning(ctx *context.Context) (bool, error) {
	/// TODO

	return false, nil
}

// IsSimulatorRunning returns true if VMware Tools is currently running in the guest OS, and false otherwise.
func (vm VirtualMachine) IsSimulatorRunning(ctx *context.Context) bool {
	return false
}

func (vm *VirtualMachine) waitForToolsRunning(ctx *context.Context) (bool, error) {
	var running bool

	/// TODO

	return running, nil
}

func (vm *VirtualMachine) ListAddresses(ctx *context.Context) ([]NetworkInterface, error) {
	addresses := make([]NetworkInterface, 0)

	/// TODO

	return addresses, nil
}

// WaitForToolsRunning wait vmware tool starts
func (vm *VirtualMachine) WaitForToolsRunning(ctx *context.Context) (bool, error) {
	var err error
	var running bool

	/// TODO

	return running, err
}

func (vm *VirtualMachine) waitForIP(ctx *context.Context) (string, error) {
	/// TODO

	return "", fmt.Errorf("VMWare tools is not running on the VM:%s, unable to retrieve IP", vm.Name)
}

// WaitForIP wait ip
func (vm *VirtualMachine) WaitForIP(ctx *context.Context) (string, error) {
	var err error
	var ip string

	/// TODO

	return ip, err
}

// PowerOn power on a virtual machine
func (vm *VirtualMachine) PowerOn(ctx *context.Context) error {
	var err error

	/// TODO

	return err
}

// PowerOff power off a virtual machine
func (vm *VirtualMachine) PowerOff(ctx *context.Context) error {
	var err error

	/// TODO

	return err
}

// ShutdownGuest power off a virtual machine
func (vm *VirtualMachine) ShutdownGuest(ctx *context.Context) error {
	var err error

	/// TODO

	return err
}

// Delete delete the virtual machine
func (vm *VirtualMachine) Delete(ctx *context.Context) error {
	var err error

	/// TODO

	return err
}

// Status refresh status virtual machine
func (vm *VirtualMachine) Status(ctx *context.Context) (*Status, error) {
	var err error
	var status *Status = &Status{}

	/// TODO

	return status, err
}

// SetGuestInfo change guest ingos
func (vm *VirtualMachine) SetGuestInfo(ctx *context.Context, guestInfos *GuestInfos) error {
	var err error

	/// TODO

	return err
}

func (vm *VirtualMachine) cloudInit(ctx *context.Context, hostName string, userName, authKey string, cloudInit interface{}, network *Network, nodeIndex int) (GuestInfos, error) {
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
