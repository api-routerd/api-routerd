REST API based systemd management

[![Build Status CircleCI](https://circleci.com/gh/api-routerd/api-routerd.svg?style=svg)](https://circleci.com/gh/api-routerd/api-routerd)
[![Build Status](https://travis-ci.org/api-routerd/api-routerd.svg?branch=master)](https://travis-ci.org/api-routerd/api-routerd)
[![HitCount](http://hits.dwyl.io/ssahani/api-routerd/api-routerd.svg)](http://hits.dwyl.io/ssahani/api-routerd/api-routerd)
[![CodeFactor](https://www.codefactor.io/repository/github/restgw/api-routerd/badge)](https://www.codefactor.io/repository/github/restgw/api-routerd)
[![Go Report Card](https://goreportcard.com/badge/github.com/api-routerd/api-routerd)](https://goreportcard.com/report/github.com/api-routerd/api-routerd)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Coverage Status](https://coveralls.io/repos/github/api-routerd/api-routerd/badge.svg?branch=master)](https://coveralls.io/github/api-routerd/api-routerd?branch=master)


api-routerd is a cloud-enabled, mobile-ready, a super light weight remote management tool which uses REST API for real time configuration and performance as well as health monitoring for systems (containers) and applications. It provides fast API based monitoring without affecting the system it's running on.

- Proactive Monitoring and Analytics
  api-routerd saves network administrators time and frustration because it makes it easy to gather statistics and perform analyses.
- Platform independent REST APIs can be accessed via any application (curl, chrome, PostMan ...) from any OS (Linux, IOS, Android, Windows ...)
- An [iphone app demo](https://www.linkedin.com/feed/update/urn:li:activity:6493161973260357632) using REST APIS . See [Source Code](https://github.com/api-routerd/iDevOps) .
- Runs on [Rasberry pi](https://www.linkedin.com/feed/update/urn:li:activity:6492312857231888384).
- Minimal data transfer using JSON.
- Plugin based Architechture. See how to write plugin section for more information.

# Features!

|Feature| Details |
| ------ | ------ |
| socket activation | supports systemd socket activation
systemd  | information, services (start, stop, restart, status), service properties for example CPUShares
networkd |config (.network, .netdev, .link)
hostnamed | set hostname
logind |(list-sessions, list-users and terminate-user etc)
timdate| set time, zone
nameserver | add/delete/modify ```/etc/resolv.conf```
timesynd | set configs
systemd-machined | see info about images/machines. start stop machines
journald | ```journald.conf```
systemd conf | ```system.conf```
coredumpd |```coredump.conf```
systemd-resolved |```systemd-resolved.conf```
kernel modules |(modprobe, lsmod, rmmod)
network | via netlink . Link: mtu, up, down, Create bridge and enslave links, Create bond and enslave links, Adddress: Set, Get, Delete, Gateway: Default Gateway Add and Delete
group | add/delete/modify
users |add/delete/modify (requires newuser)
sysctl |add/delete/modify and apply
see information from ```/proc``` fs| netstat, netdev, memory and much more
configure ```/proc/sys/net``` | (core/ipv4/ipv6), VM
ethtool | see information and configure offload features
firewalld | see and configure firewalld
See confs | sudoers and sshd conf


### api-routerd JSON APIs

 Refer spreadsheet [APIs](https://docs.google.com/spreadsheets/d/e/2PACX-1vTl2Vmp-BdTE5Vgi_PiW-qKPJnbLxdSso9kT2GAkAxCu_iWrw3_PZLlEuyXz0lbFgd7DoofXlmmb3dP/pubhtml
)

### Dependencies

api-routerd uses a following open source projects to work properly:

* [logrus](https://github.com/sirupsen/logrus)
* [gorilla mux](https://github.com/gorilla/mux)
* [netlink](https://github.com/vishvananda/netlink)
* [gopsutil](https://github.com/shirou/gopsutil)
* [coreos go-systemd](https://github.com/coreos/go-systemd)
* [dbus](https://github.com/godbus/dbus)
* [ethtool](https://github.com/safchain/ethtool)
* [viper](https://github.com/spf13/viper)
* [go-ini](https://github.com/go-ini/ini)


### Installation

First configure your ```$GOPATH```. If you have already done this skip this step.

```sh
# keep in ~/.bashrc
```

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export OS_OUTPUT_GOPATH=1
```

Clone inside src dir of ```$GOPATH```. In my case

```sh
$ pwd
/home/sus/go/src
```

### Installation

```sh
$ go get github.com/api-routerd/api-routerd
$ git clone https://github.com/api-routerd/api-routerd
$ cd api-routerd
$ pwd
/home/sus/go/src/api-routerd
$ go build -v

$ sudo ./api-routerd
INFO[0000] api-routerd: v0.1 (built go1.11.4)
INFO[0000] Start Server at 0.0.0.0:8080
INFO[0000] Starting api-routerd in plain text mode

```

### How to configure IP and Port ?

Conf dir: ```/etc/api-routerd/```
Conf File: ```api-routerd.toml```

```sh
$ cat /etc/api-routerd/api-routerd.toml
[Network]
IPAddress="0.0.0.0"
Port="8080"
```

### How to configure users ?

Add user name and authentication string in space separated lines

```sh
# cat /etc/api-routerd/api-routerd-auth.conf
Susant secret
Max bbbb
Joy ccccc
```

### How to configure TLS ?

Generate private key (.key)

```sh
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.......................+++++
.+++++
e is 65537 (0x010001)

openssl genrsa -out server.key 2048
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```sh
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:

```

Place ```server.crt``` and ```server.key``` in the dir ```/etc/api-routerd/tls```

```sh
[root@Zeus tls]# ls
server.crt  server.key
[root@Zeus tls]# pwd
/etc/api-routerd/tls

```

Use case: https

```sh
$ curl --header "X-Session-Token: secret" --request GET https://localhost:8080/api/network/ethtool/vmnet8/get-link-features -k --tlsv1.2

```
## Use cases

Refer usecase document [use cases](https://github.com/api-routerd/api-routerd/blob/master/examples.md)

## How to write your own plugin ?

api-routerd is designed with robust plugin based architecture in mind. You can always add and remove modules to it with minimal effort
You can implement and incorporate application features very quickly. Because plug-ins are separate modules with well-defined interfaces,
you can quickly isolate and solve problems. You can create custom versions of an application with minimal source code modifications.

* Choose namespace under cmd directory ( network, proc, system etc) whare you want to put your module.
* Write sub router see for example ```api-routerd/cmd/system/login```
* Write your module ```module.go``` and  ```module_router.go```
* Write ```RegisterRouterModule```
* Register ```RegisterRouterModule``` with parent router for example for ```login``` registered with
  ```RegisterRouterSystem``` under ```system``` namespace as ```login.RegisterRouterLogin```
* See examples directory how to write on your own plugin.

### Todos

 - Write Tests
 - networkd
 - iptables

License
----

Apache 2.0


[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [git-repo-url]: <https://github.com/api-routerd/api-routerd.git>
