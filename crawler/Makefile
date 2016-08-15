# These assume you hav Docker installed and running
SHELL = bash
all: build start setupelk

build:
	@echo -ne "Building crawler container..."
	@sudo docker build -t crawler . 1>/dev/null
	@if [ $$? -eq 0 ];then echo "[OK]";else echo "[FAILED]"; fi

start:
	@echo -ne "Starting crawler container..."
	@sudo docker run \
		--privileged \
		--net=host \
		--pid=host \
		--name agentless-crawler\
		-v /cgroup:/cgroup:ro \
		-v /var/lib/docker:/var/lib/docker:ro \
		-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-it -d crawler \
		 --url "http://127.0.0.1:8080/crawler/data"\
		 --features os,disk,process\
		 --frequency 5\
		 --environment radiant\
		 --format json\
		 --crawlmode OUTCONTAINER\
		 1>/dev/null
	@echo "[OK]"

setupelk:
	@echo -ne "Building elk container..."
	@cd agentless-system-crawler && sudo docker build -t crawler-elk -f Dockerfile.elk . 1>/dev/null
	@if [ $$? -eq 0 ];then echo "[OK]";else echo "[FAILED]"; fi
	@echo -ne "Starting elk container..."
	@docker run -d -p 5601:5601 -p 9200:9200 -p 5044:5044 -p 5000:5000 -p 8080:8080 --name elk crawler-elk 1>/dev/null
	@echo "[OK]"
	@#sleep is added to allow time for elk components to bootstrap
	@#before crawler can send data to it
	@sleep 2
	@echo "You can view the crawl data at http://localhost:5601/app/kibana"
	@echo "Please create an index by @timestamp on kibana dashboard"
	
clean:
	@echo -ne "Stopping and removing crawler container..."
	@docker stop agentless-crawler 1>/dev/null
	@docker rm agentless-crawler 1>/dev/null
	@echo "[OK]" 
	@echo -ne "Stopping and removing elk container..."
	@docker stop elk>/dev/null 
	@docker rm elk 1>/dev/null
	@echo "[OK]" 
