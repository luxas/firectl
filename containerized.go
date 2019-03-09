package main

import (
	"net"
)

/*
ip r list src 172.17.0.3

ip addr del "$IP" dev eth0

ip link add name br0 type bridge
ip tuntap add dev vm0 mode tap

ip link set br0 up
ip link set vm0 up

ip link set eth0 master br0
ip link set vm0 master br0
*/

func PrepareContainerNetworking() error {
	commands := []string{}
}

func getNetworkStats() (net.IP, net.IPNet,  {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		return err
	}
	var ipaddr net.IP
	for _, addr := range iface.Addrs() {
		ip4 := net.ParseIP(addr.String()).To4()
		if ip4 != nil {
			ipaddr = ip4
			break
		}
	}

}
