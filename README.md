[![Build Status](https://github.com/Fred78290/kubernetes-desktop-autoscaler/actions/workflows/tag.yml/badge.svg?branch=autoscaler-v1.26)](https://github.com/Fred78290/kubernetes-desktop-autoscaler/actions/workflows/tag.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Fred78290_kubernetes-desktop-autoscaler&metric=alert_status)](https://sonarcloud.io/dashboard?id=Fred78290_kubernetes-desktop-autoscaler)
[![Licence](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/Fred78290/kubernetes-desktop-autoscaler/blob/master/LICENSE)

# kubernetes-desktop-autoscaler

Kubernetes autoscaler for VMWare Fusion 13 or VMWare workstation including a custom resource controller to create managed node without code

It use [vmware-desktop-autoscaler-utility](https://github.com/Fred78290/vmware-desktop-autoscaler-utility) to pilot VMWare workstation or VMWare Fusion.
### Supported releases ###

* 1.26.6
    - This version is supported kubernetes v1.26 and support k3s

* 1.27.3
    - This version is supported kubernetes v1.27 and support k3s

## How it works

This tool will drive vSphere to deploy VM at the demand. The cluster autoscaler deployment use vanilla cluster-autoscaler or my enhanced version of [cluster-autoscaler](https://github.com/Fred78290/autoscaler).

This version use grpc to communicate with the cloud provider hosted outside the pod. A docker image is available here [cluster-autoscaler](https://hub.docker.com/r/fred78290/cluster-autoscaler)

A sample of the cluster-autoscaler deployment is available at [examples/cluster-autoscaler.yaml](./examples/cluster-autoscaler.yaml). You must fill value between <>

### Before you must create a kubernetes cluster on VMWare Fusion or Workstation

You can do it from scrash or you can use script from project [autoscaled-masterkube-desktop](https://github.com/Fred78290/autoscaled-masterkube-desktop) to create a kubernetes cluster in single control plane or in HA mode with 3 control planes.

## Commandline arguments

| Parameter | Description |
| --- | --- |
| `version` | Print the version and exit  |
| `save`  | Tell the tool to save state in this file  |
| `config`  |The the tool to use config file |

## Build

The build process use make file. The simplest way to build is `make container`

# New features

## Use k3s

Instead using **kubeadm** as kubernetes deployment tool, it is possible to use **k3s**
## Use the vanilla autoscaler with extern gRPC cloud provider

You can also use the vanilla autoscaler with the [externalgrpc cloud provider](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/externalgrpc)

A sample of the cluster-autoscaler deployment with vanilla autoscaler is available at [examples/cluster-autoscaler-vanilla.yaml](./examples/cluster-autoscaler-vanilla.yaml). You must fill value between <>

## Network

Now it's possible to disable dhcp-default routes and custom route


## CRD controller

This new release include a CRD controller allowing to create kubernetes node without use of code. Just by apply a configuration file, you have the ability to create nodes on the fly.

As exemple you can take a look on [artifacts/examples/example.yaml](artifacts/examples/example.yaml) on execute the following command to create a new node

```
kubectl apply -f artifacts/examples/example.yaml
```

If you want delete the node just delete the CRD with the call

```
kubectl delete -f artifacts/examples/example.yaml
```

You have the ability also to create a control plane as instead a worker

```
kubectl apply -f artifacts/examples/controlplane.yaml
```

The resource is cluster scope so you don't need a namespace. The name of the resource is not the name of the managed node.

The minimal resource declaration

```
apiVersion: "nodemanager.aldunelabs.com/v1alpha1"
kind: "ManagedNode"
metadata:
  name: "desktop-imac-k3s-managed-01"
spec:
  nodegroup: desktop-imac-k3s
  vcpus: 2
  memorySizeInMb: 2048
  diskSizeInMb: 10240
```

The full qualified resource including networks declaration to override the default controller network management and adding some node labels & annotations. If you specify the managed node as controller, you can also allows the controlplane to support deployment as a worker node

```
apiVersion: "nodemanager.aldunelabs.com/v1alpha1"
kind: "ManagedNode"
metadata:
  name: "desktop-imac-k3s-managed-01"
spec:
  nodegroup: desktop-imac-k3s
  controlPlane: false
  allowDeployment: false
  vcpus: 2
  memorySizeInMb: 2048
  diskSizeInMb: 10240
  labels:
  - demo-label.acme.com=demo
  - sample-label.acme.com=sample
  annotations:
  - demo-annotation.acme.com=demo
  - sample-annotation.acme.com=sample
  networks:
    -
      vnet: vmnet8
      device: vmxnet3
      address: 192.168.172.80
      netmask: 255.255.255.0
      gateway: 192.168.172.2
    -
      vnet: vmnet1
      device: vmxnet3
      address: 172.16.251.80
      netmask: 255.255.255.0
```

# Declare additional routes and disable default DHCP routes

The release 1.24 and above allows to add additionnal route per interface, it also allows to disable default route declared by DHCP server.

As example of use generated by autoscaled-masterkube-vmware scripts

```
{
  "use-external-etcd": false,
  "src-etcd-ssl-dir": "/etc/etcd/ssl",
  "dst-etcd-ssl-dir": "/etc/kubernetes/pki/etcd",
  "use-k3s": true,
  "kubernetes-pki-srcdir": "/etc/kubernetes/pki",
  "kubernetes-pki-dstdir": "/etc/kubernetes/pki",
  "network": "unix",
  "listen": "/var/run/cluster-autoscaler/vmware.sock",
  "secret": "desktop",
  "minNode": 0,
  "maxNode": 9,
  "maxNode-per-cycle": 2,
  "node-name-prefix": "autoscaled",
  "managed-name-prefix": "managed",
  "controlplane-name-prefix": "master",
  "nodePrice": 0,
  "podPrice": 0,
  "optionals": {
    "pricing": false,
    "getAvailableMachineTypes": false,
    "newNodeGroup": false,
    "templateNodeInfo": false,
    "createNodeGroup": false,
    "deleteNodeGroup": false
  },
  "kubeadm": {
    "address": "172.16.216.20:6443",
    "token": "K10bf745c67aaceb81bddfc13d2a3029b75dc9f88b95356f05b262bc3f899a5f29e::server:5868d627685b30ae4d5eb4b18d810a31",
    "ca": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "extras-args": [
      "--ignore-preflight-errors=All"
    ]
  },
  "k3s": {
    "datastore-endpoint": "",
    "extras-commands": []
  },
  "default-machine": "medium",
  "machines": {
    "tiny": {
      "memsize": 1024,
      "vcpus": 1,
      "disksize": 10240
    },
    "small": {
      "memsize": 2048,
      "vcpus": 2,
      "disksize": 20480
    },
    "medium": {
      "memsize": 3072,
      "vcpus": 2,
      "disksize": 20480
    },
    "large": {
      "memsize": 4096,
      "vcpus": 4,
      "disksize": 51200
    },
    "xlarge": {
      "memsize": 16384,
      "vcpus": 4,
      "disksize": 102400
    },
    "2xlarge": {
      "memsize": 16384,
      "vcpus": 8,
      "disksize": 102400
    },
    "4xlarge": {
      "memsize": 32768,
      "vcpus": 8,
      "disksize": 102400
    }
  },
  "node-labels": [
    "topology.kubernetes.io/region=home",
    "topology.kubernetes.io/zone=office",
    "topology.csi.vmware.com/k8s-region=home",
    "topology.csi.vmware.com/k8s-zone=office"
  ],
  "cloud-init": {
    "package_update": false,
    "package_upgrade": false,
    "growpart": {
      "mode": "auto",
      "devices": [
        "/"
      ],
      "ignore_growroot_disabled": false
    },
    "runcmd": [
      "echo '172.16.216.20 desktop-imac-k3s-masterkube desktop-imac-k3s-masterkube.aldunelabs.com' >> /etc/hosts"
    ]
  },
  "ssh-infos": {
    "wait-ssh-ready-seconds": 180,
    "user": "kubernetes",
    "ssh-private-key": "/etc/ssh/id_rsa"
  },
  "vmware": {
    "api": {
      "address": "10.0.0.53:5323",
      "key": "/etc/ssl/certs/autoscaler-utility/client.key",
      "cert": "/etc/ssl/certs/autoscaler-utility/client.crt",
      "cacert": "/etc/ssl/certs/autoscaler-utility/ca.crt"
    },
    "config": {
      "nodegroup": "desktop-imac-k3s",
      "timeout": 300,
      "template": "1QMNNKCEH6HVNBH4AJ13HARQG7AS0EJJ",
      "linked": false,
      "autostart": true,
      "network": {
        "domain": "aldunelabs.com",
        "dns": {
          "search": [
            "aldunelabs.com"
          ],
          "nameserver": [
            "10.0.0.5"
          ]
        },
        "interfaces": [
          {
            "primary": false,
            "exists": true,
            "vnet": "vmnet8",
            "type": "nat",
            "device": "vmxnet3",
            "mac-address": "generate",
            "nic": "eth0",
            "dhcp": true,
            "use-dhcp-routes": true,
            "routes": []
          },
          {
            "primary": true,
            "exists": true,
            "vnet": "vmnet1",
            "type": "hostOnly",
            "device": "vmxnet3",
            "mac-address": "generate",
            "nic": "eth1",
            "dhcp": false,
            "use-dhcp-routes": false,
            "address": "172.16.216.24",
            "gateway": "172.16.74.2",
            "netmask": "255.255.255.0",
            "routes": []
          }
        ]
      }
    }
  }
}
```