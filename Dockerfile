FROM quay.io/prometheus/busybox:latest
LABEL maintainer="Kamil Kayumov <kamil.kayumov@gmail.com>"

ARG TARGETARCH
ARG BIN_DIR=.build/linux-${TARGETARCH}/

COPY ${BIN_DIR}/kafka_exporter /bin/kafka_exporter

EXPOSE 9308

USER nobody

ENTRYPOINT [ "/bin/kafka_exporter" ]
