# Proxy Environment Setup and Development
Steps below explain the basic configuration and setup to run Proxy standalone.


## Step 1: Authentication config file

This project contains example `creds.json` file. This supercedes the authentication done by CCSAPI.
In other words if authentication succeeds against a record in this fie, CCSAPI will not be contacted.

Each record is a json object with the following syntax:
```json
{
  "Status":200, "Node":"10.140.28.132:2379", "Docker_id":"", "Container":"", "Swarm_shard":true,
  "Tls_override":true, "Space_id":"85cdc7e0-32d8-4552-9bae-907c3f1d98d9", "Reg_namespace":"swarm", "Apikey":"c3d87893a5b7f56991fd328f655f25cce286591c3ce4a558",
  "Orguuid":"9013217d-0abf-40fe-bd35-bb625066408c",
  "Userid":"924fc412d1004528b90007e898aeb0d8"
}
```

The important fields are:
- Space_id: Matching happens against this field. X-Auth-Project-Id header in the request must have identical Space_id for a match to happen.
- Status: If it is not 200, auth will fail.
- Node: address and port to redirect request to
- Tls_override: Do not use tls (even if proxy is configured to use tls)
- Swarm_shard: Radiant shard if true, otherwise nova-docker shard.

The other fields can be ignored in a standalone setup.

If running as container, the location of the `creds.json` file is specified in the following Dockerfile line:
```
ENV stub_auth_file="/opt/tls_certs/creds.json"
```

If NOT running as a docker container, then the default `creds.json` file location is in the current directory.

## Step 2: Compile  and run
Use `build.sh` script to compile, then run directly as a process on the host by invoking bin/hijack
```shell
chmod +x build.sh
./build.sh
bin/hijack
```

Alternatively, use `builddocker.sh` to build a docker container named hijack, and use
`rundocker.sh` to start the container with the right -v and -p options.

**NOTE**: All default config options are defined in the Dockerfile, and can be overridden using the docker -e option on startup.

**NOTE**: If running standalone, i.e., using file auth, be mindful of where creds.json should be placed. See comments in Step 2 above.

## Step 3: Setup client headers to execute HTTP requests
Your request must carry these the specified HTTP headers,otherwise File
authentication will fail and CCSAPI will be invoked:
- X-Auth-Project-Id: <Space_id>
where <Space_id> is the Bluemix Space Id for that user.
- X-Auth-TenantId: <Space_id>
where <Space_id> is the Blumix Space Id required by Swarm cluster.

For debugging purposes it is handy to crate docker config file to reference with
docker calls. Use either `~/.docker/config.json` for all the docker calls or
custom `~/projects/radiant/proxy/config.json` to reference with explicit calls

Content of the file must include the above headers.
Example:
```json
{
  "HttpHeaders": {
    "X-Auth-Project-Id": "abcdefg",
    "X-Auth-TenantId": "abcdefg"
  }
}
```
Then reference this file when executing docker calls:
`DOCKER_TLS_VERIFY=""  docker -H $PROXY --config ~/projects/radiant/proxy/ ps`


## Step 4: Invoke client HTTP calls:
The requests can be executed using `docker` client, directly using `curl` or
for Kubernetes using `kubectl`

### Using docker client
Setup the location of the Proxy server endpoint. When running locally, its localhost
on port 8087

```shell
export PROXY=localhost:8087
```
Sample calls:
```shell
DOCKER_TLS_VERIFY="" docker -H $PROXY --config ~/projects/radiant/proxy/ ps
DOCKER_TLS_VERIFY="" docker -H $PROXY --config ~/projects/radiant/proxy/ run -d --name MS-test1 -m 128m --env test=1 10.140.132.215:5001/mrsabath/web-ms
DOCKER_TLS_VERIFY="" docker -H $PROXY --config ~/projects/radiant/proxy/ inspect <cont_id>
```


## Errors and Hints

`An error occurred trying to connect: Get https://localhost:8087/v1.21/images/4a419cdeaf69/json: tls: oversized record received with length 20527[]`

In order to fix this problem, use `DOCKER_TLS_VERIFY=""` prefix for running 'docker' command
