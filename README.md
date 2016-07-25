## OpenRadiant

[![Travis Lint Status](https://travis.innovate.ibm.com/alchemy-containers/openradiant.svg?token=hs5iLEHWzyL9jLf6acy1)](https://travis.innovate.ibm.com/alchemy-containers/openradiant)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

OpenRadiant is a modular platform for enterprise container-native
devops.  The OpenRadiant platform can be subsetted and/or extended to
create the solution you desire.

Features of the OpenRadiant platform include:
* Kubernetes
* Swarm (original, not experimental)
* Mesos
* Multi-tenancy - with or without Bring-Your-Own-IPv4
* Multi-sharding
* HA control plane in each shard
* Ansible-based create/update/destroy tooling
* Support for a variety of sources of authentication
* Control plane secured by TLS
* Support for a variety of Software-Defined-Network technologies
* Live container crawling

OpenRadiant is a work in progress.  The above features are not yet
available in all combinations.

* [Architecture Overview](#architecture-overview)
* [Ansible Structure](docs/ansible.md)
* [Code of Conduct](#code-of-conduct)
* [Contributing to the project](#contributing-to-the-project)
* [Maintainers](#maintainers)
* [Communication](#communication)
* [Tiny Exaple Solution](examples/tiny-example.md)
* [Proxy documentation](proxy/README.md)
* [Learn concepts and commands](#learn-concepts-and-commands)
* [License](#license)
* [Issues](#issues)

### Architecture Overview

OpenRadiant is software that you can subset and/or extend to produce
software that you use to operate an enterprise container-native
devops service.

One operating instantiation of the full platform is called an
environment, and it contains one or more shards that operate
independently of each other.  Each shard provides Kubernetes, Swarm, and/or
Mesos service.  There is an outer control plane with a proxy API
server that implements the Kubernetes and Swarm APIs --- with
appropriate restrictions and extensions --- by appropriately
transforming each RPC and dispatching it to the appropriate shard.

You can subset OpenRadiant so that it creates just one shard and there
is no API proxy.

In a shard there are worker nodes and control plane nodes.  The
Kubernetes and Swarm workload is dispatched to the worker nodes.  The
control plane nodes run the Kubernetes, Swarm, and/or Mesos control
planes in an HA configuration.  We use Mesos to coordinate resource
allocation decisions between Kubernetes and Swarm.


### Code of Conduct
Participation in the OpenRadiant community is governed by the OpenRadiant [Code of Conduct](CONDUCT.md)

### Contributing to the project
We welcome contributions to the OpenRadiant Project in many forms. There's always plenty to do! Full details of how to contribute to this project are documented in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

### Maintainers
The project's [maintainers](MAINTAINERS.txt): are responsible for reviewing and merging all pull requests and they guide the over-all technical direction of the project.

### Communication
We use \[TODO] [OpenRadiant Slack](https://OpenRadiant.slack.org/) for communication between developers.

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
