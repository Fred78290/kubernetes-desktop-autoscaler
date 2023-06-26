package server

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Fred78290/kubernetes-desktop-autoscaler/constantes"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/context"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/desktop"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/types"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/utils"
	glog "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	uid "k8s.io/apimachinery/pkg/types"
)

// AutoScalerServerNodeState VM state
type AutoScalerServerNodeState int32

// AutoScalerServerNodeType node class (external, autoscaled, managed)
type AutoScalerServerNodeType int32

// autoScalerServerNodeStateString strings
var autoScalerServerNodeStateString = []string{
	"AutoScalerServerNodeStateNotCreated",
	"AutoScalerServerNodeStateCreating",
	"AutoScalerServerNodeStateRunning",
	"AutoScalerServerNodeStateStopped",
	"AutoScalerServerNodeStateDeleted",
	"AutoScalerServerNodeStateUndefined",
}

const (
	// AutoScalerServerNodeStateNotCreated not created state
	AutoScalerServerNodeStateNotCreated = iota

	// AutoScalerServerNodeStateCreating running state
	AutoScalerServerNodeStateCreating

	// AutoScalerServerNodeStateRunning running state
	AutoScalerServerNodeStateRunning

	// AutoScalerServerNodeStateStopped stopped state
	AutoScalerServerNodeStateStopped

	// AutoScalerServerNodeStateDeleted deleted state
	AutoScalerServerNodeStateDeleted

	// AutoScalerServerNodeStateUndefined undefined state
	AutoScalerServerNodeStateUndefined
)

const (
	// AutoScalerServerNodeExternal is a node create out of autoscaler
	AutoScalerServerNodeExternal = iota
	// AutoScalerServerNodeAutoscaled is a node create by autoscaler
	AutoScalerServerNodeAutoscaled
	// AutoScalerServerNodeManaged is a node managed by controller
	AutoScalerServerNodeManaged
)

// AutoScalerServerNode Describe a AutoScaler VM
type AutoScalerServerNode struct {
	NodeGroupID      string                    `json:"group"`
	NodeName         string                    `json:"name"`
	NodeIndex        int                       `json:"index"`
	VMUUID           string                    `json:"vm-uuid"`
	CRDUID           uid.UID                   `json:"crd-uid"`
	Memory           int                       `json:"memory"`
	CPU              int                       `json:"cpu"`
	Disk             int                       `json:"disk"`
	IPAddress        string                    `json:"address"`
	State            AutoScalerServerNodeState `json:"state"`
	NodeType         AutoScalerServerNodeType  `json:"type"`
	ControlPlaneNode bool                      `json:"control-plane,omitempty"`
	AllowDeployment  bool                      `json:"allow-deployment,omitempty"`
	ExtraLabels      types.KubernetesLabel     `json:"labels,omitempty"`
	ExtraAnnotations types.KubernetesLabel     `json:"annotations,omitempty"`
	Configuration    *desktop.Configuration    `json:"vmware"`
	serverConfig     *types.AutoScalerServerConfig
}

func (s AutoScalerServerNodeState) String() string {
	return autoScalerServerNodeStateString[s]
}

func (vm *AutoScalerServerNode) waitReady(c types.ClientGenerator) error {
	glog.Debugf("AutoScalerNode::waitReady, node:%s", vm.NodeName)

	return c.WaitNodeToBeReady(vm.NodeName)
}

