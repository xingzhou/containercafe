## Ansible

The Ansible playbooks in OpenRadiant are modular in three dimensions.
One is componentry: for each major component there is one or a few
roles.  That will be changing, so that for each major component there
is an includable playbook.  The second dimension is the inventory: the
playbooks do *not* assume that your inventory is entirely dedicated to
one shard or environment but rather allow your inventory to include
multiple environments and also things unrelated to OpenRadiant.

The third dimension is a factoring between provisioning of machines
and deployment of software on those machines.  There are independent
playbooks for each of those two phases, with a written contract
between them.  Any provisioning technique that results in inventory
content meeting the contract can be used with any of the software
deployment playbooks (because they assume only the contract).

Because the playbooks support multiple environments and shards, the
inventory contract applies to a given environment or shard.  Currently
there are no playbooks for environments, only shards.  A shard is also
known as a cluster.  The contract for a cluster specifies that certain
inventory groups exist.  They are as follows.

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
* You must be able to `ssh -l root` from your controller machine to provisioned machines without supplying a password.

Because a cluster name contributes to host names, to Ansible group
names, and to various things in various compute provider clouds (all
of which impose restrictions on the characters used), a cluster name
should start with a lowercase letter and only contain lower case
letters, numbers, and the hyphen ('-').

To support the structure of environments and shards, an environment's
name must have the form
`{{environment_role}}-{{environment_location}}` and a cluster name
must have the form ``{{environment_name}}-{{cluster_short_name}}`
(i.e.,
`{{environment_role}}-{{environment_location}}-{{cluster_short_name}}`).
There may be no dash in the `environment_kind` nor the
`environment_location`.  OpenRadiant presently leaves it up to you to
decide how to choose the environment role and location names.  In the
tiny example you find a cluster named `dev-vbox-radiant01`.

OpenRadiant currently includes no playbooks for provisioning machines,
and one playbook for installing software one them.  That is
`ansible/site.yml`.

The playbooks that deploy software on machines are parameterized by a
large collection of Ansible variables.  Five have no meaningful
defaults.  The others are defined by defaults in
`ansible/group_vars/all` and can be overridden by settings in an
environment-specific file and a cluster-specific file.

Following are the five with no meaningful default.

* `cluster_name`: this identifies the cluster being processed, as
  discussed above.

* `envs`: this identifies the parent directory under which the
  environment- and cluster-specific files are found.  The settings for
  the environment named `A-B` are found in `{{envs}}/A-B/defaults.yml`
  (if relative, the base is the filename of the playbook).

* `master_vip`: this identifies the virtual IP (VIP) for most of the
  master components.

* `master_vip_net` and `cidr_prefix`: these two variables identify
  the subnet (`{{master_vip_net}}/{{cidr_prefix}}`) that contains the
  `master_vip`.  That should be either (1) identical to the subnet of
  the network interface (on the master machines) identified by the
  `network_inteface` variable or (2) a subnet, disjoint from all the
  others in the system, that all the machines in the cluster know to
  be on the same network as (1).  These two variables are needed only
  if the master components are being deployed in an HA configuration.

When working on a cluster named `{{env_name}}-{{cluster_short_name}}`,
the following files are relevant to the settings of the Ansible
variables.

* `ansible/group_vars/all` (and all the rest of the usual Ansible
  story, which is mostly not exercised).  This provides the baseline
  defaults.

* `{{envs}}/{{env_name}}/defaults.yml`.  The playbook explicitly loads
  this, and its settings override those from `ansible/group_vars/all`.

* `{{envs}}/{{env_name}}/{{cluster_short_name}}.yml`.  The playbook
  explicitly loads this, after the previous file --- so its settings
  override the others.

The settings for the master components VIP should appear in the
cluster-specific file.  The `cluster_name` and `envs` should be
supplied on the command line invoking the playbook.


