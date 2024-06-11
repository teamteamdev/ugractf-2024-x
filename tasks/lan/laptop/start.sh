#!/bin/sh
if [ "$1" = --internal ]; then
	shift

	hostname "$NAME"

	ip address add "$IP/20" dev wlan0
	ip link set wlan0 up
	ip route add default via 10.13.0.1

	while ! curl http://10.13.0.1/new-client -F "name=$NAME"; do
		sleep 1
	done

	rm /start.sh

	exec "$@"
	exit
fi

umount /etc/resolv.conf
umount /etc/hostname
echo "127.0.0.1 localhost" >/etc/hosts

printf "nameserver 77.88.8.8\nnameserver 77.88.8.1\n" >/etc/resolv.conf
echo "$NAME" >/etc/hostname

id="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 8)"

ip netns add "net-$id"

ip link add "veth-$id" type veth peer name wlan0 netns "net-$id"
while ! ip link set "veth-$id" master br0; do
	sleep 1
done

ip link set "veth-$id" up

exec ip netns exec "net-$id" unshare -u ./start.sh --internal "$@"
