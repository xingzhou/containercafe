**Step 1**: Use GVM to install go version 1.4.1+ (Jenkins doesn't work with go1.6 yet, so use go1.5 at most)
- `gvm install go1.5`
- `gvm use go1.5`

**Step 2**: Get code
- `go get github.ibm.com/alchemy-containers/remoteabac`

**Step 3**: Compile and install code
- `godep go install github.ibm.com/alchemy-containers/remoteabac/cmd/remoteabac`

**Step 4**: Run code
- `$GOBIN/remoteabac --address=:8888 --tls-cert-file=cert.pem --tls-private-key-file=key.pem --authorization-policy-file=etcd@http://<ip>:<port>/abac-policy`

**Step 5**: Add/delete user
- `curl -k https://<ip>:<port>/users` to list all the users
- `curl -XPOST -k https://<ip>:<port>/user/admin?privileged=true` to add an admin user
- `curl -XPOST -k https://<ip>:<port>/user/haih/haihns` to add user `haih` to namespace `haihns`
- `curl -XDELETE -k https://<ip>:<port>/user/haih` to delete user `haih`
