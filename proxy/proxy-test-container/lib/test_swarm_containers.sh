#!/bin/bash

source "$(dirname "$0")/common_test_func.sh"

function main {
    local LOGS_PATH="$1"
    local TENANT_ID="$2"
    local PARALLEL="$4"
    local NUM_NO_NET_CONTAINERS=$3
    local NUM_NET_CONTAINERS=0 # No network for now

    init_tests "$TENANT_ID" "swarm" "$LOGS_PATH"
    
    local NONET_COUNTER=1
    while [ $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do

        if [[ $PARALLEL == true ]]; then
            test_creation_container_without_network "$TENANT_ID" $NONET_COUNTER &
            increment_test_count
        else
            test_creation_container_without_network "$TENANT_ID" $NONET_COUNTER
        fi

        let NONET_COUNTER++
    done

    local NET_COUNTER=1
    while [ $NET_COUNTER -le $NUM_NET_CONTAINERS ]; do
    
        if [[ $PARALLEL == true ]]; then
            test_creation_container_with_network "$TENANT_ID" $NET_COUNTER &
            increment_test_count
        else
            test_creation_container_with_network "$TENANT_ID" $NET_COUNTER
        fi

        let NET_COUNTER++
    done
    
    wait

    NONET_COUNTER=1
    while [ $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do

        if [[ $PARALLEL == true ]]; then
            test_deletion_container_without_network "$TENANT_ID" $NONET_COUNTER &
            increment_test_count
        else
            test_deletion_container_without_network "$TENANT_ID" $NONET_COUNTER
        fi
    
        let NONET_COUNTER++
    done
    
    wait

    assert_network_ls 0

    if [[ "$try_net_id" != "" ]]; then
        assert_network_inspect "$try_net_id" 0
    fi 

    # Make sure everything is clean at the end
    assert_ps_a 0
    assert_ps 0

    complete_tests "Docker Containers"
}

function test_creation_container_without_network {
    local TENANT_ID="$1"
    local NONET_COUNTER="$2"

    begin_test_block

    assert_ps 0

    assert_inspect "${TENANT_ID}_test1" 1

    assert_create_nonet "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
    assert_ps 0

    assert_inspect "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0

    assert_stop "${TENANT_ID}_test1" 1 # Should fail because doesn't exist
    assert_stop "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
    assert_ps_a 0
    assert_ps 0

    assert_start "${TENANT_ID}_test1" 1 # Doesn't exist
    assert_start "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0

    assert_ps_a 0

    # Test can't remove without stopping first
    assert_rm "${TENANT_ID}_nonet_test${NONET_COUNTER}" 1
    
    end_test_block
}

function test_deletion_container_without_network {
    local TENANT_ID="$1"
    local NONET_COUNTER="$2"

    if [ $(( $NONET_COUNTER % 2 )) -eq 0 ]; then
        begin_test_block
    
        assert_stop "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        assert_rm "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
        
        end_test_block
    else
        assert_rm_f "${TENANT_ID}_nonet_test${NONET_COUNTER}" 0
    fi
}

function test_creation_container_with_network {
    local TENANT_ID="$1"
    local NONET_COUNTER="$2"

    begin_test_block

    assert_ps 0

    assert_create_net "${TENANT_ID}_net_test${NET_COUNTER}" 0
    assert_ps 0

    assert_inspect "${TENANT_ID}_net_test${NET_COUNTER}" 0

    assert_stop "${TENANT_ID}_net_test${NET_COUNTER}" 0
    assert_ps_a 0
    assert_ps 0

    assert_start "${TENANT_ID}_net_test${NET_COUNTER}" 0

    assert_ps_a 0

    # Test can't remove without stopping first
    assert_rm "${TENANT_ID}_net_test${NET_COUNTER}" 1
    
    end_test_block
}

function assert_ps {
    assert docker ps \
           --equal $1 \
           --log "docker ps"
}

function assert_ps_a {
    assert docker ps -a \
           --equal $1 \
           --log "docker ps -a"
}

function assert_create_nonet {
    assert docker run -d --name "$1" --net none -m 128m mrsabath/web-ms \
           --equal $2 \
           --log "container creation w/o net; name = $1; docker run"
}

function assert_create_net {
    assert docker run -d --name "$1" -m 128m mrsabath/web-ms \
           --equal $2 \
           --log "container creation w/ net; name = $1; docker run"
}

function assert_inspect {
    assert docker inspect "$1" \
           --equal $2 \
           --log "docker inspect; name = $1"
}

function assert_stop {
    assert docker stop "$1" \
           --equal $2 \
           --log "docker stop; name = $1"
}

function assert_start {
    assert docker start "$1" \
           --equal $2 \
           --log "docker start; name = $1"
}

function assert_rm {
    assert docker rm "$1" \
           --equal $2 \
           --log "docker rm; name = $1"
}

function assert_rm_f {
    assert docker rm -f "$1" \
           --equal $2 \
           --log "docker rm -f; name = $1"
}

function assert_network_ls {
    assert docker network ls \
           --equal $1 \
           --log "docker network ls"
}

function assert_network_inspect {
    assert docker network inspect "$1" \
           --equal $2 \
           --log "docker network inspect; name = $1"
}

main "$@"
