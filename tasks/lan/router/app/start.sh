#!/bin/sh
ip link add br0 type bridge
ip address add 10.13.0.1/20 dev br0
ip link set br0 up

iptables -t nat -A POSTROUTING -o tun0 -j MASQUERADE --random
iptables -A FORWARD -i br0 -o tun0 -j ACCEPT
iptables -A FORWARD -i tun0 -o br0 -m conntrack --ctstate RELATED,ESTABLISHED,DNAT -j ACCEPT
iptables -A FORWARD -j DROP

iptables -t nat -A PREROUTING -i tun0 -p tcp --dport 80 -j DNAT --to-destination 10.13.0.1:80

python3 emerg.py &

exec "$@"
