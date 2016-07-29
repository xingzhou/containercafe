# OpenRadiant API Wrapper - Proxy
Proxy is the component of OpenRadiant that intercepts the communication between
the clients (Docker or Kubernetes) and the OpenRadiant cluster, using HTTP session
hijacking. It validates the tenant
and provided TLS certificates. It also redirects to proper shard when cluster
sharing is used.


![Image of Proxy](media/2016-07.OpenRadiantProxy.png)

![Image of Proxy details](media/2016-05.Proxy-details.png)

* Single integration point for all the services
* Handles the multi-tenant authentication (framework to support various auth interfaces)
* Supports common certification, creates tenants and their TLS certificates
* Stateless microservice deployed as a container
* Handles re-direct to appropriate framework API handler
* Validates and filters the requests; Masquerades the internal information that should not be public
* Supports sharding
* Can execute quota validation across multiple frameworks and shards
* Executes annotation injections to support viewing K8s pods via Docker interface
* API version mapping
* API Extensions:
  * Quota management
  * Health checking, status, HA and informational APIs
  * Metering and chargeback management
* Prototyping new features (e.g. Powerstrip for Docker)
