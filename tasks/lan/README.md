## Как перезагружать роутер

Если участники удалят форвардинг интерфейса роутера наружу, его можно восстановить командой:

```
./bring-router-back
```

При этом нужно, чтобы устройство было подключено к VPN.


## Как подключаться к VNC с Zoom

```
ssh -o ProxyCommand='ssh incubator.ttc.tf "docker exec -i \`docker container ls -f name=marina-laptop -q\` sshd -i"' -i laptop/ssh/id_ed25519 -L 5902:localhost:5902 work@Marina-PC-XFCE410Pro
```

После этого можно подключаться на `vnc://localhost:5902` с паролем `9m4oxjdU`.

VNC, в который изначально попадают участники -- аналогично на порту `5901` с паролем `m4sha819`.
