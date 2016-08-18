###Overview

This is a directory defined for the tiny example, it includes the configuration files that are specific to a given shard. In the tiny example, the shard is named `dev-vbox-radiant01`):

#####defaults.yml
Specifies the ssh information that ansible uses for the deployment within the machines of the shard. (In the case of tiny examples, thoses machines are the vagrant VMs).

#####radiant01.hosts
Includes the groups and the IP addresses of the machines.

#####radiant01.yml
Includes the virtual IP (VIP) configuration that will be implemented by Ucarp. The VIP will be used as a front end address for the load balancer (HA proxy) binds to.

#####vagrant_config.rb
This is a configuration file used by our tiny example to dimension the vagrant VMs. Through this file, a user can define how many VMs vagrant should instantiate, and the amount of memory allocated for each VM.
