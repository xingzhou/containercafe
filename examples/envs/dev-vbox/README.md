###Overview

This is a directory defined for the tiny example, it includes the configuration files that are specific to a given shard. 
You can define your own configuration by making a copy of this folder and setting the value of the environment variable `ENV_NAME`
to the name of the folder before running `vagrant up`. For the tiny example `ENV_NAME` is `dev-vbox`.
In the tiny example, the shard is named `dev-vbox-radiant01`:

#####defaults.yml
Specifies the variable settings that are particular to the tiny example’s environment. In our example these settings include the ssh information that ansible uses for the deployment within the machines (or virtual machines) of the shard. 

#####radiant01.hosts
Includes the groups and the IP addresses of the machines.

#####radiant01.yml
Specifies the environment variables specific to a shard. In our example it Includes the virtual IP (VIP) configuration that will be implemented by Ucarp. The VIP will be used as a front end address for the load balancer (HA proxy) binds to.

#####vagrant_config.rb
This is a configuration file used by our tiny example to dimension the vagrant VMs. This file defines how many VMs vagrant should instantiate, and the amount of memory allocated for each VM.
