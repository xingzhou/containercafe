## Ansible

The Ansible playbooks in OpenRadiant are modular in three dimensions.
One is componentry: there is an includable Ansible playbook for each
major component.  The second dimension is the inventory: the playbooks
do *not* assume that your inventory is entirely dedicated to one shard
or environment but rather allow your inventory to include multiple
environments and also things unrelated to OpenRadiant.

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

* `cluster-{{cluster_name}}`: all the machines in the cluster

* `workers-{{cluster_name}}`: the worker nodes in the cluster

* `etcd-{{cluster_name}}`: where to run etcd servers

* `zk-{{cluster_name}}`: where to run ZooKeeper servers

* `lb-{{cluster_name}}`: where to run the load balancer(s) for the master components

* `k8sm_auth-{{cluster_name}}`: where to run OpenRadiant's RemoteABAC server(s)

* `k8sm_master-{{cluster_name}}`: where to run the Kubernetes master components

* `swarm_master-{{cluster_name}}`: where to run the Swarm master components

* `mesos_master-{{cluster_name}}`: where to run the Mesos master components

By defining a group per HA component this contract enables a variety
of deployment patterns: a set of machines per component, all
components on each of a set of machines, and many in between.

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
