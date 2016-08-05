The flannel daemon picks a subnet for this host and writes it, along
with the correct MTU to use, into envar settings in
`/run/flannel/subnet.env`.  These need to get into the configuration
of Docker.

Coreos publishes a Docker container image
(`quay.io/coreos/flanneld-amd64`) containing the flannel daemon ---
but pulling it requires credentials.

The approach taken here is to create a system service to run flanneld
and make the Docker engine service depend on it.  At first we are only
supporting upstart.
