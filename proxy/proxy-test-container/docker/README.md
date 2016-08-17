# Testing, for API Proxy and tests running as a Container
Steps below explain configuration and setup to run Swarm & Kubernetes testing as a container, while API Proxy is running as a container.

Note: Steps below detail multi-user testing. Testing can also be done with a single-user.

## Step 1: Build the Docker image

Move to the proxy-test-container directory and run the build script
```bash
cd ~/workspace/openradiant/proxy/proxy-test-container
./docker/build.sh
```

The build script has 3 possible flags (all optional):
 1. `-e <environment>`: the environment name.
 2. `-d <docker_version>`: the Docker version that will be installed.
 3. `-k <kubectl_version>`: the Kubernetes version that will be installed.

## Step 2: Certificates

After starting the API Proxy, execute the script to create test tenants:
```bash
docker exec api-proxy /api-proxy/make_TLS_certs.sh test1 radiant01 192.168.10.2
docker exec api-proxy /api-proxy/make_TLS_certs.sh test2 radiant01 192.168.10.2
docker exec api-proxy /api-proxy/make_TLS_certs.sh test3 radiant01 192.168.10.2
```
To view all the accounts valid for this proxy: 
```bash
docker exec api-proxy cat /api-proxy/creds.json
```
The certificate creation script will output few export statements. For example:
```bash
# Setup docker environment:
export DOCKER_HOST=localhost:8087
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=~/.openradiant/envs/dev-vbox/radiant01/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I

# Setup kubernetes environment:
export KUBECONFIG=~/.openradiant/envs/dev-vbox/radiant01/7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I/kube-config
```
For each user you need to remember the API key (e.g. `7uJNzJqK5T33A4j9XkH6Fd1dQwCza0zHGHeFokmRJOWfz87I`).

## Step 3: Run the tests

The test image has 6 possible flags (all optional except for `-t`):
 1. `-l <proxy_location>`: either local or dev-mon01. Indicates where proxy should be run - on user's local machine, or in the "dev-mon01" data center. For container testing, the proxy runs locally. Local is the default.
 2. `-n <network_id>`: id of the network element to be inspected. Currently not supported; therefore, default is empty.
 3. `-t <shard>:<tenant_id>:<api_key>`: you can pass as many tenants as you want.
 4. `-k <test_kube?>`: `true` or `false`. Default is true; kube tests will be run. 
 5. `-c <num_containers>`: the total number of containers to be created, all without network. If no argument passed, default value is 5.
 6. `-p <num_pods>`: the total number of pods to be created. If no argument is passed, default value is 5.

Suppose each test wants to examine a lifecycle of 3 containers and 3 pods.
Then run:
```bash
./docker/run.sh -c 3 -p 3 \
                -t radiant01:test1:nVdIhfJKFtEM2G9hpkJI5EVZ5VGeFNLBoFBA2B6zJqaSZ71W \
                -t radiant01:test2:JpjnQqnDQQCnNlc9bWJJoznhjr4awLdVe10B45LRCE31CqDh \
                -t radiant01:test3:PBJ7ZKbi0hDSBFStFY5K9TnNdojEhUZ1goE1Swn3G6fle5iR
```

## Step 4: Results

Both Docker Swarm & Kubernetes tests will be executed. Summaries of the test results will be printed to console.
```
The env is dev-vbox
The log tag is 2016-08-16_02-11-29
Completed: 3/3

Test results for tenant test1:
Docker Containers Test summary:
Total = 46 tests
Passed = 46 tests
Failed = 0 tests
Total time = 55 sec
Kube Test summary:
Total = 22 tests
Passed = 22 tests
Failed = 0 tests
Total time = 21 sec

Test results for tenant test2:
Docker Containers Test summary:
Total = 46 tests
Passed = 46 tests
Failed = 0 tests
Total time = 54 sec
Kube Test summary:
Total = 22 tests
Passed = 22 tests
Failed = 0 tests
Total time = 21 sec

Test results for tenant test3:
Docker Containers Test summary:
Total = 46 tests
Passed = 46 tests
Failed = 0 tests
Total time = 54 sec
Kube Test summary:
Total = 22 tests
Passed = 22 tests
Failed = 0 tests
Total time = 21 sec

0 failed tests.
```

To show the full log report run:
```bash
./docker/log.sh --latest
```

Otherwise it's possible to print a specific log identified by the log tag (the log timestamp):
```bash
./docker/log.sh -t 2016-08-16_02-11-29
```

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

