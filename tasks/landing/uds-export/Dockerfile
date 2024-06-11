FROM alpine
RUN apk update && apk add socat
CMD ["/usr/bin/socat", "UNIX-LISTEN:/tmp/app.sock,mode=0666,fork", "TCP:web:3000"]
