package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Fred78290/kubernetes-desktop-autoscaler/api"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/constantes"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/desktop"
	managednodeClientset "github.com/Fred78290/kubernetes-desktop-autoscaler/pkg/generated/clientset/versioned"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/types"
	"github.com/Fred78290/kubernetes-desktop-autoscaler/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	apiv1 "k8s.io/api/core/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const vmnotfound = "vm not found"

var pcislotnumber = []int{160, 192, 161, 193, 225}
var ip4address = []string{"192.168.172.%d", "172.16.251.%d", "192.168.205.%d"}

type MockupEthernetCard struct {
	AddressType          string `json:"addresType"`
	BsdName              string `json:"bsdName"`
	ConnectionType       string `json:"connectionType"`
	DisplayName          string `json:"displayName"`
	MacAddress           string `json:"macaddress"`
	MacAddressOffset     int    `json:"macaddressOffset"`
	LinkStatePropagation bool   `json:"linkStatePropagation"`
	PciSlotNumber        int    `json:"pcislotNumber"`
	Present              bool   `json:"present"`
	VirtualDev           string `json:"virtualDev"`
	Vnet                 string `json:"vnet"`
	IP4Address           string `json:"address"`
}

type MockupVirtualMachine struct {
	Path     string               `json:"path"`
	Uuid     string               `json:"uuid"`
	Name     string               `json:"name"`
	Vcpus    int                  `json:"vcpus"`
	Memory   int                  `json:"memory"`
	Powered  bool                 `json:"powered"`
	Address  string               `json:"address"`
	Ethernet []MockupEthernetCard `json:"ethernet"`
}

type MockupAutoscalerUtility struct {
	VMIdentifiers   []string                         `json:"uuid"`
	VirtualMachines []*MockupVirtualMachine          `json:"vm"`
	vmbyuuid        map[string]*MockupVirtualMachine `json:"-"`
	vmbyname        map[string]*MockupVirtualMachine `json:"-"`
}

func newVMWareDesktopAutoscalerServiceClient(configFile string) (*MockupAutoscalerUtility, error) {
	var config MockupAutoscalerUtility

	if configStr, err := os.ReadFile(configFile); err != nil {
		return nil, err
	} else if err = json.Unmarshal(configStr, &config); err != nil {
		return nil, err
	}

	config.vmbyname = make(map[string]*MockupVirtualMachine)
	config.vmbyuuid = make(map[string]*MockupVirtualMachine)

	for _, vm := range config.VirtualMachines {
		config.vmbyname[vm.Name] = vm
		config.vmbyuuid[vm.Uuid] = vm
	}

	return &config, nil
}

func generateMacAddress() string {
	buf := make([]byte, 3)

	if _, err := rand.Read(buf); err != nil {
		return ""
	}

	return fmt.Sprintf("00:16:3e:%02x:%02x:%02x", buf[0], buf[1], buf[2])
}

func (m *MockupAutoscalerUtility) findInstanceID(nodeName string) string {
	if vm, found := m.vmbyname[nodeName]; found {
		return vm.Uuid
	}

	return m.VMIdentifiers[len(m.VirtualMachines)]
}

func (m *MockupAutoscalerUtility) getMacAddress(address string) string {

	if strings.ToLower(address) == "generate" {
		address = generateMacAddress()
	} else if strings.ToLower(address) == "ignore" {
		address = ""
	}

	return address
}

func (m *MockupAutoscalerUtility) removeVM(vm *MockupVirtualMachine) {
	delete(m.vmbyname, vm.Name)
	delete(m.vmbyuuid, vm.Uuid)

	for index := 0; index < len(m.VirtualMachines); index++ {
		if m.VirtualMachines[index].Uuid == vm.Uuid {
			m.VirtualMachines = append(m.VirtualMachines[:index], m.VirtualMachines[index+1:]...)
		}
	}
}

