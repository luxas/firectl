#!/bin/bash

if [[ ! -d /install ]]; then
    echo "Please volume mount /install to the host directory you want firecracker and firectl copied to"
fi

cp /firectl /firecracker /install
echo "firectl and firecracker installed to the host directory of your choice!"
