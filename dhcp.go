package main

import (
	"fmt"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
)

func run() error {
	gatewayIP := net.ParseIP(*gatewayIPFlag)
	if gatewayIP == nil {
		return fmt.Errorf("--gateway-ip is invalid")
	}
	clientIP := net.ParseIP(*clientIPFlag)
	if clientIP == nil {
		return fmt.Errorf("--client-ip is invalid")
	}
	subnetMask := net.ParseIP(*subnetMaskFlag)
	if subnetMask == nil {
		return fmt.Errorf("--subnet-mask is invalid")
	}
	clientMAC, err := net.ParseMAC(*clientMACFlag)
	if err != nil {
		return fmt.Errorf("--client-mac is invalid: %v", err)
	}
	duration, err := time.ParseDuration("168h")
	if err != nil {
		return err
	}
	d := 2 * time.Month
	handler := &DHCPHandler{
		gatewayIP:     gatewayIP.To4(),
		clientIP:      clientIP.To4(),
		clientMAC:     clientMAC,
		dnsServer:     net.IP{1, 1, 1, 1},
		subnetMask:    subnetMask.To4(),
		leaseDuration: duration,
	}

}

type DHCPServer struct {
	GatewayIP     net.IP
	ClientIP      net.IP
	ClientMAC     net.HardwareAddr
	DNSServer     net.IP
	SubnetMask    net.IP
	LeaseDuration time.Duration
	Debug         bool
}

// DHCPServer implements the dhcp.Handler interface
var _ dhcp.Handler = &DHCPServer{}

// ListenAndServe starts the DHCP server and blocks forever
func (s *DHCPServer) ListenAndServe(iface string) error {
	if err := s.DefaultAndValidate(); err != nil {
		return err
	}
	listener, err := conn.NewUDP4BoundListener(iface, ":67")
	if err != nil {
		return err
	}
	return dhcp.Serve(listener, s)
}

// DefaultAndValidate defaults the server's fields, and validates that the required parameters are set
func (s *DHCPServer) DefaultAndValidate() error {
	if s.DNSServer == nil {
		s.DNSServer = net.IP{1, 1, 1, 1}
	}
	if s.SubnetMask == nil {
		s.SubnetMask = net.IP{255, 255, 255, 0}
	}
	if s.LeaseDuration == 0 {
		s.LeaseDuration = 7 * time.Day
	}
	if s.GatewayIP == nil || s.ClientIP == nil || s.ClientMAC == nil {
		return fmt.Errorf("DHCPServer requires GatewayIP, ClientIP and ClientMAC to be set")
	}
}

// ServeDHCP implements the dhcp.Handler interface
func (s *DHCPServer) ServeDHCP(p dhcp.Packet, request dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	// Answer DISCOVER and REQUEST calls only for this specific MAC address
	var response dhcp.MessageType
	switch request {
	case dhcp.Discover:
		response = dhcp.Offer
	case dhcp.Request:
		response = dhcp.ACK
	}
	if h.Debug {
		fmt.Printf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, request.String(), options, response.String())
	}
	if response != 0 {
		opts := dhcp.Options{
			dhcp.OptionSubnetMask:       []byte(s.SubnetMask),
			dhcp.OptionRouter:           []byte(s.GatewayIP),
			dhcp.OptionDomainNameServer: []byte(s.DNSServer),
		}
		optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
		requestingMAC := p.CHAddr().String()
		if requestingMAC == h.clientMAC.String() {
			if h.Debug {
				fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", response.String(), s.GatewayIP.String(), s.ClientIP.String(), optSlice, requestingMAC)
			}
			return dhcp.ReplyPacket(p, response, h.GatewayIP, h.ClientIP, h.LeaseDuration, optSlice)
		}
	}
	return nil
}
