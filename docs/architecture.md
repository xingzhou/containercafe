# Architecture

OpenRadiant is software that you can use, or extend to produce more
capable software that you use, to operate an enterprise
container-native devops service.

One operating instantiation of the full platform is called an
environment, and it contains one or more shards that operate
independently of each other.  See
[a picture of a deployed environment](media/DeployedTopology.svg).
The multiplicity of shards is intended primarily for operational
flexibility.  It allows an environment to evolve by addition and
removal of shards when modifying an existing shard would be more
difficult.  It also can work around scalability limits of one shard,
to a limited degree; OpenRadiant is not engineered for a large number
of shards.  OpenRadiant includes an API proxy that proxies each
request to the appropriate shard.

Each shard provides Kubernetes, eventually Swarm (V1, not the "swarm
mode" introduced in Docker 1.12), and any platforms added by
extensions.  The operator adds extensions by deploying them in a shard
after the OpenRadiant installer does its work.  The API proxy
currently proxies only Kubernetes and SwarmV1 requests.

The API proxy has two main jobs: multi-sharding and multi-tenancy.
The API proxy implements the Kubernetes and Swarm APIs --- with
restrictions and extensions needed for multi-tenancy --- by
appropriately transforming each RPC and proxying it to the appropriate
shard.  Each tenant is confined to one shard.  The API proxy includes
a plugin framework for authentication of users and tenants, and
keeping track of which tenant is in which shard.  When the API proxy
proxies a request to the Kubernetes API server or SwarmV1 manager, it
does so impersonating the original client.  The authentication plugin
keeps track of the credentials needed to do so.

OpenRadiant currently includes one authentication plugin, which keeps
its data in files.  This is a very early point-in-time statement.

You can operate OpenRadiant with just one shard.

In a shard there are worker nodes and master nodes.  The Kubernetes
(and eventually Swarm workload) is run on the worker nodes.  The
master nodes run the Kubernetes, eventually Swarm, and any extension's
master components in an HA configuration.  We are working on a better
solution than Mesos to coordinate resource allocation decisions
between Kubernetes and Swarm.  See
[a picture of a shard](media/DeployedShard.svg) for a picture of a
shard with Kubernetes, SwarmV1, and Mesos.

The layout of control plane components onto physical machines is not
shown in that picture, because there are many options for how that is
done.  The operator indicates her choice by appropriately populating
the relevant inventory groups in her Ansible inventory for the shard;
see
[the inventory contract for a shard](ansible.md#the-inventory-contract-for-a-shard)
for details.

Not clearly visible in that picture is the networking technology used
for the workload (although the kube-proxy --- which is part of that in
some cases --- _is_ shown).  As noted
[elsewhere](ansible.md#networking-plugins), this is the purview of the
_networking plugin_ chosen by the service operator.

OpenRadiant includes Ansible-based installation technology to
instantiate an OpenRadiant environment.  An installer machine acts as
Ansible controller to install OpenRadiant in a target environment.
The installation is parameterized by some Ansible variables files that
describe the desired target environment.

OpenRadiant deploys Kubernetes in containers.  You can choose whatever
version you want.  Your configure your choice of image (including
tag).  See the doc about
[the relevant configuration variables](docs/ansible.md#primary-shard-variables-that-have-defaults)
and
[where to put your settings for those variables](docs/ansible.md#additional-files-for-setting-ansible-variable-values).

The Ansible playbooks strive to meet the Ansible ideal of achieving a
prescribed desired state, and can thus be used to update as well as
install.  However, there are limits to the space of initial states
with which these playbooks can cope.
