
* To test remote abac authorization

$ curl -k -XPOST -d @<input file> https://<ip>:<port>/authorize

* To test user management (this is REST-ful)

** To add regular user

$ curl -k -XPUT https://<ip>:<port>/user/<username>/<namespace>

** To add privileged user

$ curl -k -XPUT https://<ip>:<port>/user/<username>?privileged=true

** To show user policy

$ curl -k -XGET https://<ip>:<port>/user/<username>
$ curl -k -XGET https://<ip>:<port>/user/<username>/<namespace>

** To delete user policy in all namespaces

$ curl -k -XDELETE https://<ip>:<port>/user/<username>

** To delete user policy in a specific namespace

$ curl -k -XDELETE https://<ip>:<port>/user/<username>/<namespace>
