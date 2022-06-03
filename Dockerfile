FROM ubuntu:20.04

WORKDIR /var/opt

COPY ./metrics-server-prometheus-exporter ./

RUN apt-get update \
    && apt-get install -y curl \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.22.5/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

EXPOSE 9100

ENTRYPOINT [ "/var/opt/metrics-server-prometheus-exporter" ]