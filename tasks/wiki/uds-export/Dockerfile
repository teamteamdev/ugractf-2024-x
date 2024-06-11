FROM alpine
RUN apk update && apk add socat
CMD ["/usr/bin/socat", "UNIX-LISTEN:/tmp/app.sock,mode=0666,fork", "TCP:mediawiki:80"]
