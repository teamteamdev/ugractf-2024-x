FROM alpine
RUN apk update && \
    apk add --no-cache openssh socat && \
    adduser -D teaching-materials-key && \
    echo teaching-materials-key:this-is-a-teemoorka-test-assignment-for-job-candidates | chpasswd && \
    ssh-keygen -A && \
    sed -i "s/#Port 22/Port 5890/" /etc/ssh/sshd_config && \
    rm -rf /var/cache/apk/*
COPY secret /secret
COPY motd /etc/motd
COPY su /usr/bin/su
COPY note.txt /home/teaching-materials-key/note.txt
RUN touch "/home/teaching-materials-key/Что вы здесь забыли?"
CMD ["socat", "unix-listen:/tmp/app.sock,mode=0666,fork", "exec:/usr/sbin/sshd -i"]
# CMD ["/usr/sbin/sshd", "-D"]
