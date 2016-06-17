## OpenRadiant

[![Travis Lint Status](https://travis.innovate.ibm.com/alchemy-containers/radiant-ansible.svg?token=hs5iLEHWzyL9jLf6acy1)](https://travis.innovate.ibm.com/alchemy-containers/openradiant)
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

export ANSIBLE_INVENTORY=
export ANSIBLE_LIBRARY=$ANSIBLE_INVENTORY

create cluster definition
```
ansible-playbook -e cluster_name=dev-mon01-radiant01-pd -e envs=../radiant-envs/envs site.yml
```
For more information see ...

### Custom Configurations

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
