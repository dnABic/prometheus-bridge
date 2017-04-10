FROM scratch

ADD prometheus-bridge /usr/sbin/prometheus-bridge

EXPOSE 9091

ENTRYPOINT ["prometheus-bridge"]
