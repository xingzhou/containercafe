## OpenRadiant

[![Travis Lint Status](https://travis.innovate.ibm.com/alchemy-containers/openradiant.svg?token=hs5iLEHWzyL9jLf6acy1)](https://travis.innovate.ibm.com/alchemy-containers/openradiant)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

OpenRadiant is an open source containers management and orchestration platform based on best of breed
containers management technologies such as Kubernetes, Swarm and Mesos. OpenRadiant provides automation
to configure Kubernetes and Swarm for a secure and highly available deployment, provides integration
with multiple network SDN solutions and cloud providers, and configurations to support multitenancy.

* [Code of Conduct](#code-of-conduct)
* [Contributing to the project](#contributing-to-the-project)
* [Maintainers](#maintainers)
* [Communication](#communication)
* [Quick Start](#quick-start)
* [Proxy Setup](#proxy-setup)
* [Other Configurations](#custom-configurations)
* [Learn concepts and commands](#learn-concepts-and-commands)
* [License](#license)
* [Issues](#issues)

### Code of Conduct
Participation in the OpenRadiant community is governed by the OpenRadiant [Code of Conduct](CONDUCT.md)

### Contributing to the project
We welcome contributions to the OpenRadiant Project in many forms. There's always plenty to do! Full details of how to contribute to this project are documented in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

### Maintainers
The project's [maintainers](MAINTAINERS.txt): are responsible for reviewing and merging all pull requests and they guide the over-all technical direction of the project.

### Communication
We use \[TODO] [OpenRadiant Slack](https://OpenRadiant.slack.org/) for communication between developers.

### Quick Start

This produces a very simple demonstration cluster of two VirtualBox
VMs, one master and one worker.  They have Mesos installed, and
Kubernetes and Swarm playing nicely together thanks to Mesos.  The
networking is Docker bridge networking.  The Swarm master is modified
for multi-tenant use.

Install [Vagrant](https://www.vagrantup.com/) and
[VirtualBox](https://www.virtualbox.org/wiki/Downloads). You will need
at least Vagrant 1.8.4 and VirtualBox 5.0.24.

Checkout this project:

```
git clone git@github.ibm.com:alchemy-containers/openradiant.git
cd openradiant
```

Install ansible:

```
pip install -r requirements.txt
```

Ansible 1.9.6 is reccomended. If you have a different version and run into issues
try the following:

```
pip install --upgrade ansible==1.9.6
```

Set up your ansible inventory to use the sample project:

```
export ANSIBLE_INVENTORY=../examples/envs/dev-vbox/radiant01.hosts
export ANSIBLE_LIBRARY=$ANSIBLE_INVENTORY
```

Create a new cluster with Vagrant:

```
cd examples/vagrant
vagrant up
cd -
```

Deploy OpenRadiant:

```
cd ansible
ansible-playbook site.yml -e cluster_name=dev-vbox-radiant01
```

Now you can open an SSH connection to the master node:
```
ssh -i ~/.vagrant.d/insecure_private_key vagrant@192.168.10.2
```

On the master you will find both the `kubectl` and `docker`
(currently 1.11) commands on your `$PATH`.

The Swarm master is configured for multi-tenant use.  To prepare to
use it as a tenant, do this on the master:
```
cd; mkdir -p radiant/configs/demo; cat > radiant/configs/demo/config.json <<EOF
{
    "HttpHeaders": {
          "X-Auth-TenantId": "demo"
    }
}
EOF
```

Then you will want these commands on the master:
```
export DOCKER_TLS_VERIFY=""
export DOCKER_CONFIG=~/radiant/configs/demo
export DOCKER_HOST=localhost:2375
```

To get a listing of this tenant's containers, issue the following command on the master:
```
docker ps
```
At first, there will be none.  So create one, like this:
```
docker run --name s1 -d -m 128m busybox sleep 864000
```

Then, get a list of containers, with `docker ps`.  You can inspect its network
configuration from inside, like this:
```
docker exec s1 ifconfig
```

You can create a Kubernetes "deployment" with a command like this:
```
kubectl run k1 --image=busybox sleep 864000
```

The containers in this deployment will be invisible to Swarm because
they lack the label identifying your tenant.  To make containers
visible to Swarm, make a kubernetes pod as follows.  Create a YAML
file prescribing the pod:
```
cat > sleepy-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: sleepy-pod
  annotations:
    containers-label.alpha.kubernetes.io/com.swarm.tenant.0: demo
spec:
  containers:
    - name: sleeper
      image: busybox
      args:
      - sleep
      - "864000"
```

Then create the pod:
```
kubectl create -f sleepy-pod.yaml
```

Then you can watch for it to come up, with
```
kubectl get pod
```

Once it is up, you can inspect its network configuration from inside,
like this:
```
kubectl exec sleepy-pod ifconfig
```

In a similar vein, you can ping one of these containers from the
other.  For example (in which the Swarm container's IP address
is 172.17.0.5):
```
kubectl exec sleepy-pod ping -- -c 2 172.17.0.5
```
For more information see ... (TBD)

####To check the HAproxy statistics using the GUI:

On your local browser, enter the following URL:

master_ip:harproxy_GUI_port/haproxy_stats

Example:
http://192.168.10.2:9000/haproxy_stats (port 9000 is statically assigned)
When prompt for the user_namer:password  use  vagrant:radiantHA

To setup Proxy, follow the steps in [Proxy Setup](#proxy-setup)

### Proxy Setup
For details, please follow [Proxy documentation](proxy/README.md)

### Other Configurations

TBD

### Learn concepts and commands

Browse TBD to learn more. Here are some topics you may be
interested in:

TBD


### License

Copyright 2015-2016 IBM Corporation

Licensed under the [Apache License, Version 2.0 (the "License")](http://www.apache.org/licenses/LICENSE-2.0.html).

Unless required by applicable law or agreed to in writing, software distributed under the license is distributed on an "as is" basis, without warranties or conditions of any kind, either express or implied. See the license for the specific language governing permissions and limitations under the license.

### Issues

Report bugs, ask questions and request features [here on GitHub](../../issues).
