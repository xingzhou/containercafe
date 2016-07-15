
* To test remote abac authorization

$ curl -k -XPOST -d @<input file> https://<ip>:<port>/authorize

* To test user management (this is REST-ful)

** To add regular user

$ curl -k -XPUT https://<ip>:<port>/<username>/<namespace>

** To add privileged user

$ curl -k -XPUT https://<ip>:<port>/<username>?privileged=true

** To show user policy

$ curl -k -XGET https://<ip>:<port>/<username>
$ curl -k -XGET https://<ip>:<port>/<username>/<namespace>

** To delete user policy in all namespaces

$ curl -k -XDELETE https://<ip>:<port>/<username>

** To delete user policy in a specific namespace

$ curl -k -XDELETE https://<ip>:<port>/<username>/<namespace>
