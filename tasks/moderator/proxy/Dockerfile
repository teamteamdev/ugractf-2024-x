FROM alpine:3

RUN apk add --no-cache socat

EXPOSE 8080

CMD ["socat", "-d", "UNIX-LISTEN:/tmp/moderator.sock,fork,umask=0000,unlink-early", "TCP:adminbot:8080"]