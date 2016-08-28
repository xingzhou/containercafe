Agentless System Crawler 
========================

**Disclaimer:**
---------------

```
"The strategy is definitely: first make it work, then make it right, and, finally, make it fast."
```
The current state of this project is in the middle of "make it right".


Run Crawler in OpenRadiant platform**Prereqs and Building:**
-----------

Simplest trial is to run crawler in a container. 
Crawler also supports an output emmitter to standard elk-stack. 
So you can view crawl data on kibana dashboard.

To setup crawler with elk-stack:

```bash
$ make
Building crawler container...[OK]
Building elk container...[OK]
Starting elk container...[OK]
Starting crawler container...[OK]
You can view the crawl data at http://localhost:5601/app/kibana
Please create an index by @timestamp on kibana dashboard

```

Be default crawler collects only __os-info__,__cpu__, __memory__ metrics for 
containers every 1 min. For these crawl features the size of output produced
is less than 1 KB.  

If you want to crawl more features (packages, files, configs etc.) you can configure
crawler accordingly. For more details see [agentless-system-crawler](https://github.com/cloudviz/agentless-system-crawler)

Stop  crawler**Prereqs and Building:**
-----------
To stop crawling, simply do _make clean_. It will stop crawler and elk containers and remove them.

```bash
$make clean
Stopping and removing crawler container...[OK]
Stopping and removing elk container...[OK]
```

