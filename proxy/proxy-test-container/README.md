# Testing, for Proxy running as a Container
Steps below explain configuration & setup to run Swarm & Kubernetes testing, while Proxy is running as a container. <br />
Note: Steps below detail multi-user testing. Testing can also be done with a single-user. 


## Step 1: Terminals
Open as many terminals as test tenants desired. In this example, we will have 3 tenants - test1, test2, test3.
Therefore, open 3 different terminal windows. cd to location where you have cloned this openradiant project. 
Then cd to proxy directory.
```
cd ~/workspace/openradiant/proxy
```

## Step 2: Certificates
In every terminal, execute the script to create test tenants: <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test1` <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test2` <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test3`. <br />

If you are using multiple openradiant shards, your master cluster is different from 
default VIP: `192.168.10.2`, pass that VIP as the 2nd argument. For example, if address is 
`192.168.10.10`, execute `docker exec hjproxy /hijack/make_TLS_certs.sh test1 192.168.10.10`. 
Otherwise, the default address `192.168.10.2` will be used. 

This will create the certificates for this tenant and the configurations necessary to run the tests. 
To view all the accounts valid for this proxy: 
```
docker exec hjproxy cat /hijack/creds.json
```

## Step 3: Environment Variables
The certificate creation script will output few export statements. For example: <br />
```
# Setup docker environment:
export DOCKER_HOST=localhost:8087
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I

# Setup kubernetes environment:
export KUBECONFIG=dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I/kube-config
```
Copy and paste the first 2 lines. For the variable `DOCKER_CERT_PATH`, first execute <br />
`pwd` <br />
Copy the resulting path, and paste it before `dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I`
The final command should look something like:
`export DOCKER_CERT_PATH=/Users/atarng/workspace/openradiant/proxy/dockerize/OpenRadiant/dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I`. 
Execute this command. <br />

For Kubernetes, there will be something like <br />
`export KUBECONFIG=dockerize/OpenRadiant/dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I/kube-config`. 

Do the same copy and pasting of the path to the outer directory before the presented path. The final command will look something like
`export KUBECONFIG=/Users/atarng/workspace/openradiant/proxy/dockerize/OpenRadiant/dockerize/OpenRadiant/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I/kube-config`.
Execute this command. <br />  


Do this for all users in their respective terminals. 



## Step 4: Run Test Script
`cd proxy-test-container`
In each corresponding window, run the test script, `test_containers.sh` 

`test_containers.sh` has 6 possible flags (all optional):
1) -l (proxy_Location): either local or dev-mon01. Indicates where proxy should be run - on user's local machine, or in the "dev-mon01" data center. For container testing, the proxy runs locally. Local is the default. 
2) -n (network_id): id of the network element to be inspected. Currently not supported; therefore, default is empty. 
3) -t (tenant_id): the tenant id to be used in testing. Default is `test1`. 
4) -k (test_kube): `true` or `false`. Default is true; kube tests will be run. 
5) -c (num_containers): the total number of containers to be created, all without network. If no argument passed, default value is 5. 
6) -p (num_pods): the total number of pods to be created. If no argument is passed, default value is 5. 


Suppose each test wants to examine a lifecycle of 3 containers and 3 pods.
Then, in each corresponding window, run: <br />
`./test_containers.sh -t test1 -c 3 -p 3` <br />
`./test_containers.sh -t test2 -c 3 -p 3` <br />
`./test_containers.sh -t test3 -c 3 -p 3`.


## Step 5: Results

Both Docker Swarm & Kubernetes tests will be executed. Summaries of the test results can be found in the `logs` folder, under `<tenant_id>_test_kube_pods_results_<timestamp>.log` for the Kube tests, and `<tenant_id>test_swarm_results_<timestamp>.log` for the Swarm tests. 

Each line of the results file will look similar to this: 
`20160705.144709,1,test1,swarm,test1,Docker ps,PASS,OK`. 
Each field is comma-seperated, and represents the following:

1. Timestamp of when the command was executed. Read as: YYYYMMDD.HHMMSS. 
2. Execution time of the command, in seconds. 
3. Tenant_id 
4. Test type: either docker `swarm` or `kube`. 
5. Test number (sequential). 
6. Description of test run. 
7. Result: either `PASS` or `FAIL`. 
8. If test passed, will just be `OK`. Otherwise, if fail, further information on what went wrong (the error displayed to the user).


