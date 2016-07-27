# Virtual IP

This role will add virtual IP support to the nodes on which in runs.  Virtual IP is acheived by using UCARP.
Currently this role is applied to the controller (API) hosts.

## How it works

When this role is applied, it installs the UCARP package and copies 3 files to the machine:

* /etc/network/if-up.d/vip - This file starts the UCARP process on the specified interface when the interface
started.
* /etc/network/vip-up.sh - This script is called from the UCARP process when it becomes the master.  Currently this
script just adds the virtual ip address to the interface.  This script could do other things, but that's all it's doing
at this point.
* /etc/network/vip-down.sh - This script is called from UCARP when it loses the master role.  Currently this script
just removes the virtual ip address from the interface.

Once the files are copied to the system, the playbook will start the UCARP process on each of the hosts and wait until
one of them becomes the master and the virtual ip address is available.

When the systems are started, UCARP listens for an advertisement from the master, if it doesn't receive one, then it
becomes the master and runs the vip-up script which adds the virtual ip to the interface.  As the other systems
come online, their UCARP processes will be in backup mode.  If the master host goes down, then one of the other host
systems will become the master.  There are lots of settings in UCARP to adjust how and when this happens.  So far
the default settings seem to work fine.

Note: There is a setting to have one of the systems be the primary master, the problem with this is that if that system
goes down and comes back online again, it will become the master, this results is a brief connection failures while
the switch occurs. It doesn't seem like we want that behavior, so currently the master setting is no on all instances.
Therefore, once a system becomes the master it stays the master until it goes down or someone manually forces it to
lose the master status.

## How to enable
In your group_vars/all.yml file, update the following attributes, in this example, we are using 10.140.23.246 in subnet
10.140.23.240/28.

    master_vip: 10.140.23.246  << Your floating ip address
    undercloud_cidr:
     - 10.140.179.0/24        
     - 10.140.23.240/28        << Add your floating ip subnet (if different)

    undercloud_master_vip: '{{ master_vip }}'  << Set the undercloud to your floating ip address

    vip:
      enabled: true          << Set the vip enabled flag to true
      interface: eth0        << Set this interface if different, defaults to eth0
      cidr_prefix: '/32'     << Update the virtual ip cidr prefix for the subnet.  Defaults to /32
      ucarp_id: '18'         << Update the ucarp id for your environment.  The value is 1-255.  Note that this 
                                value needs to be unique within the same multicast IP group. 
  
Depending on the value of the fqdn attribute, you may need to re-generate your ssl certificates.

## How to tell which host is master
The easiest way is to simply ssh to the virtual ip address.  You can also check the interface for the virtual ip
address:

    ip addr

## How to test

Once all the systems are online, i started a nova list in a loop on one of the non-controller systems.

    while true; do date; nova list; done
    
You'll see the nova list returning the list of server instances. The date call is there simply to show it's running. 

Now power off the master using the slcli command to simulate a failure

    slci vs power_off <host>

Now you will see the connection errors from nova list, in about 5s another API host will become the master, and
nova list will work again.  You can now power on the host and it will be in backup mode.

    slcli vs power_on <host>

You can also force a master to give up control by issuing a SIGUSR2 to the ucarp process:

    ps -ef | grep ucarp
    kill -s SIGUSR2 <ucarp process id>
    
## Debugging

By default UCARP uses multicast IP address of 224.0.0.18.  If the ucarp id is the same as another environment (a
different virtual ip address), then you will get these warnings in the syslog:

    [WARNING] Bad digest - md2=[b02dfc0e...] md=[84b05a8c...] - Check vhid, password and virtual IP address
    
The easiest fix is to use a different ucarp_id value.  You can see which ids are in use on your subnet by doing a tcpdump:

    tcpdump -i eth0 -n 'net 224.0.0.0/8'
    
You'll see packets like this:

    21:13:56.556490 IP 192.168.130.9 > 224.0.0.18: VRRPv2, Advertisement, vrid 130, prio 254, authtype none, intvl 1s, length 28
    21:13:56.597069 IP 10.140.174.10 > 224.0.0.18: VRRPv2, Advertisement, vrid 1, prio 0, authtype none, intvl 1s, length 36
    21:13:56.874891 IP 10.140.179.137 > 224.0.0.18: VRRPv2, Advertisement, vrid 18, prio 0, authtype none, intvl 1s, length 36
    21:13:57.556664 IP 192.168.130.9 > 224.0.0.18: VRRPv2, Advertisement, vrid 130, prio 254, authtype none, intvl 1s, length 28
    21:13:57.596708 IP 10.140.174.10 > 224.0.0.18: VRRPv2, Advertisement, vrid 1, prio 0, authtype none, intvl 1s, length 36
    21:13:57.874878 IP 10.140.179.137 > 224.0.0.18: VRRPv2, Advertisement, vrid 18, prio 0, authtype none, intvl 1s, length 36

The vrid is the ucarp id.  The master node transmits the advertisement.
