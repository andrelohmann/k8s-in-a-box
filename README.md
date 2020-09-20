# k8s-in-a-box
Raspberry Pi K8s Cluster

## Hardware

4 x Raspberry Pi 4 2GB

3 x Raspberry Pi 4 8GB

7 x SanDisk microSD 32GB U3 A1

3 x 500GB HDD USB3

1 x 8 Port Gigabit Switch

1 x GL-MT300N-V2 Mini Smart Router

https://openwrt.org/toh/hwdata/gl.inet/gl.inet_gl-mt300n_v2

## Setup

GL-MT300N-V2 (router) <- connect via WLAN-Client to your public WLAN, offer WLAN access Point, offer LAN

IP Addresses

Router - 192.168.199.1
HA Proxy - 192.168.199.10
Master Node1 - 192.168.199.11
Master Node2 - 192.168.199.12
Master Node3 - 192.168.199.13
Worker Node1 - 192.168.199.21
Worker Node2 - 192.168.199.22
Worker Node3 - 192.168.199.23



## Installation

### OS Preparation

Flash SD Cards with Raspberry Pi Imager (https://www.raspberrypi.org/downloads/). Choose Ubuntu 20.04.1 LTS for arm64.

Edit network-config on the SD Card.

```
version: 2
ethernets:
  eth0:
    addresses:
      - 192.168.199.10/24
    dhcp4: no
    gateway4: 192.168.199.1
    nameservers:
        addresses: [192.168.199.1]
```

### Kubernetes

  * https://phoenixnap.com/kb/how-to-install-kubernetes-on-a-bare-metal-server
  * https://opensource.com/article/20/6/kubernetes-raspberry-pi
  * https://www.serverlab.ca/tutorials/containers/kubernetes/deploying-kubernetes-ubuntu-18-04/
  * https://thenewstack.io/how-to-deploy-a-kubernetes-cluster-with-ubuntu-server-18-04/
  * https://www.youtube.com/watch?v=qv3_gLvjITk

