FROM luxas/firectl
RUN apk add --update iproute2 dhcp
COPY dhcpd.conf /etc/dhcpd.conf.tmpl
RUN mkdir -p /var/lib/dhcp/ /var/run/dhcp && touch /var/lib/dhcp/dhcpd.leases
COPY run.sh fc-dhcpd /
ENTRYPOINT /run.sh
