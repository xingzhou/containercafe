#!/bin/bash

docker run -v ~/.openradiant/envs:/tests/certs:ro -v api-proxy-tests-logs:/tests/logs --net="host" api-proxy-tests "$@"
