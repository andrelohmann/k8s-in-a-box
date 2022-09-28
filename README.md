# k8s-in-a-box

Raspberry Pi K8s Cluster

<div align="center">
  <img src="cluster.jpg" />
</div>

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

IP Addresses:
  * Router - 192.168.199.1
  * HA Proxy - 192.168.199.10
  * Master Node1 - 192.168.199.11
  * Master Node2 - 192.168.199.12
  * Master Node3 - 192.168.199.13
  * Worker Node1 - 192.168.199.21
  * Worker Node2 - 192.168.199.22
  * Worker Node3 - 192.168.199.23
 
### Router Steps

  * Install Open-WRT
  * Connect via cable with the right port (left is still configured as WAN port)
  * 192.168.1.1 <- root : admin
  * Network -> Interfaces -> delete both WAN ports
  * Network -> Interfaces -> Devices -> edit br-lan -> add eth0.2 to bridge ports
  * Network -> Wireless -> radio0 -> Scan -> Select network and connect to it
  * Configure a second SSID as Master
  * !you can only have one client active!
  * If the client can't connect, the AP (master) SSID will not be provided
  * Everytime you disconnect the cluster, disable the client network first, otherwise next time, you are not able to connect to the routers AP (master)
  
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

Edit meta-data on the SD Card.

```
instance_id: proxy.k8s.lan
```

Edit user-data on the SD Card.

```
# On first boot, set the (default) ubuntu user's password to "ubuntu" and
# expire user passwords
chpasswd:
  #expire: true
  list:
  - ubuntu:ubuntu
users:
- name: deploy
  gecos: K8s Deployment User
  sudo: ALL=(ALL) NOPASSWD:ALL
  groups: users, admin
  ssh_import_id: None
  lock_passwd: true
  ssh_authorized_keys:
  - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCpSyhHsctydITvC2XeqLYOZtbkKhOAf3f9sy8MYMpFcKQ2CRJ5DVMgRJyUR6yYLqlMTxZW7i9UaB0r+Bzgis3ay3N7EubJgZDPSNe3RyVvS1EShahEZeijZL0mhU4xq8Ui/LGjpOhGEtSCV/5CIqxPINlpVKlXfHgInJvEYA+hY6rns8+x8shq9KYb/Frpj2DftgZJoEVfEgFxrUIaiZA68KKPMVdTxL5B4xmBIofPvbhEZbQCsysjjJGTA+SUCSC3rKyM1UsqGntz+oftd5HN2XNfCDNVKFLOkKTRIfPZc8MrYlC5bB6cx02HYY9fm/1UiOuihZdCklkySQ+B7Igj

# Enable password authentication with the SSH daemon
ssh_pwauth: true
```

### Ansible

Join the vagrant machine.

```
cd vagrant
vagrant up
vagrant ssh
```

Then join the ansible folder and run ansible-playbook.

```
cd ansible
ansible-playbook playbook.yml
```

#### Installation Steps (Ansible Playbooks)

##### K8s Nodes - Prerequesites

Add docker daemon.json

```
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
```

Add k8s.conf

```
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
```

Add cgroup configs to /boot/firmware/cmdline.txt

```
net.ifnames=0 dwc_otg.lpm_enable=0 console=serial0,115200 console=tty1 root=LABEL=writable rootfstype=ext4 elevator=deadline rootwait fixrtc cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1 swapaccount=1
```

Add deb repository

```
deb https://apt.kubernetes.io/ kubernetes-xenial main
```

Install packages and set apt-mark to hold

```
apt install kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl
```

##### Loadbalancer

Install HA Proxy with following configuration

```
# /etc/haproxy/haproxy.cfg
#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
    log /dev/log local0
    log /dev/log local1 notice
    daemon

#---------------------------------------------------------------------
# common defaults that all the 'listen' and 'backend' sections will
# use if not designated in their block
#---------------------------------------------------------------------
defaults
    mode                    http
    log                     global
    option                  httplog
    option                  dontlognull
    option http-server-close
    option forwardfor       except 127.0.0.0/8
    option                  redispatch
    retries                 1
    timeout http-request    10s
    timeout queue           20s
    timeout connect         5s
    timeout client          20s
    timeout server          20s
    timeout http-keep-alive 10s
    timeout check           10s

#---------------------------------------------------------------------
# apiserver frontend which proxys to the masters
#---------------------------------------------------------------------
frontend apiserver
    bind *:6443
    mode tcp
    option tcplog
    default_backend apiserver

#---------------------------------------------------------------------
# round robin balancing for apiserver
#---------------------------------------------------------------------
backend apiserver
    option httpchk GET /healthz
    http-check expect status 200
    mode tcp
    option ssl-hello-chk
    balance     roundrobin
        server master1.k8s.lan 192.168.199.11:6443 check
        server master2.k8s.lan 192.168.199.12:6443 check
        server master3.k8s.lan 192.168.199.13:6443 check

#---------------------------------------------------------------------
# round robin balancing for worker nodes
#---------------------------------------------------------------------
listen worker
    bind *:30000-32767
    mode tcp
    server worker1.k8s.lan 192.168.199.21
    server worker2.k8s.lan 192.168.199.22
    server worker3.k8s.lan 192.168.199.23
```

