FROM alpine
RUN apk add --update iproute2
ARG FIRECRACKER_VERSION=v0.15.2
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /firecracker
# This Dockerfile's context is root of this repo
COPY firectl docker/install.sh /
RUN chmod +x /firecracker /firectl /install.sh
RUN ln -s /firecracker /usr/local/bin/firecracker && ln -s /firectl /usr/local/bin/firectl
CMD ["/install.sh"]
