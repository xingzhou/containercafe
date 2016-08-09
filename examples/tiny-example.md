## Tiny Example

This produces a very simple demonstration cluster of two VirtualBox
VMs, one master and one worker.  They have Mesos installed, and
Kubernetes and Swarm playing nicely together thanks to Mesos.  The
networking is Docker bridge networking.  The Swarm master is modified
for multi-tenant use.

Install [Vagrant](https://www.vagrantup.com/) and
[VirtualBox](https://www.virtualbox.org/wiki/Downloads). You will need
at least Vagrant 1.8.4 and VirtualBox 5.0.24.

Checkout this project:

```
git clone git@github.ibm.com:alchemy-containers/openradiant.git
cd openradiant
```

If you are running Ubuntu in your host, you may need to install the following
python packages:

```
sudo apt-get install python-pip python-dev
```

Install ansible:

```
pip install -r requirements.txt
```

Ansible 1.9.6 is reccomended. If you have a different version and run into issues
try the following:

```
pip install --upgrade ansible==1.9.6
```

Create a new cluster with Vagrant:

```
cd examples/vagrant
vagrant up
cd -
```
In case you face any issues, please follow [vagrant troubleshooting](vagrant/README.md)

Deploy OpenRadiant:

```
cd ansible
export ANSIBLE_INVENTORY=../examples/envs/dev-vbox/radiant01.hosts
export ANSIBLE_LIBRARY=$ANSIBLE_INVENTORY
ansible-playbook -v shard.yml -e cluster_name=dev-vbox-radiant01 -e network_kind=bridge
```

Now you can open an SSH connection to the master node:
```
ssh -i ~/.vagrant.d/insecure_private_key vagrant@192.168.10.2
```

On the master you will find both the `kubectl` and `docker`
(currently 1.11) commands on your `$PATH`.

The Swarm master is configured for multi-tenant use.  To prepare to
use it as a tenant, do this on the master:
```
cd; mkdir -p radiant/configs/demo; cat > radiant/configs/demo/config.json <<EOF
{
    "HttpHeaders": {
          "X-Auth-TenantId": "demo"
    }
}
EOF
```

Then you will want these commands on the master:
```
export DOCKER_TLS_VERIFY=""
export DOCKER_CONFIG=~/radiant/configs/demo
export DOCKER_HOST=localhost:2375
```

To get a listing of this tenant's containers, issue the following command on the master:
```
docker ps
```
At first, there will be none.  So create one, like this:
```
docker run --name s1 -d -m 128m busybox sleep 864000
```

Then, get a list of containers, with `docker ps`.  You can inspect its network
configuration from inside, like this:
```
docker exec s1 ifconfig
```

You can create a Kubernetes "deployment" with a command like this:
```
kubectl run k1 --image=busybox sleep 864000
```

The containers in this deployment will be invisible to Swarm because
they lack the label identifying your tenant.  To make containers
visible to Swarm, make a kubernetes pod as follows.  Create a YAML
file prescribing the pod:
```
cat > sleepy-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: sleepy-pod
  annotations:
    containers-annotations.alpha.kubernetes.io: "{ \"com.ibm.radiant.tenant.0\": \"demo\",  \"OriginalName\": \"sleeper\" }"
spec:
  containers:
    - name: sleeper
      image: busybox
      args:
      - sleep
      - "864000"
```

Then create the pod:
```
kubectl create -f sleepy-pod.yaml
```

Then you can watch for it to come up, with
```
kubectl get pod
```

Once it is up, you can inspect its network configuration from inside,
like this:
```
kubectl exec sleepy-pod ifconfig
```

In a similar vein, you can ping one of these containers from the
other.  For example (in which the Swarm container's IP address
is 172.17.0.5):
```
kubectl exec sleepy-pod ping -- -c 2 172.17.0.5
```
For more information see ... (TBD)

### To check the HAproxy statistics using the GUI:

On your local browser, enter the following URL:

master_ip:harproxy_GUI_port/haproxy_stats

Example:
http://192.168.10.2:9000/haproxy_stats (port 9000 is statically assigned)
When prompt for the user_namer:password  use  vagrant:radiantHA

To setup Proxy, follow the steps in [Proxy Setup](#proxy-setup)

### Proxy Setup
For details, please follow [Proxy documentation](../proxy/README.md)
