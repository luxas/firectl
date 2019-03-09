#!/bin/sh

set -x

DEV=eth0
IFACE=$(ip link | grep -o "$DEV@[[:alnum:]]*")
IP=$(ip addr show "$DEV" | grep "inet" | awk '{ print $2 }')
BRD=$(ip addr show "$DEV" | grep "inet" | awk '{ print $4 }')
MAC=$(echo $FQDN|md5sum|sed 's/^\(..\)\(..\)\(..\)\(..\)\(..\).*$/02:\1:\2:\3:\4:\5/')
GATEWAY=172.17.0.1
SOURCE=172.17.0.254

IPNOSUBNET=$(echo "$IP" | cut -d/ -f1)

echo "Device: $DEV"
echo "Interface: $IFACE"
echo "IP: $IP"
echo "Broadcast: $BRD"

echo "Deleting ip from container!"
ip a del "$IP" brd "$BRD" dev "$DEV"

echo "Got these args: $@"

echo "Creating bridge stuff..."
ip link add name br0 type bridge
ip link set br0 up
ip link set "$DEV" master br0
#ip addr change ${SOURCE}/24 dev br0

ip tuntap add dev vm0 mode tap
ip link set vm0 up
ip link set vm0 master br0

/fc-dhcpd --source-ip ${SOURCE} --gateway-ip ${GATEWAY} --client-ip ${IPNOSUBNET} --client-mac ${MAC} & > dhcp.log

#sed -e "s|MAC|${MAC}|g;s|IP|${IPNOSUBNET}|g;s|GATEWAY|${GATEWAY}|g" /etc/dhcpd.conf.tmpl > /etc/dhcpd.conf
#cat /etc/dhcpd.conf
#dhcpd br0 &

echo "Launching firectl in 5 seconds..."

sleep 5

/firectl --root-drive /rootfs.ext4 --kernel /vmlinux --firecracker-binary /firecracker --add-network "vm0/${MAC}"
