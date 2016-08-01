# Proxy Testing 
Steps below explain configuration & setup to run Swarm & Kube testing. <br />
Note: Steps below detail multi-user testing. Testing can also be done with a single-user. 


## Step 1: Terminals
Open as many terminals as users desired. In this example, we will have 3 users - test1, test2, test3. Therefore, open 3 different terminal windows. 


## Step 2: Certificates
In seperate terminals, create certificates for each user. Execute: <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test1` <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test2` <br />
`docker exec hjproxy /hijack/make_TLS_certs.sh test3`. <br />

This will create the certificates, and setup the configurations necessary to run the tests for the given users. 



## Step 3: Environment Variables
The certificate creation script will output a few export statements; for example: <br />
`export DOCKER_HOST=localhost:8087
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=dockerize/OpenRadiant/YlCajJeJ3IJeBwTXVhQ7jCqiTC437NjlZooFBWuST2MqWk9q` <br />

Copy and paste the first 2 lines. For the variable `DOCKER_CERT_PATH`, first execute <br />
`cd ..`
`pwd` <br />
Copy the resulting path, and paste it before `dockerize/OpenRadiant/YlCajJeJ3IJeBwTXVhQ7jCqiTC437NjlZooFBWuST2MqWk9q`. The final command should look something like 
`export DOCKER_CERT_PATH=/Users/atarng/workspace/openradiant/proxy/dockerize/OpenRadiant/YlCajJeJ3IJeBwTXVhQ7jCqiTC437NjlZooFBWuST2MqWk9q`. 
Execute this command. <br />

For Kubernetes, there will be something like <br />
`export KUBECONFIG=dockerize/OpenRadiant/YlCajJeJ3IJeBwTXVhQ7jCqiTC437NjlZooFBWuST2MqWk9q/kube-config`. 

Do the same copy and pasting of the path to the outer directory before the presented path. The final command will look something like
`export KUBECONFIG=/Users/atarng/workspace/openradiant/proxy/dockerize/OpenRadiant/YlCajJeJ3IJeBwTXVhQ7jCqiTC437NjlZooFBWuST2MqWk9q/kube-config`.
Execute this command. <br />  


Do this for all users in their respective terminals. 



## Step 4: Run Test Script
In each corresponding window, run the test script, `test_containers.sh`. 

`test_containers.sh` has 6 possible flags (all optional):
1) -l (proxy_Location): either local or dev-mon01. Indicates where proxy should be run - on user's local machine, or in the "dev-mon01" data center. For container testing, the proxy runs locally. Local is the default. 
2) -n (network_id): id of the network element to be inspected. Currently not supported; therefore, default is empty. 
3) -t (tenant_id): the tenant id to be used in testing. Default is `test1`. 
4) -k (test_kube): `true` or `false`. Default is true; kube tests will be run. 
5) -c (num_containers): the total number of containers to be created, all without network. If no argument passed, default value is 5. 
6) -p (num_pods): the total number of pods to be created. If no argument is passed, default value is 5. 


Suppose each user wants to examine the network element `kitties`, and create 3 containers and 3 pods each. Then, in each corresponding window, run: <br />
`./test_containers.sh -n kitties -t test1 -c 3 -p 3` <br />
`./test_containers.sh -n kitties -t test2 -c 3 -p 3` <br />
`./test_containers.sh -n kitties -t test3 -c 3 -p 3`.


## Step 5: Results

Both Docker Swarm & Kube tests will be executed. Summaries of the test results can be found in the `logs` folder, under `<tenant_id>_test_kube_pods_results_<timestamp>.log` for the Kube tests, and `<tenant_id>test_swarm_results_<timestamp>.log` for the Swarm tests. 

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


