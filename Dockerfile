FROM hub.iot.chinamobile.com/library/centos:7.4

WORKDIR /data/app
COPY ./build/ocb-mep-mepmgr /data/app/ocb-mep-mepmgr
COPY  ./conf /data/app/conf

ENTRYPOINT ["/data/app/ocb-mep-mepmgr"]

