FROM harbor.haodai.net/base/alpine:3.7cgo

WORKDIR /app

MAINTAINER wenzhenglin(http://g.haodai.net/wenzhenglin/service-prober)

# curl http://fs.haodai.net/soft/uploadapi -F file=@cloudprober -F truncate=yes
RUN wget http://fs.haodai.net/soft/cloudprober -O /bin/cloudprober && \
        chmod +x /bin/cloudprober

COPY service-prober /app
COPY cloudprober.cfg /app
COPY dockerstart.sh /app/start.sh

CMD /app/start.sh
# ENTRYPOINT ["/app/start.sh"]

EXPOSE 9314