func (m *MockupAutoscalerUtility) virtualMachineFromVM(vm *MockupVirtualMachine) *api.VirtualMachine {
	return &api.VirtualMachine{
		Uuid:    vm.Uuid,
		Name:    vm.Name,
		Vmx:     vm.Path,
		Vcpus:   int32(vm.Vcpus),
		Memory:  int64(vm.Memory),
		Address: vm.Address,
		Powered: vm.Powered,
	}
}

func (m *MockupAutoscalerUtility) Create(ctx context.Context, in *api.CreateRequest, opts ...grpc.CallOption) (*api.CreateResponse, error) {
	if _, found := m.vmbyname[in.Name]; found {
		return &api.CreateResponse{
			Response: &api.CreateResponse_Error{
				Error: &api.ClientError{
					Code:   409,
					Reason: "Machine already",
				},
			},
		}, nil
	} else {
		uuid := m.VMIdentifiers[len(m.VirtualMachines)]

		vm := &MockupVirtualMachine{
			Name:     in.Name,
			Uuid:     uuid,
			Path:     fmt.Sprintf("/Users/Me/Library/Masterkube/%s.vmwarevm/%s.vmx", uuid, uuid),
			Vcpus:    int(in.Vcpus),
			Memory:   int(in.Memory),
			Powered:  false,
			Ethernet: make([]MockupEthernetCard, 0, len(in.Networks)),
		}

		if len(in.Networks) > 0 {
			nodeIndex := len(m.VirtualMachines)

			for index, net := range in.Networks {
				vm.Ethernet = append(vm.Ethernet, MockupEthernetCard{
					AddressType:          net.Macaddress,
					BsdName:              net.BsdName,
					ConnectionType:       net.Type,
					DisplayName:          net.DisplayName,
					MacAddress:           m.getMacAddress(net.Macaddress),
					MacAddressOffset:     0,
					Vnet:                 net.Vnet,
					IP4Address:           fmt.Sprintf(ip4address[index], nodeIndex+10),
					LinkStatePropagation: true,
					PciSlotNumber:        pcislotnumber[index],
				})
			}

			vm.Address = vm.Ethernet[0].IP4Address
		}

		m.VirtualMachines = append(m.VirtualMachines, vm)
		m.vmbyname[vm.Name] = vm
		m.vmbyuuid[vm.Uuid] = vm

		return &api.CreateResponse{Response: &api.CreateResponse_Result{
			Result: &api.CreateReply{
				Machine: m.virtualMachineFromVM(vm),
			},
		}}, nil
	}
}

