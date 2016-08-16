##Overview

OpenRadiant can be deployed in different environments. the tiny example is a vagrant based environment that can get you started on your local machines with just a couple of commands.
If you wish to use OpenRadiant in a different environment, there are a few changes that you need to make in order to customize the deployment to you environment.


The following is a quick list of things to consider when deploying OpenRandiant in your environment.

####Inventory

The inventory file is a way to instruct ansible on where to deploy the various OpenRadiant components. In our tiny example this file is located at: openradiant/examples/envs/dev-vbox/radiant01.hosts
Since OpenRadiant allows for a modular where you can pick and chose which/where the components can be deployed, the inventory file have several groups where some of them default to the value of others.
This file should be modified by adding the specific IP addresses for the machines you wish OpenRadiant to be deployed on.

You can also use your own inventory file (as long as you preserve the [group_names]), simply let  ANSIBLE_INVENTORY point to the path of your inventory file instead.
E.g. in the tiny example, use: export ANSIBLE_INVENTORY=../path/to/your/inventory
For detailed information about the ansible variables and the inventory, check [our Ansible doc](docs/ansible.md).

####User SSH info

Ansible needs the ability to ssh into the machines where OpenRadiant will be deployed. You have to make sure that the ssh keys, file permissions, etc. are properly configured on your machines.
You also need to make sure that the ansible configuration in openradiant/examples/envs/dev-vbox/defaults.yml corresponds to your environment (your ssh user and ssh key location)

####Environment specific variables

Our ansible uses the default values for certain environment specific variables. You might need to modify those default values in our openradiant/ansible/group_vars/all file acccording to your environment.
For example, if your control plane IP address is configured on Eth0 instead of Eth1, you would need to modify "network_interface: eth1" accordingly.
