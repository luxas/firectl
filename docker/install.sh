#!/bin/sh

if [ ! -d /install ]; then
    echo "Please volume mount /install to the host directory you want firecracker and firectl copied to"
    exit 1
fi

cp /firectl /firecracker /install
echo "firectl and firecracker installed to the host directory of your choice!"