Testing the loadbalancer

```
nc -v proxy.k8s.lan 6443
```

##### K8s ControlPlane

On the first node, generate a Cluster Token

```
TOKEN=$(sudo kubeadm token generate)
```

Initialize the Controlplane on the first node

```
sudo kubeadm init --token=${TOKEN} --kubernetes-version=v1.18.2 --pod-network-cidr=10.244.0.0/16
```
## Testing the cluster

### General Health

```
kubectl cluster-info
kubectl get nodes -o wide
kubectl get pods --all-namespaces
kubectl config view
kubectl explain pod
kubectl explain service
kubectl explain deployment
```

### Node Failure

```
kubectl get nodes -o wide
kubectl describe nodes worker1.k8s.lan
kubectl cordon worker1.k8s.lan
kubectl drain worker1.k8s.lan
kubectl delete node worker1.k8s.lan
```

### Testing a deployment

#### Deployment

```
kubectl create deployment echo-server --image=ghcr.io/andrelohmann/k8s-in-a-box:latest
kubectl get deployments
kubectl get pods
kubectl get events
```

#### Service

```
kubectl expose deployment echo-server --type=NodePort --port=8000 --name=echo
kubectl get services
kubectl get services echo
```

#### Delete everything

```
kubectl delete service echo
kubectl delete deployment echo-server
```

#### Namespaces

https://kubernetes.io/docs/tasks/administer-cluster/namespaces/

```
kubectl get namespaces
kubectl create -f https://k8s.io/examples/admin/namespace-dev.json
kubectl create -f https://k8s.io/examples/admin/namespace-prod.json
kubectl get namespaces --show-labels
kubectl create deployment echo-server --image=ghcr.io/andrelohmann/k8s-in-a-box:latest  -n=development --replicas=3
kubectl get deployment echo-server -n=development
kubectl expose deployment echo-server --type=NodePort --port=8000 --name=echo -n=development
kubectl get services echo -n=development
kubectl get pods -o wide -n=development
```

## Sources

### Kubernetes

  * https://phoenixnap.com/kb/how-to-install-kubernetes-on-a-bare-metal-server
  * https://opensource.com/article/20/6/kubernetes-raspberry-pi
  * https://www.serverlab.ca/tutorials/containers/kubernetes/deploying-kubernetes-ubuntu-18-04/
  * https://thenewstack.io/how-to-deploy-a-kubernetes-cluster-with-ubuntu-server-18-04/
  * https://www.youtube.com/watch?v=qv3_gLvjITk
  * https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
  * https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/
  * https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/#stacked-control-plane-and-etcd-nodes
  * https://kubernetes.io/docs/concepts/storage/
  * https://linuxconfig.org/how-to-install-kubernetes-on-ubuntu-18-04-bionic-beaver-linux
  * https://bee42.com/de/blog/tutorials/kubernetes-cluster-on-embedded/
  * https://rook.io/
  * https://rook.io/docs/rook/v1.4/ceph-quickstart.html
  * https://kubernetes.io/de/docs/reference/kubectl/cheatsheet/
  * https://medium.com/@chamilad/load-balancing-and-reverse-proxying-for-kubernetes-services-f03dd0efe80
  * https://levelup.gitconnected.com/step-by-step-slow-guide-kubernetes-cluster-on-raspberry-pi-4b-part-3-899fc270600e
  * https://vitobotta.com/2020/03/20/haproxy-kubernetes-hetzner-cloud/

### Docker Security

  * https://www.computerwoche.de/a/7-security-tools-fuer-docker-und-kubernetes,3546931

### Kubernetes Application deployments

  * https://wkrzywiec.medium.com/deployment-of-multiple-apps-on-kubernetes-cluster-walkthrough-e05d37ed63d1
  * https://wkrzywiec.medium.com/how-to-deploy-application-on-kubernetes-with-helm-39f545ad33b8
  * https://medium.com/swlh/how-to-declaratively-run-helm-charts-using-helmfile-ac78572e6088
