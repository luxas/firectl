FROM alpine
ARG FIRECRACKER_VERSION=v0.14.0
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /firecracker
# This Dockerfile's context is root of this repo
COPY firectl docker/install.sh /
CMD ["/install.sh"]
