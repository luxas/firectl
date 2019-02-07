#!/bin/sh

if [ ! -d /install ]; then
    echo "Please volume mount /install to the host directory you want firecracker and firectl copied to"
    echo "Example: docker run -it -v /usr/local/bin:/install luxas/firectl"
    echo "Alternatively, you can run /firectl and /firecracker directly from the container"
    exit 1
fi

cp /firectl /firecracker /install
echo "firectl and firecracker installed to the host directory of your choice!"
