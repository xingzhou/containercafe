## OpenRadiant

[![Travis Lint Status](https://travis.innovate.ibm.com/alchemy-containers/openradiant.svg?token=hs5iLEHWzyL9jLf6acy1)](https://travis.innovate.ibm.com/alchemy-containers/openradiant)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

OpenRadiant is an open source containers management and orchestration platform based on best of breed
containers management technologies such as Kubernetes, Swarm and Mesos. OpenRadiant provides automation
to configure Kubernetes and Swarm for a secure and highly available deployment, provides integration
with multiple network SDN solutions and cloud providers, and configurations to support multitenancy.

* [Quick Start](#quick-start)
* [Other Configurations](#custom-configurations)
* [Learn concepts and commands](#learn-concepts-and-commands)
* [License](#license)
* [Issues](#issues)

### Quick Start

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

Set up your ansible inventory to use the sample project:

```
export ANSIBLE_INVENTORY=examples/envs/dev-vbox/radiant01.hosts
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
ansible-playbook site.yml -e cluster_name=dev-vbox-radiant01 -e envs=examples/envs
```
For more information see ... (TBD)

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
