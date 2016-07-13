## Vagrant Quick Deployment

Vagrant is a tool for building complete virutalized development environments. 
Vagrant uses a base image to quickly clone a virtual machine. These base images are known as "boxes" in Vagrant
A "VagrantFile" is used to describe how to configure and provision our machines.
To provision a development environment with vagrant we use the "vagrant up" command in the same directory branch of the VagrantFile.

$ vagrant up

After the completion of the above command, we should have a running vagrant environment. We can verify this environment by using the "vagrant status" command.

$ vagrant status

A sample output would look like:

Current machine states:

radiant2                  running (virtualbox)
radiant3                  running (virtualbox)

There are multiple ways to ssh into a vagrant machine, the easiest is to use "vagrant ssh".

$ vagrant ssh {machine_name}  e.g. vagrant ssh radiant2

Alternatively we can use the typical way to ssh into a machine:

$ ssh user@machine_address e.g. ssh vagrant@192.168.10.2

We can use the vagrant ssh to externally execute commands inside the machine. For instance to get the ip address of the machine we can use:

$ vagrant ssh radiant2 -c 'ifconfig'

### Troubleshooting the deployment

##### If you try to run any vagrant command and you get the following error:

A Vagrant environment or target machine is required to run this
command. Run `vagrant init` to create a new Vagrant environment. Or,
get an ID of a target machine from `vagrant global-status` to run
this command on. A final option is to change to a directory with a
Vagrantfile and to try again.

This could either be:

1- you are NOT under the same user that created the vagrant environment, e.g. if you are not root, try sudo su
2- you are in a different directory than the one where the vagrant environment was created, change directory to your /path/to/vagrant

##### If you try to run the "vagrant up" command and you get an error related to the uid:

Most likely you do not have enough permissions, if you are not root, try sudo su, and then run "vagrant up"

##### If you try to run the "vagrant up" command and you get an error related to exisiting machines:

A VirtualBox machine with the name 'radiant3' already exists

Then you have to clean up this machine from your environment. there are two way to do that:

1- the soft way:

    $ vagrant destroy radiant3
    
    or simply "vagrant destroy", which will destroy all the machines
    
2- the hard  way:
    The hard way should be used if the soft way does not resolve the problem. In this case we clean up the VMs using virtual box.
    To clean up a VM with vitual box, we have to 1) power off 2) detach the hard disk 3) delete the VM. It is very important to remove the hard disk "--hdd none" as shown below:
    
    $ vboxmanage list runningvms
    $ vboxmanage controlvm radiant3 poweroff
    $ vboxmanage modifyvm radiant3 --hdd none
    $ vboxmanage unregistervm radiant3  -delete

After cleaning up the VM, we can retry the "vagrant up" command.
