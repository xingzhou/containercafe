##Overview

OpenRadiant can be deployed in different ways.  The [tiny-example](tiny-example.md) uses Vagrant and Virtual Box to get you started on your local machines with just a couple of commands.
If you wish to use OpenRadiant for a larger deployment, there are a few changes that you need to make in order to configure your deployment.


The following is a quick list of things to consider when configuring to deploy OpenRandiant.

####Inventory

OpenRadiant is installed by Ansible playbooks, which determine where
to install OpenRadiant based on certain Ansible inventory groups whose
names are derived from the names of the OpenRadiant environment and
shard being deployed.  See
[the inventory contract](../docs/ansible.md#the-inventory-contract)
for full details.

In our tiny example the Ansible inventory is in a file located at
[openradiant/examples/envs/dev-vbox/radiant01.hosts](envs/dev-vbox/radiant01.hosts).
In that file you will see all the required inventory groups (plus one
that is not required).  In this example many of the groups have the
same membership, but you can have a different layout in your other
deployments.  Identify machines by their IP addresses.  These
addresses should be reachable from each other and also from the
Ansible controller machine.

As long as your inventory uses the prescribed group names, you can use
any of Ansible's ways of supplying inventory (for examples: one file,
a directory containing several files, a program providing dynamic
inventory) and ways of identifying the desired source (for examples:
the `ANSIBLE_INVENTORY` environment variable, the `-i` argument, the
`/etc/ansible/hosts` location).

For example:
```bash
export ANSIBLE_INVENTORY=path/to/your/inventory_as_one_file
```


####User SSH info

Ansible needs the ability to ssh into the machines where OpenRadiant
will be deployed. You have to make sure that the ssh keys, file
permissions, etc. are properly configured on your machines.  You also
need to tell the OpenRadiant playbooks the remote user account to
login to and the private key file (on the Ansible controller machine)
to use to authenticate.  The usual Ansible rules and options apply
here.  You can use the `ansible_ssh_user` or `remote_user` variable
(or `ansible_user` if you are using Ansible 2), or set none of those
variables and take the default (which is that the remote username is
the same as your local username).  Similarly, you can let Ansible use
your default private key or identify the desired one in the
`ansible_ssh_private_key_file` variable.

In the tiny example, that configuration is done in the environment
variables file at
[openradiant/examples/envs/dev-vbox/defaults.yml](envs/dev-vbox/defaults.yml).

If you are going to be repeatedly creating and destroying the machines
in your environment (and who isn't?), you may also want to configure
SSH to not check host SSH keys.  OpenRadiant includes an `ansible.cfg`
file that does this.


####Environment specific variables

OpenRadiant defines a bunch of Ansible variables and their default
values in
[openradiant/ansible/group_vars/all](../ansible/group_vars/all).  The
tiny example accepts all the default values, but you might want to
override some of them.  Put your overrides in your
environment-specific and/or shard-specific variables files.  For
example, if `eth0` is the right network interface for OpenRadiant to
use on all the machines in your environment then you would set the
following in your environment-specific variables file.

```
network_interface: eth0
```
