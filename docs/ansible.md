# Ansible

* [Overview](#overview)
* [Ansible versions and bugs and configuration](#ansible-versions-and-bugs-and-configuration)
* [The inventory contract](#the-inventory-contract)
* [The playbooks](#the-playbooks)
* [Ansible variables](#ansible-variables)
* [Networking plugins](#networking-plugins)
* [Temporary and not-so-temporary file locations on the installer](#temporary-and-not-so-temporary-file-locations-on-the-installer)
* [Considerations for deploying on machines that can not pull software from public sources](#considerations-for-deploying-on-machines-that-can-not-pull-software-from-public-sources)


## Overview

The Ansible playbooks in OpenRadiant are modular in three dimensions.
One is componentry: for each major deployed component there is one or
a few roles, and an Ansible variable that controls whether that
component is deployed.  As part of this dimension, OpenRadiant is
extensible regarding the networking technology used.  OpenRadiant has
the concept of a "networking plugin", which supplies the Ansible roles
that deploy the chosen networking technology.  OpenRadiant includes
two networking plugins: (a) one that uses Docker bridge networking
(which is not really functional if you have more than one worker) and
(b) one that uses Flannel with its "host-gw" backend.  As a solution
developer you can supply and configure OpenRadiant to use an alternate
networking plugin that you supply.

The second dimension is the inventory: the playbooks do *not* assume
that your inventory is entirely dedicated to one shard or environment
but rather allow your inventory to include multiple environments and
also things unrelated to OpenRadiant.  This is part of a larger
modularity statement, which is that OpenRadiant aims to constrain the
operator's use of Ansible as little as possible.  In addition to
recognizing that other stuff may be in the Ansible inventory, this
also includes placing minimal requirements on the version of Ansible
used and various other Ansible settings (e.g., config file settings,
module library path).

The third dimension is a factoring between provisioning of machines
and deployment of software on those machines.  There are independent
playbooks for each of those two phases, with a written contract
between them.  Any provisioning technique that results in inventory
content meeting the contract can be used with any of the software
deployment playbooks (because they assume only the contract).


## Ansible versions and bugs and configuration

OpenRadiant aims to place few constraints on the version of Ansible
used and the Ansible configuration.  It has been shown to work with
(a) versions 1.9.5, 1.9.6 (the latest, at the time of this writing,
lower than Ansible 2), and 2.1.1.0 (the latest, at the time of this
writing, in the Ansible 2 line) and (b) no `ansible.cfg`.  Earlier
versions of Ansible 2 have a bug
(https://github.com/ansible/ansible/issues/15930), which can be worked
around by setting the `roles_path` to include the OpenRadiant roles.
Extending OpenRadiant with a network plugin that is not part of
OpenRadiant also requires setting `roles_path` to include the roles
directory(s) containing the network plugin.

Some early versions of Ansible 2 have another bug, for which we have
no work-around.  We suspect it is the 2.0.* versions.  The bug is in
the `docker` module.

The following table summarizes what we know about Ansible versions and
bugs and configuration.

| Version | Usable? | Needs `roles_path` workaround? |
|--------------|-----|-------------------------------|
| 1.9.5, 1.9.6 | Yes | No  |
| 2.0.*        | No  (due to bug in `docker` module) | |
| 2.1.0.*      | Yes | Yes |
| >= 2.1.1.0   | Yes | No  |

OpenRadiant includes an `ansible.cfg` that sets `roles_path` to the
OpenRadiant roles directory and has some other settings whose purpose
is speeding up the Ansible execution.  A platform operator is free to
use a different `ansible.cfg`, provided `roles_path` is set if and as
necessary.

If you use OpenRadiant with a provisioning technology that involves a
dynamic Ansible inventory then you will need a config file and/or
environment variable settings as appropriate for that dynamic
inventory.


## The inventory contract

Installing OpenRadiant is factored into two orthogonal phases:
provisioning machines and deploying software on them.  Any acceptable
technique for provisioning works with all the software deployment
playbooks.  What enables this orthogonality is a contract between the
two phases.  The contract concerns "groups" in the Ansible inventory.
OpenRadiant does not care how your inventory is
produced/stored/written.  You could keep one big `/etc/ansible/hosts`
file.  You could keep many files in an `/etc/ansible/hosts/`
directory.  You could keep various inventory files in various places
and pass the relevant one via the `-i` argument or the
`ANSIBLE_INVENTORY` environment variable.  You could use an Ansible
dynamic inventory.  All that matters is what groups exist and what
machines are in them.

Because the playbooks support multiple environments and shards, the
inventory contract applies to a given environment or shard.  A shard
is also known as a cluster.


### The structure of environment and shard names

To support the structure of environments and shards, an environment's
name must have the form
`{{environment_role}}-{{environment_location}}` and a shard (AKA cluster) name
must have the form `{{environment_name}}-{{cluster_short_name}}`
(i.e.,
`{{environment_role}}-{{environment_location}}-{{cluster_short_name}}`).
There may be no dash in the `environment_kind` nor the
`environment_location`.  OpenRadiant presently leaves it up to you to
decide how to choose the environment role and location names.  In the
tiny example you find a cluster named `dev-vbox-radiant01`.

Because a shard (cluster) name contributes to host names, to file
names, to Ansible group names, and to various things in various
compute provider clouds (all of which impose restrictions on the
characters used), a shard name should start with a lowercase letter
and only contain lower case letters, numbers, and the hyphen ('-').


### The inventory contract for a shard

The contract for a shard specifies that certain inventory groups
exist.  All these groups are required, even if their purpose is
related to a component that is not being deployed.  The groups are as
follows.

* `cluster-{{cluster_name}}`: all the machines in the cluster ---
  i.e., the union of the following groups

* `workers-{{cluster_name}}`: the worker nodes in the cluster

* `etcd-{{cluster_name}}`: where to run etcd servers; should be an odd
  number of machines

* `zk-{{cluster_name}}`: where to run ZooKeeper servers; should be an
  odd number of machines

* `lb-{{cluster_name}}`: where to run the load balancers for the master
  components

* `k8s_auth-{{cluster_name}}`: where to run OpenRadiant's Kubernetes
  authorization servers

* `k8s_master-{{cluster_name}}`: where to run the Kubernetes master
  components

* `swarm_master-{{cluster_name}}`: where to run the Swarm master components

* `mesos_master-{{cluster_name}}`: where to run the Mesos master components

By defining a group per HA component this contract enables a variety
of deployment patterns: a set of machines per component, all
components on each of a set of machines, and many in between.

Except as noted above, you can put any number of machines in each of
those groups.

There are also some requirements on the provisioned machines.

* They must run Ubuntu 14.04.

* The Ansible inventory hostname of each machine must be an IPv4
  address.

* They can send IP packets to each other using the IP address by which
  they are known in the Ansible inventory.

* You (the user running the Ansible playbook) must be able to `ssh`
  from the Ansible controller machine to the configured user on the
  provisioned machines, and that user must be able to `sudo` to `root`
  without supplying a password.  The configured user is the remote
  user determined by the usual Ansible rules --- including the
  settings (if any) of the `ansible_ssh_user`, `ansible_user`, and
  `remote_user` Ansible variable(s) in any of the places where Ansible
  allows them to be set.


### The inventory contract for an environment

The contract for an environment includes, beyond the shards of that
environment, the following inventory group.

* `proxy-{{env_name}}`: where to run the API proxy for that environment

Currently that inventory group must contain exactly one machine.  In
the future we will support an HA cohort.

The requirements on machines on that group are as follows.

* They must run Ubuntu 14.04.

* You (the user running the Ansible playbook) must be able to `ssh`
  from the Ansible controller machine to the configured user on the
  provisioned machines, and that user must be able to `sudo` to `root`
  without supplying a password.  The configured user is the remote
  user determined by the usual Ansible rules --- including the
  settings (if any) of the `ansible_ssh_user`, `ansible_user`, and
  `remote_user` Ansible variable(s) in any of the places where Ansible
  allows them to be set.

* The proxy machines must be able to open TCP connections into the
  shards using the Ansible inventory hostnames (which are IPv4
  addresses) for the shard machines.


## The playbooks

OpenRadiant currently includes no playbooks for provisioning machines,
and two playbooks for installing software one them.  These are
`ansible/env-basics.yml` and `ansible/shard.yml`.  We are creating
additional playbooks for use at the environment level.


## Ansible variables

The playbooks that deploy software on machines are parameterized by a
large collection of Ansible variables.  OpenRadiant introduces a
distinction between _primary_ variables, which a platform operator may
set, and _secondary_ variables, which are really just intermediates in
the logic of the Ansible roles.

### Primary Ansible variables for the shard playbook

Six have no meaningful defaults.  The others are defined by defaults
in `ansible/group_vars/all` and can be overridden by settings in an
environment-specific file and a cluster-specific file.

#### Primary shard variables that have no default

Following are the six with no meaningful default.

* `cluster_name`: this identifies the cluster being processed, as
  discussed above.

* `envs`: an absolute or relative pathname for the directory under
  which the environment- and shard-specific files are found (see
  [below](#additional-files-for-setting-ansible-variable-values)).  If
  relative then it is interpreted relative to the directory containing
  the shard playbook.  The settings for the shard named `A-B-C` are
  found in `{{envs}}/A-B/C.yml`.

* `network_kind`: this is the name of the network plugin to use.  This
  must be supplied on the Ansible command line (because an Ansible
  playbook can not invoke a role whose name is not computable at the
  start of the playbook); in the future we will provide a shell script
  that reads this variable's setting according to our conventions and
  then invokes the playbook.

* `master_vip`: this identifies the virtual IP (VIP) address for most
  of the master components.

* `master_vip_net` and `master_vip_net_prefix_len`: these two
  variables identify the subnet
  (`{{master_vip_net}}/{{master_vip_net_prefix_len}}` in CIDR
  notation) that contains the `master_vip`.  That should be either (1)
  identical to the subnet of the network interface (on the master
  machines) identified by the `network_inteface` variable or (2) a
  subnet, disjoint from all the others in the system, that all the
  machines in the cluster know to be on the same network as (1).
  These two variables are needed only if the master components are
  being deployed in an HA configuration.

The settings for the master components VIP can be in the environment-
and/or cluster-specific variables files (see next section).  The
`cluster_name`, `envs`, and `network_kind` must be supplied on the
command line invoking the playbook.


#### Additional files for setting Ansible variable values

The OpenRadiant shard playbook follows a convention for reading
variable values from additional files.  When working on a cluster
named `{{env_name}}-{{cluster_short_name}}`, the following files are
relevant to the settings of the Ansible variables.

* `ansible/group_vars/all` (and all the rest of the usual Ansible
  story, which is mostly not exercised).  This provides the baseline
  defaults.  This is maintained by the OpenRadiant developers and is
  not expected to be modified by the platform operator.

* `{{envs}}/{{env_name}}/defaults.yml`.  The playbook explicitly loads
  this, and its settings override those from `ansible/group_vars/all`.
  The platform operator is expected to provide this file.

* `{{envs}}/{{env_name}}/{{cluster_short_name}}.yml`.  The playbook
  explicitly loads this, after the previous file --- so its settings
  override the others.  The platform operator is expected to provide
  this file.


#### Primary shard variables that have defaults

Following are the most important of those variables.

* `kubernetes_deploy`, `swarm_deploy`, `mesos_deploy`: these control
  which of these three components are deployed.  The usual Ansible
  conventions for values for booleans apply.

* `k8s_hyperkube_image`, `k8s_hyperkube_version`, `kube_image`: these
  configure the kubernetes image to use; the first two are put
  together as image name and tag when Mesos is *not* involved; the
  third contains both image name and tag and is used when Mesos *is*
  involved.

* `mesos_master_image`, `mesos_slave_image`: image name&tag to use for
  Mesos on master and worker nodes (respectively).

* `swarm_image`: image name and tag to use for the Swarm manager.

* `ha_deploy`: controls whether a load balancer is deployed in front
  of the master components.  You pretty much want this when the master
  components are deployed in a highly available manner; defaults to
  true.

* `network_interface`: the network interface to use on the target
  machines, for all control and data plane activity.  Yes, currently
  there is support only for using a single network interface for all
  purposes.

See `ansible/group_vars/all` for the definitions and defaults of all
the variables that have defaults.


### Secondary Ansible variables for the shard playbook

These are of no interest to a deployer of OpenRadiant or a solution
built on it.  These _are_ of interest to a developer of OpenRadiant or
a solution that extends OpenRadiant.  Mostly the secondary variables
are just set in certain roles and read in others.

Following are the secondary variables of interest to a developer.

* `etcd_deploy`: A boolean that controls whether etcd is deployed.
  Initially set to `False`.  After the networking plugin's
  `-variables` role has had a chance to set this variable to `True`,
  this variable is also set to `True` if any other settings imply the
  need for etcd.

* `k8s_worker_kubelet_network_args`: An array of strings.  If
  Kubernetes is being used without Mesos then these are arguments to
  add to the command line that runs the Kubernetes kubelet on worker
  nodes; if Mesos is involved then these instead need to be the
  arguments added to the k8s scheduler command line to influence the
  way the kubelets are run on worker nodes.  Initially set to the
  empty array; the neworking plugin may set it otherwise.  Moot if
  Kubernetes is not being deployed.

* `use_kube_system_kubedns`: A boolean that controls whether the KubeDNS
  application is deployed as usual in the `kube-system` namespace.

* `kube_proxy_deploy`: A boolean that controls whether the
  `kube-proxy` is used.  Defaults to `True`; consulted only after the
  networking plugin has a chance to set it to `False`.

* `is_systemd` and `is_upstart`: these identify the service manager
  used on the target machine; not defined for localhost.  Currently
  these are the only two recognized.  The systemd case is recognized
  but not yet fully supported.


### Primary Ansible variables for the enviornment playbooks

* `env_name`: a string of the sort discussed above.  Must be supplied
  by the platform operator in the command that invokes the environment
  playbook.

* `proxy_deploy`: controls whether the API proxy is deployed in the
  environment.  *NB*: this is a hypothetical, not implemented right
  now.

The environment playbook(s) follow a convention for reading variable
values from an additional file --- the environment-specific file
discussed
[above](#additional-files-for-setting-ansible-variable-values).


## Networking plugins

The job of a networking plugin is to deploy the networking technology
used for the workload.  This starts with the way containers are
networked and includes, if Kubernetes is in the picture, the DNS
service expected in a Kubernetes cluster and the implicit load
balancers for non-headless Kubernetes services.

As noted above, OpenRadiant includes the following networking plugins.

* `bridge`: this uses Docker bridge networking and thus is not really
  functional when there are multiple worker nodes in a shard.
  Provides non-multi-tenant DNS to users via the kube DNS application.
  Uses kube-proxy to provide implicit load balancers.

* `flannel`: this uses Flannel networking with its `host-gw` backend.
  This supports multiple worker nodes and uses ordinary IP routing and
  connectivity.  It does not support the Kubernetes API for network
  filtering nor any implicit network filtering for containers created
  through the Docker API.  Provides non-multi-tenant DNS to users via
  the kube DNS application.  Uses kube-proxy to provide implicit load
  balancers.

To create a networking plugin, the developer needs to define three
Ansible roles.  A networking plugin named `fred` supplies the
following Ansible roles.

* `fred-variables`: This role is invoked on localhost and on all
  cluster machines after the global default settings and environment-
  and cluster-specific settings, before anything else is done.

* `fred-local-prep`: This role is invoked on localhost after all the
  variables have been set and before any work is done on the managed
  machines.

* `fred-remote-prep`: This role is invoked on all the cluster
  machines, after the local prep and after the deployment of Docker
  and etcd and/or ZooKeeper, before the deployment of higher level
  frameworks such as Mesos or Kubernetes.


## Temporary and not-so-temporary file locations on the installer

The playbooks use (a) a temporary directory on the installer machine
for files needed for the duration of one run of a playbook and (b) a
non-temporary directory on the installer machine for files that should
persist throughout all the playbook runs regarding a given
environment.

The temporary files are kept in `{{ lookup('env','HOME') }}/tmp/`.
The platform operator is allowed to retain them between playbook runs, and is
allowed to delete them between playbook runs.  To allow the operator
to work on multiple deployments concurrently, it is recommended that
the temporary files be kept in a subdirectory that is specific to the
environment or shard being deployed.

The non-temporary files are kept in `{{ lookup('env','HOME')
}}/.openradiant/envs/{{ env_name }}/` or, in the case of files
specific to a shard, `{{ lookup('env','HOME') }}/.openradiant/envs/{{
env_name }}/{{cluster_tail}}/` (where `cluster_tail` is the shard
short name, the part after the environment name).

For example, there is a CA and admin user that are shared throughout
an environment.  Their certs and keys are kept in `{{
lookup('env','HOME') }}/.openradiant/envs/{{ env_name
}}/admin-certs/`.  There are additional certs and keys that are
specific to each shard, and they are kept in `{{ lookup('env','HOME')
}}/.openradiant/envs/{{ env_name }}/{{cluster_tail}}/admin-certs/`
(even though none of them is for "admin").

The platform operator is responsible for keeping the contents of
`~/.openradiant/envs/{{ env_name }}/` persistent throughout the
lifetime of the named environment and sharing it among all the
installer machine(s) for that environment.


## Considerations for deploying on machines that can not pull software from public sources

There are two sorts of sources considered here: apt repositories and
Docker registries (is it really true that there is no solution for
pypi here?).  In both cases the platform operator must create a mirror
from which the target machines _can_ pull.  The operator also needs to
configures OpenRadiant to use those mirrors.

For apt the answer is to configure an Ansible variable named `repos`
that contains a list of records, one for each alternate apt repo to
use.  Each record must contain a field that is named `repo` and has a
value like `deb http://10.20.30.40:8082/ trusty main`.  A record may
also have a field named `validate_certs`.  These two fields are used
as the values of the corresponding arguments of the `apt_repository`
module in Ansible.  A record may also have a field that is named
`key_url` and has a value like `http://10.20.30.40:8082/Release.pgp`;
in this case the `key_url` and `validate_certs` (if any) fields are
used as the `url` and `validate_certs` (respectively) arguments of the
`apt_key` module in Ansible.  Additionally, a record may also have a
field that is named `key_package`; in this case Ansible's `apt` module
is invoked with the `key_package` value passed through the `pkg`
argument.

For Docker the answer is to redefine the relevant Ansible variables
that indicate the name and tag of the image to pull.  These variables
appear in `ansible/group_vars/all`.
