FROM alpine
ARG FIRECRACKER_VERSION=v0.15.0
ADD https://github.com/firecracker-microvm/firecracker/releases/download/${FIRECRACKER_VERSION}/firecracker-${FIRECRACKER_VERSION} /firecracker
# This Dockerfile's context is root of this repo
COPY firectl docker/install.sh /
RUN chmod +x /firecracker /firectl /install.sh
CMD ["/install.sh"]
