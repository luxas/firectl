package main

import (
	"fmt"
	"log"
	"net"

	"github.com/docker/libcontainer/netlink"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

/*
ip addr del "$IP/$IPNET" dev eth0

ip link add name br0 type bridge
ip tuntap add dev vm0 mode tap

ip link set br0 up
ip link set vm0 up

ip link set eth0 master br0
ip link set vm0 master br0
*/

func run() error {
	routes, err := netlink.NetworkGetRoutes()
	if err != nil {
		return err
	}
	var defaultIface *net.Interface
	var defaultIPNet net.IPNet
	var defaultIPAddr net.Addr
	for _, route := range routes {
		if route.Default {
			if route.Iface == nil {
				return fmt.Errorf("unexpected default iface was nil")
			}
			if defaultIface != nil {
				return fmt.Errorf("unexpected two default interfaces")
			}
			defaultIface = route.Iface
			break
		}
	}
	if defaultIface == nil {
		return fmt.Errorf("couldn't find default interface")
	}
	addrs, err := defaultIface.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if addr.IP != nil && addr.IP.To4() != nil {
			defaultIPAddr = addr
			break
		}
	}
	fmt.Printf("Would run 'ip addr del %s dev %s'\n", defaultIPAddr.String(), defaultIface.Name)
	return nil
}
