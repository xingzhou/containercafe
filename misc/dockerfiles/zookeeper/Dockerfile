FROM alpine:3.4

ENV ZK_VER zookeeper-3.4.8

ENV JAVA_HOME=/usr/lib/jvm/default-jvm

RUN  apk add --no-cache openjdk7 bash && \
 ln -sf "${JAVA_HOME}/bin/"* "/usr/bin/" &&\
 wget http://mirrors.ocf.berkeley.edu/apache/zookeeper/$ZK_VER/$ZK_VER.tar.gz &&\
 mkdir opt &&\
 tar xzvf $ZK_VER.tar.gz -C /opt &&\
 mkdir -p /opt/$ZK_VER/logs &&\
 cp /opt/$ZK_VER/conf/zoo_sample.cfg /opt/$ZK_VER/conf/zoo.cfg &&\
 rm $ZK_VER.tar.gz

# Make soft link to /opt/zookeeper, so we did not need to change it
RUN ln -s /opt/$ZK_VER /opt/zookeeper
ENV ZK_HOME /opt/zookeeper
ENV PATH $ZK_HOME/bin:$PATH
COPY bootstrap.sh /opt/zookeeper

ENTRYPOINT ["/opt/zookeeper/bootstrap.sh"]
