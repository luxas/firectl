

docker run -it --privileged --device /dev/kvm -v $(pwd)/vmlinux:/vmlinux -v $(pwd)/rootfs.ext4:/rootfs.ext4 --name fc$(date +%s) fc


sed -i "s|eth0 inet manual|eth0 inet dhcp|" /etc/network/interfaces
/etc/init.d/networking restart

cat > /etc/systemd/network/dhcp.network <<EOF
[Match]
Name=eth0

[Network]
DHCP=ipv4
EOF


systemctl daemon-reload && systemctl restart systemd-networkd

LinkLocalAddressing=ipv4