import os
import socket


sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind(("", 64931))

while True:
    message = sock.recv(1024)
    if message == b"letmeiiiiin":
        print("Bringing configuration port forwarding back", flush=True)
        os.system("iptables -t nat -F PREROUTING")
        os.system("iptables -t nat -A PREROUTING -i tun0 -p tcp -m tcp --dport 80 -j DNAT --to-destination 10.13.0.1:80")
        os.system("wget http://10.13.0.1/reset-tables")