func (m *MockupAutoscalerUtility) Delete(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.DeleteResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.DeleteResponse{
			Response: &api.DeleteResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else if vm.Powered {
		return &api.DeleteResponse{
			Response: &api.DeleteResponse_Error{
				Error: &api.ClientError{
					Code:   409,
					Reason: "vm powered",
				},
			},
		}, nil

	} else {
		m.removeVM(vm)
	}

	return &api.DeleteResponse{
		Response: &api.DeleteResponse_Result{
			Result: &api.DeleteReply{
				Done: true,
			},
		},
	}, nil
}

func (m *MockupAutoscalerUtility) PowerOn(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.PowerOnResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.PowerOnResponse{
			Response: &api.PowerOnResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else if vm.Powered {
		return &api.PowerOnResponse{
			Response: &api.PowerOnResponse_Error{
				Error: &api.ClientError{
					Code:   409,
					Reason: "already powered",
				},
			},
		}, nil
	} else {
		vm.Powered = true
	}

	return &api.PowerOnResponse{
		Response: &api.PowerOnResponse_Result{
			Result: &api.PowerOnReply{
				Done: true,
			},
		},
	}, nil
}

func (m *MockupAutoscalerUtility) PowerOff(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.PowerOffResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.PowerOffResponse{
			Response: &api.PowerOffResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else if !vm.Powered {
		return &api.PowerOffResponse{
			Response: &api.PowerOffResponse_Error{
				Error: &api.ClientError{
					Code:   409,
					Reason: "already power off",
				},
			},
		}, nil
	} else {
		vm.Powered = false
	}

	return &api.PowerOffResponse{
		Response: &api.PowerOffResponse_Result{
			Result: &api.PowerOffReply{
				Done: true,
			},
		},
	}, nil
}

func (m *MockupAutoscalerUtility) PowerState(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.PowerStateResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.PowerStateResponse{
			Response: &api.PowerStateResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.PowerStateResponse{
			Response: &api.PowerStateResponse_Powered{
				Powered: vm.Powered,
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) ShutdownGuest(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.ShutdownGuestResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.ShutdownGuestResponse{
			Response: &api.ShutdownGuestResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else if !vm.Powered {
		return &api.ShutdownGuestResponse{
			Response: &api.ShutdownGuestResponse_Error{
				Error: &api.ClientError{
					Code:   409,
					Reason: "vm not powered",
				},
			},
		}, nil
	} else {
		vm.Powered = false
	}

	return &api.ShutdownGuestResponse{
		Response: &api.ShutdownGuestResponse_Result{
			Result: &api.PowerOffReply{
				Done: true,
			},
		},
	}, nil
}

func (m *MockupAutoscalerUtility) Status(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.StatusResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.StatusResponse{
			Response: &api.StatusResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		ethernet := make([]*api.Ethernet, 0, len(vm.Ethernet))

		for _, net := range vm.Ethernet {
			ethernet = append(ethernet, &api.Ethernet{
				GeneratedAddress: net.MacAddress,
				Address:          net.IP4Address,
				AddressType:      net.AddressType,
				BsdName:          net.BsdName,
				ConnectionType:   net.ConnectionType,
				DisplayName:      net.DisplayName,
				PciSlotNumber:    int32(net.PciSlotNumber),
				Present:          net.Present,
				VirtualDev:       net.VirtualDev,
				Vnet:             net.Vnet,
			})
		}

		return &api.StatusResponse{
			Response: &api.StatusResponse_Result{
				Result: &api.StatusReply{
					Powered:  vm.Powered,
					Ethernet: ethernet,
				},
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) WaitForIP(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.WaitForIPResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.WaitForIPResponse{
			Response: &api.WaitForIPResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.WaitForIPResponse{
			Response: &api.WaitForIPResponse_Result{
				Result: &api.WaitForIPReply{
					Address: vm.Address,
				},
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) WaitForToolsRunning(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.WaitForToolsRunningResponse, error) {
	if _, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.WaitForToolsRunningResponse{
			Response: &api.WaitForToolsRunningResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.WaitForToolsRunningResponse{
			Response: &api.WaitForToolsRunningResponse_Result{
				Result: &api.WaitForToolsRunningReply{
					Running: true,
				},
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) SetAutoStart(ctx context.Context, in *api.AutoStartRequest, opts ...grpc.CallOption) (*api.AutoStartResponse, error) {
	if _, found := m.vmbyuuid[in.Uuid]; !found {
		return &api.AutoStartResponse{
			Response: &api.AutoStartResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.AutoStartResponse{
			Response: &api.AutoStartResponse_Result{
				Result: &api.AutoStartReply{
					Done: true,
				},
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) VirtualMachineByName(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.VirtualMachineResponse, error) {
	if vm, found := m.vmbyname[in.Identifier]; !found {
		return &api.VirtualMachineResponse{
			Response: &api.VirtualMachineResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.VirtualMachineResponse{
			Response: &api.VirtualMachineResponse_Result{
				Result: m.virtualMachineFromVM(vm),
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) VirtualMachineByUUID(ctx context.Context, in *api.VirtualMachineRequest, opts ...grpc.CallOption) (*api.VirtualMachineResponse, error) {
	if vm, found := m.vmbyuuid[in.Identifier]; !found {
		return &api.VirtualMachineResponse{
			Response: &api.VirtualMachineResponse_Error{
				Error: &api.ClientError{
					Code:   404,
					Reason: vmnotfound,
				},
			},
		}, nil
	} else {
		return &api.VirtualMachineResponse{
			Response: &api.VirtualMachineResponse_Result{
				Result: m.virtualMachineFromVM(vm),
			},
		}, nil
	}
}

func (m *MockupAutoscalerUtility) ListVirtualMachines(ctx context.Context, in *api.VirtualMachinesRequest, opts ...grpc.CallOption) (*api.VirtualMachinesResponse, error) {
	machines := make([]*api.VirtualMachine, 0, len(m.VirtualMachines))

	for _, vm := range m.VirtualMachines {
		machines = append(machines, m.virtualMachineFromVM(vm))
	}

	return &api.VirtualMachinesResponse{
		Response: &api.VirtualMachinesResponse_Result{
			Result: &api.VirtualMachinesReply{
				Machines: machines,
			},
		},
	}, nil
}

type baseTest struct {
	mockup     *MockupAutoscalerUtility
	testConfig *desktop.Configuration
	t          *testing.T
}

type nodegroupTest struct {
	baseTest
}

type autoScalerServerNodeGroupTest struct {
	AutoScalerServerNodeGroup
	baseTest
}

func (ng *autoScalerServerNodeGroupTest) createTestNode(nodeName string, desiredState ...AutoScalerServerNodeState) *AutoScalerServerNode {
	var state AutoScalerServerNodeState = AutoScalerServerNodeStateNotCreated

	if len(desiredState) > 0 {
		state = desiredState[0]
	}

	node := &AutoScalerServerNode{
		NodeGroupID:   testGroupID,
		NodeName:      nodeName,
		VMUUID:        ng.findInstanceID(nodeName),
		CRDUID:        testCRDUID,
		Memory:        ng.Machine.Memory,
		CPU:           ng.Machine.Vcpu,
		Disk:          ng.Machine.Disk,
		IPAddress:     "127.0.0.1",
		State:         state,
		NodeType:      AutoScalerServerNodeAutoscaled,
		NodeIndex:     1,
		Configuration: ng.testConfig,
		serverConfig:  ng.configuration,
	}

	if vmuuid := node.findInstanceUUID(); len(vmuuid) > 0 {
		node.VMUUID = vmuuid
	}

	ng.Nodes[nodeName] = node
	ng.RunningNodes[len(ng.RunningNodes)+1] = ServerNodeStateRunning

	return node
}

func (m *nodegroupTest) launchVM() {
	ng, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if err := testNode.launchVM(m, ng.NodeLabels, ng.SystemLabels); err != nil {
			m.t.Errorf("AutoScalerNode.launchVM() error = %v", err)
		}
	}
}

func (m *nodegroupTest) startVM() {
	_, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if err := testNode.startVM(m); err != nil {
			m.t.Errorf("AutoScalerNode.startVM() error = %v", err)
		}
	}
}

func (m *nodegroupTest) stopVM() {
	_, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if err := testNode.stopVM(m); err != nil {
			m.t.Errorf("AutoScalerNode.stopVM() error = %v", err)
		}
	}
}

func (m *nodegroupTest) deleteVM() {
	_, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if err := testNode.deleteVM(m); err != nil {
			m.t.Errorf("AutoScalerNode.deleteVM() error = %v", err)
		}
	}
}

func (m *nodegroupTest) statusVM() {
	_, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if got, err := testNode.statusVM(); err != nil {
			m.t.Errorf("AutoScalerNode.statusVM() error = %v", err)
		} else if got != AutoScalerServerNodeStateRunning {
			m.t.Errorf("AutoScalerNode.statusVM() = %v, want %v", got, AutoScalerServerNodeStateRunning)
		}
	}
}

func (m *nodegroupTest) addNode() {
	ng, err := m.newTestNodeGroup()

	if assert.NoError(m.t, err) {
		if _, err := ng.addNodes(m, 1); err != nil {
			m.t.Errorf("AutoScalerServerNodeGroup.addNode() error = %v", err)
		}
	}
}

func (m *nodegroupTest) deleteNode() {
	ng, testNode, err := m.newTestNode(launchVMName)

	if assert.NoError(m.t, err) {
		if err := ng.deleteNodeByName(m, testNode.NodeName); err != nil {
			m.t.Errorf("AutoScalerServerNodeGroup.deleteNode() error = %v", err)
		}
	}
}

func (m *nodegroupTest) deleteNodeGroup() {
	ng, err := m.newTestNodeGroup()

	if assert.NoError(m.t, err) {
		if err := ng.deleteNodeGroup(m); err != nil {
			m.t.Errorf("AutoScalerServerNodeGroup.deleteNodeGroup() error = %v", err)
		}
	}
}

func (m *baseTest) KubeClient() (kubernetes.Interface, error) {
	return nil, nil
}

func (m *baseTest) NodeManagerClient() (managednodeClientset.Interface, error) {
	return nil, nil
}

func (m *baseTest) ApiExtentionClient() (apiextension.Interface, error) {
	return nil, nil
}

func (m *baseTest) PodList(nodeName string, podFilter types.PodFilterFunc) ([]apiv1.Pod, error) {
	return nil, nil
}

func (m *baseTest) NodeList() (*apiv1.NodeList, error) {
	node := apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNodeName,
			UID:  testCRDUID,
			Annotations: map[string]string{
				constantes.AnnotationNodeGroupName:        testGroupID,
				constantes.AnnotationNodeIndex:            "0",
				constantes.AnnotationInstanceID:           m.findInstanceID(testNodeName),
				constantes.AnnotationNodeAutoProvisionned: "true",
				constantes.AnnotationScaleDownDisabled:    "false",
				constantes.AnnotationNodeManaged:          "false",
			},
		},
	}

	return &apiv1.NodeList{
		Items: []apiv1.Node{
			node,
		},
	}, nil
}

func (m *baseTest) UncordonNode(nodeName string) error {
	return nil
}

func (m *baseTest) CordonNode(nodeName string) error {
	return nil
}

func (m *baseTest) SetProviderID(nodeName, providerID string) error {
	return nil
}

func (m *baseTest) MarkDrainNode(nodeName string) error {
	return nil
}

func (m *baseTest) DrainNode(nodeName string, ignoreDaemonSet, deleteLocalData bool) error {
	return nil
}

func (m *baseTest) GetNode(nodeName string) (*apiv1.Node, error) {
	node := &apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
			UID:  testCRDUID,
			Annotations: map[string]string{
				constantes.AnnotationNodeGroupName:        testGroupID,
				constantes.AnnotationNodeIndex:            "0",
				constantes.AnnotationInstanceID:           m.findInstanceID(nodeName),
				constantes.AnnotationNodeAutoProvisionned: "true",
				constantes.AnnotationScaleDownDisabled:    "false",
				constantes.AnnotationNodeManaged:          "false",
			},
		},
	}

	return node, nil
}

func (m *baseTest) DeleteNode(nodeName string) error {
	return nil
}

func (m *baseTest) AnnoteNode(nodeName string, annotations map[string]string) error {
	return nil
}

func (m *baseTest) LabelNode(nodeName string, labels map[string]string) error {
	return nil
}

func (m *baseTest) TaintNode(nodeName string, taints ...apiv1.Taint) error {
	return nil
}

func (m *baseTest) WaitNodeToBeReady(nodeName string) error {
	return nil
}

func (m *baseTest) newTestNodeNamedWithState(nodeName string, state AutoScalerServerNodeState) (*autoScalerServerNodeGroupTest, *AutoScalerServerNode, error) {

	if ng, err := m.newTestNodeGroup(); err == nil {
		vm := ng.createTestNode(nodeName, state)

		return ng, vm, err
	} else {
		return nil, nil, err
	}
}

func (m *baseTest) newTestNode(name ...string) (*autoScalerServerNodeGroupTest, *AutoScalerServerNode, error) {
	nodeName := testNodeName

	if len(name) > 0 {
		nodeName = name[0]
	}

	return m.newTestNodeNamedWithState(nodeName, AutoScalerServerNodeStateNotCreated)
}

func (m *baseTest) newTestNodeGroup() (*autoScalerServerNodeGroupTest, error) {
	config, err := m.newTestConfig()

	if err == nil {
		if machine, ok := config.Machines[config.DefaultMachineType]; ok {
			ng := &autoScalerServerNodeGroupTest{
				baseTest: baseTest{
					mockup:     m.mockup,
					t:          m.t,
					testConfig: m.testConfig,
				},
				AutoScalerServerNodeGroup: AutoScalerServerNodeGroup{
					AutoProvision:              true,
					ServiceIdentifier:          config.ServiceIdentifier,
					NodeGroupIdentifier:        testGroupID,
					ProvisionnedNodeNamePrefix: config.ProvisionnedNodeNamePrefix,
					ManagedNodeNamePrefix:      config.ManagedNodeNamePrefix,
					ControlPlaneNamePrefix:     config.ControlPlaneNamePrefix,
					Status:                     NodegroupCreated,
					MinNodeSize:                config.MinNode,
					MaxNodeSize:                config.MaxNode,
					SystemLabels:               types.KubernetesLabel{},
					Nodes:                      make(map[string]*AutoScalerServerNode),
					RunningNodes:               make(map[int]ServerNodeState),
					pendingNodes:               make(map[string]*AutoScalerServerNode),
					configuration:              config,
					Machine:                    machine,
					NodeLabels:                 config.NodeLabels,
				},
			}

			return ng, err
		}

		m.t.Fatalf("Unable to find machine definition for type: %s", config.DefaultMachineType)
	}

	return nil, err
}

func (m *baseTest) getConfFile() string {
	if config := os.Getenv("TEST_CONFIG"); config != "" {
		return config
	}

	return "../test/config.json"
}

func (m *baseTest) getClientConfFile() string {
	if config := os.Getenv("TEST_CLIENT_CONFIG"); config != "" {
		return config
	}

	return "../test/mockup.json"
}

var phConfig *types.AutoScalerServerConfig

func (m *baseTest) newTestConfig() (*types.AutoScalerServerConfig, error) {
	var config types.AutoScalerServerConfig

	if phConfig == nil {

		if configStr, err := os.ReadFile(m.getConfFile()); err != nil {
			return nil, err
		} else if err = json.Unmarshal(configStr, &config); err != nil {
			return nil, err
		} else if mockup, err := newVMWareDesktopAutoscalerServiceClient(m.getClientConfFile()); err != nil {
			return nil, err
		} else {
			m.mockup = mockup
			m.testConfig = config.GetDesktopConfiguration()
			m.testConfig.Timeout = m.testConfig.Timeout * time.Second
			config.SSH.TestMode = true

			m.testConfig.SetClient(mockup)

			phConfig = &config
		}
	}

	return phConfig, nil
}

func (m *baseTest) ssh() {
	config, err := m.newTestConfig()

	if assert.NoError(m.t, err) {
		if _, err = utils.Sudo(config.SSH, "127.0.0.1", 1, "ls"); err != nil {
			m.t.Errorf("SSH error = %v", err)
		}
	}
}

func (m *baseTest) findInstanceID(nodeName string) string {
	return m.mockup.findInstanceID(nodeName)
}

func Test_SSH(t *testing.T) {
	createTestNodegroup(t).ssh()
}

func createTestNodegroup(t *testing.T) *nodegroupTest {
	return &nodegroupTest{
		baseTest: baseTest{
			t: t,
		},
	}
}

func TestNodeGroup_launchVM(t *testing.T) {
	createTestNodegroup(t).launchVM()
}

func TestNodeGroup_startVM(t *testing.T) {
	createTestNodegroup(t).startVM()
}

func TestNodeGroup_stopVM(t *testing.T) {
	createTestNodegroup(t).stopVM()
}

func TestNodeGroup_deleteVM(t *testing.T) {
	createTestNodegroup(t).deleteVM()
}

func TestNodeGroup_statusVM(t *testing.T) {
	createTestNodegroup(t).statusVM()
}

func TestNodeGroupGroup_addNode(t *testing.T) {
	createTestNodegroup(t).addNode()
}

func TestNodeGroupGroup_deleteNode(t *testing.T) {
	createTestNodegroup(t).deleteNode()
}

func TestNodeGroupGroup_deleteNodeGroup(t *testing.T) {
	createTestNodegroup(t).deleteNodeGroup()
}
