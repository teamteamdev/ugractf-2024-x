```
$ ./new-client <id>
```
ID будет виден участникам, так что туда нужно вставить, например, имя пользователя в Telegram.

```
$ ./get-client-ovpn <id>
```
Выдает уже готовую конфигурацию `<id>.ovpn`, включающую в себя адрес сервера. Это только для внутренних сервисов, участникам не давать.

```
$ ./get-client-crt <id>
$ ./get-client-key <id>
```
Выдает соответственно файлы `<id>.crt` и `<id>.key`, они не содержат адрес сервера.

Чтобы подключиться к серверу, нужно иметь либо конфиг `ovpn`, либо комбинацию `crt`+`key`+`static-key`+адрес сервера, где `crt` и `key` отдельные на пользователя, а `static-key` можно получить через:
```
$ ./get-static-key
```
Участник его получает через wiki, как и адрес.

Команда для подключения с этими данными выглядит так:
```
$ openvpn \
	--client \
	--dev tun \
	--remote {адрес сервера} \
	--cert {path-to-crt} \
	--key {path-to-key} \
	--tls-auth {path-to-static-key} \
	--peer-fingerprint 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00
```
Это выведет ошибку:
```
...
2024-06-16 15:02:12 TLS Error: --tls-verify/--peer-fingerprintcertificate hash verification failed. (got fingerprint: 1f:...:bf)
```
Нормальный fingerprint после этого нужно подставить в команду.
