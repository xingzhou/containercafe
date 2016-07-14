# Proxy Environment Setup and Development
Steps below explain the basic configuration and setup to run Proxy. This setup
assumes the target system already has Docker installed, the Proxy will be deployed
locally to manage Swarm and Kubernetes instances also installed locally, using
the ansible scripts described in [Quick Start](../README.md#quick-start)
For other deployments, one can still follow these steps, with adjustments provided below.

## Step 1: Get proxy code
If you have not done this already, clone the repository:

```
git clone https://github.ibm.com/alchemy-containers/openradiant.git
```
Then build and deploy it:
```
cd openradiant/proxy
./builddocker.sh
./rundocker.sh
```

This will start the Proxy as a container named `hjproxy`, running in the current terminal.
If you want to run the Proxy as background container (daemon), use the `-d` flag:
```
./rundocker.sh -d
```
The Proxy container can be seen along with its logs:
```
docker ps
docker logs -f hjproxy
```

If you don't want to run the Proxy as a container, you can run directly as GoLang
application. This requires Go libraries installed, proxy code compiled and added
to go path.
Review [Dockerfile](dockerize/Dockerfile) for more details.
Setup the environment and start the application:
```
source ./set_local_env.sh
./start_proxy.sh
```


## Step 2: Setting up the tenant
From another terminal go back to the proxy directory and then execute the script
to create TLS certificates and API key for the given tenant. The example is
using tenant `test1`
```
docker ps
docker exec hjproxy /hijack/make_TLS_certs.sh test1
```

This command will display the details about the newly created TLS certs, including
their location. Certificates are created on the mounted volumes, therefore they
are available inside the proxy container and for the current terminal.

You can create as many tenants as you like.

At the bottom of the output the script displays the docker environment setup for
the newly created tenant, including the location of the certs. Here is the sample
output:
```
Execute the following:
export DOCKER_HOST=localhost:8087
export DOCKER_TLS_VERIFY=1
export DOCKER_CONFIG=dockerize/OpenRadiant/1z5AmY6uDzqBT65mkPtffEhOutcxs3sghn9S9LrXfAOztCpR
export DOCKER_CERT_PATH=dockerize/OpenRadiant/1z5AmY6uDzqBT65mkPtffEhOutcxs3sghn9S9LrXfAOztCpR
```
Copy the last four lines and paste in a new terminal. Make sure you are in
`openradiant/proxy` directory.

Now you should be able to execute commands against OpenRadiant account that you
just created:

```
docker ps
docker run -d --name test --net none -m 128m mrsabath/web-ms
docker inspect test
```
To run the proxy against a different OpenRadiant deployment or cluster, change
the IP of the TARGET_SERVER value to this new location in `make_TLS_certs.sh` script.
By default it is set to the local server running on `radiant2` in vagrant:
`TARGET_SERVER="192.168.10.2"`, then rebuild and redeploy the Proxy.
The newly created tenant will be pointing at the new location.
You can also manually change the values in "/hijack/creds.json" file that lists
all the valid tenants.  
To modify the value in the running container:
```
docker exec -it hjproxy /bin/bash
```

**NOTE**: All default config options are defined in the [Dockerfile](dockerize/Dockerfile),
and can be overridden using the docker -e option on [startup](rundocker.sh)


## Errors and Hints

`An error occurred trying to connect: Get https://localhost:8087/v1.21/images/4a419cdeaf69/json: tls: oversized record received with length 20527[]`

In order to fix this problem, use `DOCKER_TLS_VERIFY=""` prefix for running 'docker' command



## Running Test Scripts

* Test scripts are located in `proxy-test` folder, and should be run from this directory.

* Execute `source setup_local.sh` or `source setup_CCSAPI.sh`. Former tests local file authentication, latter, CCSAPI authentication.

* Run with the command `./test_containers.sh da07` (can also be run with no argument).

* Results delineated in file `test_containers_results.txt`.