func (vm *AutoScalerServerNode) recopyEtcdSslFilesIfNeeded() error {
	var err error

	if vm.ControlPlaneNode || *vm.serverConfig.UseExternalEtdc {
		glog.Infof("Recopy etcd certs for node:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

		if err = utils.Scp(vm.serverConfig.SSH, vm.IPAddress, vm.serverConfig.ExtSourceEtcdSslDir, "."); err != nil {
			glog.Errorf("scp failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("mkdir -p %s", vm.serverConfig.ExtDestinationEtcdSslDir)); err != nil {
			glog.Errorf("mkdir failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("cp -r %s/* %s", filepath.Base(vm.serverConfig.ExtSourceEtcdSslDir), vm.serverConfig.ExtDestinationEtcdSslDir)); err != nil {
			glog.Errorf("mv failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("chown -R root:root %s", vm.serverConfig.ExtDestinationEtcdSslDir)); err != nil {
			glog.Errorf("chown failed: %v", err)
		}
	}

	return err
}

func (vm *AutoScalerServerNode) recopyKubernetesPKIIfNeeded() error {
	var err error

	if vm.ControlPlaneNode {
		glog.Infof("Recopy pki kubernetes certs for node:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

		if err = utils.Scp(vm.serverConfig.SSH, vm.IPAddress, vm.serverConfig.KubernetesPKISourceDir, "."); err != nil {
			glog.Errorf("scp failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("mkdir -p %s", vm.serverConfig.KubernetesPKIDestDir)); err != nil {
			glog.Errorf("mkdir failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("cp -r %s/* %s", filepath.Base(vm.serverConfig.KubernetesPKISourceDir), vm.serverConfig.KubernetesPKIDestDir)); err != nil {
			glog.Errorf("mv failed: %v", err)
		} else if _, err = utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, fmt.Sprintf("chown -R root:root %s", vm.serverConfig.KubernetesPKIDestDir)); err != nil {
			glog.Errorf("chown failed: %v", err)
		}
	}

	return err
}

func (vm *AutoScalerServerNode) waitForSshReady() error {
	return context.PollImmediate(time.Second, time.Duration(vm.serverConfig.SSH.WaitSshReadyInSeconds)*time.Second, func() (done bool, err error) {
		if _, err := utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, "ls"); err != nil {
			if strings.HasSuffix(err.Error(), "connection refused") {
				glog.Warnf("Wait ssh ready for node: %s, address: %s.", vm.NodeName, vm.IPAddress)
				return false, nil
			}

			return false, fmt.Errorf("unable to ssh: %s, reason:%v", vm.NodeName, err)
		}

		return true, nil
	})
}

func (vm *AutoScalerServerNode) kubeAdmJoin(c types.ClientGenerator) error {
	kubeAdm := vm.serverConfig.KubeAdm

	args := []string{
		"kubeadm",
		"join",
		kubeAdm.Address,
		"--token",
		kubeAdm.Token,
		"--discovery-token-ca-cert-hash",
		kubeAdm.CACert,
		"--apiserver-advertise-address",
		vm.IPAddress,
	}

	if vm.ControlPlaneNode {
		args = append(args, "--control-plane")
	}

	// Append extras arguments
	if len(kubeAdm.ExtraArguments) > 0 {
		args = append(args, kubeAdm.ExtraArguments...)
	}

	command := strings.Join(args, " ")

	glog.Infof("Join cluster for node:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

	if out, err := utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, command); err != nil {
		return fmt.Errorf("unable to execute command: %s, output: %s, reason:%v", command, out, err)
	}

	// To be sure, with kubeadm 1.26.1, the kubelet is not correctly restarted
	time.Sleep(5 * time.Second)

	return context.PollImmediate(5*time.Second, time.Duration(vm.serverConfig.SSH.WaitSshReadyInSeconds)*time.Second, func() (done bool, err error) {
		if node, err := c.GetNode(vm.NodeName); err == nil && node != nil {
			return true, nil
		}

		glog.Infof("Restart kubelet for node:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

		if out, err := utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, "systemctl restart kubelet"); err != nil {
			return false, fmt.Errorf("unable to restart kubelet, output: %s, reason:%v", out, err)
		}

		return false, nil
	})
}

func (vm *AutoScalerServerNode) k3sAgentJoin(c types.ClientGenerator) error {
	kubeAdm := vm.serverConfig.KubeAdm
	k3s := vm.serverConfig.K3S
	args := []string{
		fmt.Sprintf("echo K3S_ARGS='--kubelet-arg=provider-id=%s --node-name=%s --server=https://%s --token=%s' > /etc/systemd/system/k3s.service.env", vm.generateProviderID(), vm.NodeName, kubeAdm.Address, kubeAdm.Token),
	}

	if vm.ControlPlaneNode {
		args = append(args, "echo 'K3S_MODE=server' > /etc/default/k3s", "echo K3S_DISABLE_ARGS='--disable=servicelb --disable=traefik --disable=metrics-server' > /etc/systemd/system/k3s.disabled.env")

		if vm.serverConfig.UseExternalEtdc != nil && *vm.serverConfig.UseExternalEtdc {
			args = append(args, fmt.Sprintf("echo K3S_SERVER_ARGS='--datastore-endpoint=%s --datastore-cafile=%s/ca.pem --datastore-certfile=%s/etcd.pem --datastore-keyfile=%s/etcd-key.pem' > /etc/systemd/system/k3s.server.env", k3s.DatastoreEndpoint, vm.serverConfig.ExtDestinationEtcdSslDir, vm.serverConfig.ExtDestinationEtcdSslDir, vm.serverConfig.ExtDestinationEtcdSslDir))
		}
	}

	// Append extras arguments
	if len(k3s.ExtraCommands) > 0 {
		args = append(args, k3s.ExtraCommands...)
	}

	args = append(args, "systemctl enable k3s.service", "systemctl start k3s.service")

	glog.Infof("Join cluster for node:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

	command := fmt.Sprintf("sh -c \"%s\"", strings.Join(args, " && "))
	if out, err := utils.Sudo(vm.serverConfig.SSH, vm.IPAddress, vm.Configuration.Timeout, command); err != nil {
		return fmt.Errorf("unable to execute command: %s, output: %s, reason:%v", command, out, err)
	}

	return context.PollImmediate(5*time.Second, time.Duration(vm.serverConfig.SSH.WaitSshReadyInSeconds)*time.Second, func() (done bool, err error) {
		if node, err := c.GetNode(vm.NodeName); err == nil && node != nil {
			return true, nil
		}

		return false, nil
	})
}

func (vm *AutoScalerServerNode) joinCluster(c types.ClientGenerator) error {
	if vm.serverConfig.UseK3S != nil && *vm.serverConfig.UseK3S {
		return vm.k3sAgentJoin(c)
	} else {
		return vm.kubeAdmJoin(c)
	}
}

func (vm *AutoScalerServerNode) setNodeLabels(c types.ClientGenerator, nodeLabels, systemLabels types.KubernetesLabel) error {
	labels := utils.MergeKubernetesLabel(nodeLabels, systemLabels, vm.ExtraLabels)

	if err := c.LabelNode(vm.NodeName, labels); err != nil {
		return fmt.Errorf(constantes.ErrLabelNodeReturnError, vm.NodeName, err)
	}

	annotations := types.KubernetesLabel{
		constantes.AnnotationNodeGroupName:        vm.NodeGroupID,
		constantes.AnnotationInstanceID:           vm.VMUUID,
		constantes.AnnotationScaleDownDisabled:    strconv.FormatBool(vm.NodeType != AutoScalerServerNodeAutoscaled),
		constantes.AnnotationNodeAutoProvisionned: strconv.FormatBool(vm.NodeType == AutoScalerServerNodeAutoscaled),
		constantes.AnnotationNodeManaged:          strconv.FormatBool(vm.NodeType == AutoScalerServerNodeManaged),
		constantes.AnnotationNodeIndex:            strconv.Itoa(vm.NodeIndex),
	}

	annotations = utils.MergeKubernetesLabel(annotations, vm.ExtraAnnotations)

	if err := c.AnnoteNode(vm.NodeName, annotations); err != nil {
		return fmt.Errorf(constantes.ErrAnnoteNodeReturnError, vm.NodeName, err)
	}

	if vm.ControlPlaneNode && vm.AllowDeployment {
		if err := c.TaintNode(vm.NodeName,
			apiv1.Taint{
				Key:    constantes.NodeLabelControlPlaneRole,
				Effect: apiv1.TaintEffectNoSchedule,
				TimeAdded: &metav1.Time{
					Time: time.Now(),
				},
			},
			apiv1.Taint{
				Key:    constantes.NodeLabelMasterRole,
				Effect: apiv1.TaintEffectNoSchedule,
				TimeAdded: &metav1.Time{
					Time: time.Now(),
				},
			}); err != nil {
			return fmt.Errorf(constantes.ErrTaintNodeReturnError, vm.NodeName, err)
		}
	}

	return nil
}

func (vm *AutoScalerServerNode) launchVM(c types.ClientGenerator, nodeLabels, systemLabels types.KubernetesLabel) error {
	glog.Debugf("AutoScalerNode::launchVM, node:%s", vm.NodeName)

	var err error
	var status AutoScalerServerNodeState

	config := vm.Configuration
	network := config.Network
	userInfo := vm.serverConfig.SSH

	glog.Infof("Launch VM:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)

	if vm.State != AutoScalerServerNodeStateNotCreated {
		return fmt.Errorf(constantes.ErrVMAlreadyCreated, vm.NodeName)
	}

	if config.Exists(vm.NodeName) {
		glog.Warnf(constantes.ErrVMAlreadyExists, vm.NodeName)
		return fmt.Errorf(constantes.ErrVMAlreadyExists, vm.NodeName)
	}

	vm.State = AutoScalerServerNodeStateCreating

	if vm.NodeType != AutoScalerServerNodeAutoscaled && vm.NodeType != AutoScalerServerNodeManaged {

		err = fmt.Errorf(constantes.ErrVMNotProvisionnedByMe, vm.NodeName)

	} else if vm.VMUUID, err = config.Create(vm.NodeName, userInfo.GetUserName(), userInfo.GetAuthKeys(), config.TimeZone, vm.serverConfig.CloudInit, network, true, vm.Memory, vm.CPU, vm.Disk, vm.NodeIndex, config.AllowUpgrade); err != nil {

		err = fmt.Errorf(constantes.ErrUnableToLaunchVM, vm.NodeName, err)

	} else if err = config.PowerOn(vm.VMUUID); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if err = config.WaitForPowerState(vm.VMUUID, true); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if err = config.SetAutoStart(vm.VMUUID, true); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if _, err = config.WaitForIP(vm.VMUUID, config.Timeout); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if status, err = vm.statusVM(); err != nil {

		err = fmt.Errorf(constantes.ErrGetVMInfoFailed, vm.NodeName, err)

	} else if status != AutoScalerServerNodeStateRunning {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if err = vm.waitForSshReady(); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if err = vm.recopyKubernetesPKIIfNeeded(); err != nil {

		err = fmt.Errorf(constantes.ErrRecopyKubernetesPKIFailed, vm.NodeName, err)

	} else if err = vm.recopyEtcdSslFilesIfNeeded(); err != nil {

		err = fmt.Errorf(constantes.ErrUpdateEtcdSslFailed, vm.NodeName, err)

	} else if err = vm.joinCluster(c); err != nil {

		err = fmt.Errorf(constantes.ErrKubeAdmJoinFailed, vm.NodeName, err)

	} else if err = vm.setProviderID(c); err != nil {

		err = fmt.Errorf(constantes.ErrProviderIDNotConfigured, vm.NodeName, err)

	} else if err = vm.waitReady(c); err != nil {

		err = fmt.Errorf(constantes.ErrNodeIsNotReady, vm.NodeName)

	} else {
		err = vm.setNodeLabels(c, nodeLabels, systemLabels)
	}

	if err == nil {
		glog.Infof("Launched VM:%s for nodegroup: %s", vm.NodeName, vm.NodeGroupID)
	} else {
		glog.Errorf("Unable to launch VM:%s for nodegroup: %s. Reason: %v", vm.NodeName, vm.NodeGroupID, err.Error())
	}

	return err
}

func (vm *AutoScalerServerNode) startVM(c types.ClientGenerator) error {
	glog.Debugf("AutoScalerNode::startVM, node:%s", vm.NodeName)

	var err error
	var state AutoScalerServerNodeState

	glog.Infof("Start VM:%s", vm.NodeName)

	config := vm.Configuration

	if vm.NodeType != AutoScalerServerNodeAutoscaled && vm.NodeType != AutoScalerServerNodeManaged {

		err = fmt.Errorf(constantes.ErrVMNotProvisionnedByMe, vm.NodeName)

	} else if state, err = vm.statusVM(); err != nil {

		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

	} else if state == AutoScalerServerNodeStateStopped {

		if err = config.PowerOn(vm.VMUUID); err != nil {

			err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

		} else if err = config.WaitForPowerState(vm.VMUUID, true); err != nil {

			err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

		} else if _, err = config.WaitForIP(vm.VMUUID, config.Timeout); err != nil {

			err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

		} else if state, err = vm.statusVM(); err != nil {

			err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

		} else if state != AutoScalerServerNodeStateRunning {

			err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, err)

		} else {
			if err = c.UncordonNode(vm.NodeName); err != nil {
				glog.Errorf(constantes.ErrUncordonNodeReturnError, vm.NodeName, err)

				err = nil
			}

			vm.State = AutoScalerServerNodeStateRunning
		}
	} else if state != AutoScalerServerNodeStateRunning {
		err = fmt.Errorf(constantes.ErrStartVMFailed, vm.NodeName, fmt.Sprintf("Unexpected state: %d", state))
	}

	if err == nil {
		glog.Infof("Started VM:%s", vm.NodeName)
	} else {
		glog.Errorf("Unable to start VM:%s. Reason: %v", vm.NodeName, err)
	}

	return err
}

func (vm *AutoScalerServerNode) stopVM(c types.ClientGenerator) error {
	glog.Debugf("AutoScalerNode::stopVM, node:%s", vm.NodeName)

	var err error
	var state AutoScalerServerNodeState

	glog.Infof("Stop VM:%s", vm.NodeName)

	config := vm.Configuration

	if vm.NodeType != AutoScalerServerNodeAutoscaled && vm.NodeType != AutoScalerServerNodeManaged {

		err = fmt.Errorf(constantes.ErrVMNotProvisionnedByMe, vm.NodeName)

	} else if state, err = vm.statusVM(); err != nil {

		err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, err)

	} else if state == AutoScalerServerNodeStateRunning {
		if err = c.CordonNode(vm.NodeName); err != nil {
			glog.Errorf(constantes.ErrCordonNodeReturnError, vm.NodeName, err)
		}

		if err = config.PowerOff(vm.VMUUID, "soft"); err == nil {
			vm.State = AutoScalerServerNodeStateStopped
		} else if err = config.WaitForPowerState(vm.VMUUID, false); err != nil {
			err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, err)
		} else {
			err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, err)
		}

	} else if state != AutoScalerServerNodeStateStopped {

		err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, fmt.Sprintf("Unexpected state: %d", state))

	}

	if err == nil {
		glog.Infof("Stopped VM:%s", vm.NodeName)
	} else {
		glog.Errorf("Could not stop VM:%s. Reason: %s", vm.NodeName, err)
	}

	return err
}

func (vm *AutoScalerServerNode) deleteVM(c types.ClientGenerator) error {
	glog.Debugf("AutoScalerNode::deleteVM, node:%s", vm.NodeName)

	var err error
	var status *desktop.Status

	if vm.NodeType != AutoScalerServerNodeAutoscaled && vm.NodeType != AutoScalerServerNodeManaged {
		err = fmt.Errorf(constantes.ErrVMNotProvisionnedByMe, vm.NodeName)
	} else {
		config := vm.Configuration

		if status, err = config.Status(vm.VMUUID); err == nil {
			if status.Powered {
				// Delete kubernetes node only is alive
				if _, err = c.GetNode(vm.NodeName); err == nil {
					if err = c.MarkDrainNode(vm.NodeName); err != nil {
						glog.Errorf(constantes.ErrCordonNodeReturnError, vm.NodeName, err)
					}

					if err = c.DrainNode(vm.NodeName, true, true); err != nil {
						glog.Errorf(constantes.ErrDrainNodeReturnError, vm.NodeName, err)
					}

					if err = c.DeleteNode(vm.NodeName); err != nil {
						glog.Errorf(constantes.ErrDeleteNodeReturnError, vm.NodeName, err)
					}
				}

				if err = config.PowerOff(vm.VMUUID, "soft"); err != nil {
					err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, err)
				} else if err = config.WaitForPowerState(vm.VMUUID, false); err != nil {
					err = fmt.Errorf(constantes.ErrStopVMFailed, vm.NodeName, err)
				} else {
					vm.State = AutoScalerServerNodeStateStopped

					if err = config.Delete(vm.VMUUID); err != nil {
						err = fmt.Errorf(constantes.ErrDeleteVMFailed, vm.NodeName, err)
					}
				}
			} else if err = config.Delete(vm.NodeName); err != nil {
				err = fmt.Errorf(constantes.ErrDeleteVMFailed, vm.NodeName, err)
			}
		}
	}

	if err == nil {
		glog.Infof("Deleted VM:%s", vm.NodeName)
		vm.State = AutoScalerServerNodeStateDeleted
	} else {
		glog.Errorf("Could not delete VM:%s. Reason: %s", vm.NodeName, err)
	}

	return err
}

func (vm *AutoScalerServerNode) statusVM() (AutoScalerServerNodeState, error) {
	glog.Debugf("AutoScalerNode::statusVM, node:%s", vm.NodeName)

	// Get VM infos
	var status *desktop.Status
	var err error

	if vm.VMUUID == "" {
		return AutoScalerServerNodeStateUndefined, fmt.Errorf(constantes.ErrGetVMInfoFailed, vm.NodeName, "vm identifer is empty")
	}

	config := vm.Configuration
	if status, err = config.Status(vm.VMUUID); err != nil {
		glog.Errorf(constantes.ErrGetVMInfoFailed, vm.NodeName, err)
		return AutoScalerServerNodeStateUndefined, err
	}

	if status != nil {
		vm.IPAddress = config.FindPreferredIPAddress(status.Ethernet)

		if status.Powered {
			vm.State = AutoScalerServerNodeStateRunning
		} else {
			vm.State = AutoScalerServerNodeStateStopped
		}

		return vm.State, nil
	}

	return AutoScalerServerNodeStateUndefined, fmt.Errorf(constantes.ErrAutoScalerInfoNotFound, vm.NodeName)
}

// GetConfiguration method
func (vm *AutoScalerServerNode) GetConfiguration() *desktop.Configuration {
	return vm.Configuration
}

func (vm *AutoScalerServerNode) findInstanceUUID() string {
	if vmUUID, err := vm.Configuration.UUID(vm.NodeName); err == nil {
		vm.VMUUID = vmUUID

		return vmUUID
	}

	return ""
}

func (vm *AutoScalerServerNode) setProviderID(c types.ClientGenerator) error {
	if vm.serverConfig.UseK3S != nil && *vm.serverConfig.UseK3S {
		return nil
	}

	return c.SetProviderID(vm.NodeName, vm.generateProviderID())
}

func (vm *AutoScalerServerNode) generateProviderID() string {
	if vm.serverConfig.UseK3S != nil && *vm.serverConfig.UseK3S {
		return fmt.Sprintf("k3s://%s", vm.NodeName)
	} else {
		return fmt.Sprintf("desktop://%s", vm.VMUUID)
	}
}

func (vm *AutoScalerServerNode) setServerConfiguration(config *types.AutoScalerServerConfig) {
	vm.Configuration.Network.UpdateMacAddressTable(vm.NodeIndex)
	vm.serverConfig = config
}

func (vm *AutoScalerServerNode) retrieveNetworkInfos() error {
	return vm.Configuration.RetrieveNetworkInfos(vm.VMUUID, vm.NodeIndex)
}
