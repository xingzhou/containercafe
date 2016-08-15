FROM python:2.7

WORKDIR /crawler

COPY agentless-system-crawler/requirements.txt /crawler/requirements.txt
RUN pip install -r requirements.txt

ADD agentless-system-crawler/crawler /crawler
COPY plugins/* /crawler/plugins/

ENTRYPOINT [ "python2.7", "crawler.py" ]
