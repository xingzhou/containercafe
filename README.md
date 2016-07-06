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

Install [Vagrant](https://www.vagrantup.com/) and
[VirtualBox](https://www.virtualbox.org/wiki/Downloads). You will need
at least Vagrant 1.8.4 and VirtualBox 5.0.24.

Please make sure you have most recent version of [VirtulBox](https://www.virtualbox.org/wiki/Downloads)
5.0.24 or higher is required

Required level of Ansible is 1.9.6:
```
pip install --upgrade ansible==1.9.6
```

Please make sure you have most recent version of [VirtulBox](https://www.virtualbox.org/wiki/Downloads)
5.0.24 or higher is required

Required level of Ansible is 1.9.6:
```
pip install --upgrade ansible==1.9.6
```

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
