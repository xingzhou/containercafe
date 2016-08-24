# Configuration for the masters and the worker VMs
# The roles (masters/workers) will be assigined in the ansible inventory
$vm_count=2
$vm_memory=1536
$vm_cpu=1

#Configuration for the proxy_VM hosting the proxy
$proxy_num=1
$proxy_vm_memory=1024

#Configuration for the installer_VM hosting the ansible controller
$installer_num=1
$installer_vm_memory=1024
