#!/bin/bash

source "$(dirname "$0")/common_test_func.sh"

function main {
    local LOGS_PATH="$1"
    local TENANT_ID="$2"
    local NUM_NO_NET_CONTAINERS=$3
    local NUM_NET_CONTAINERS=0 # No network for now

    init_tests "$TENANT_ID" "swarm" "$LOGS_PATH"
    
    local NONET_COUNTER=1
    while [ $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do
        test_ps 0

        test_inspect "${TENANT_ID}_test1" 1

        test_create_nonet "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        test_ps 0

        test_inspect "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0

        test_stop "${TENANT_ID}_test1" 1 # Should fail because doesn't exist
        test_stop "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        test_ps_a 0
        test_ps 0

        test_start "${TENANT_ID}_test1" 1 # Doesn't exist
        test_start "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0

        test_ps_a 0

        # Test can't remove without stopping first
        test_rm "${TENANT_ID}_nonet_test${NONET_COUNTER}" 1

        let NONET_COUNTER=NONET_COUNTER+1
    done

    local NET_COUNTER=1
    while [ $NET_COUNTER -le $NUM_NET_CONTAINERS ]; do
        test_ps 0

        test_create_net "${TENANT_ID}_net_test${NET_COUNTER}" 0
        test_ps 0

        test_inspect "${TENANT_ID}_net_test${NET_COUNTER}" 0

        test_stop "${TENANT_ID}_net_test${NET_COUNTER}" 0
        test_ps_a 0
        test_ps 0

        test_start "${TENANT_ID}_net_test${NET_COUNTER}" 0

        test_ps_a 0

        # Test can't remove without stopping first
        test_rm "${TENANT_ID}_net_test${NET_COUNTER}" 1

        let NET_COUNTER=NET_COUNTER+1
    done

    NONET_COUNTER=1
    while [ $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do
        if [ $(( $NONET_COUNTER % 2 )) -eq 0 ]; then
            test_stop "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
            test_rm "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        else
            test_rm_f "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        fi

        let NONET_COUNTER=NONET_COUNTER+1
    done

    test_network_ls 0

    if [[ "$try_net_id" != "" ]]; then
        test_network_inspect "$try_net_id" 0
    fi 

    # Make sure everything is clean at the end
    test_ps_a 0
    test_ps 0

    complete_tests "Docker Containers"
}

function test_ps {
    assert docker ps \
           --equal $1 \
           --log "docker ps"
}

function test_ps_a {
    assert docker ps -a \
           --equal $1 \
           --log "docker ps -a"
}

function test_create_nonet {
    assert docker run -d --name "$1" --net none -m 128m mrsabath/web-ms \
           --equal $2 \
           --log "container creation w/o net; name = $1; docker run"
}

function test_create_net {
    assert docker run -d --name "$1" -m 128m mrsabath/web-ms \
           --equal $2 \
           --log "container creation w/ net; name = $1; docker run"
}

function test_inspect {
    assert docker inspect "$1" \
           --equal $2 \
           --log "docker inspect; name = $1"
}

function test_stop {
    assert docker stop "$1" \
           --equal $2 \
           --log "docker stop; name = $1"
}

function test_start {
    assert docker start "$1" \
           --equal $2 \
           --log "docker start; name = $1"
}

function test_rm {
    assert docker rm "$1" \
           --equal $2 \
           --log "docker rm; name = $1"
}

function test_rm_f {
    assert docker rm -f "$1" \
           --equal $2 \
           --log "docker rm -f; name = $1"
}

function test_network_ls {
    assert docker network ls \
           --equal $1 \
           --log "docker network ls"
}

function test_network_inspect {
    assert docker network inspect "$1" \
           --equal $2 \
           --log "docker network inspect; name = $1"
}

main "$@"
