FROM alpine AS base
RUN apk update && apk add --no-cache curl iproute2
COPY start.sh /start.sh
CMD ["/start.sh", "sleep", "inf"]

FROM ubuntu AS marina
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y curl dbus-x11 iproute2 openssh-server systemd tightvncserver xfce4 xfce4-goodies
RUN curl -sSL https://zoom.us/client/latest/zoom_amd64.deb -o /tmp/zoom_setup.deb && \
    apt-get install -y /tmp/zoom_setup.deb
RUN echo "PasswordAuthentication no" >>/etc/ssh/sshd_config && \
    echo "KbdInteractiveAuthentication no" >>/etc/ssh/sshd_config
RUN useradd --create-home --shell /bin/bash --password "$(perl -e "print crypt('m4sha819', 'sa');")" --user-group marina && \
    mkdir /home/marina/.vnc && \
    echo m4sha819 | vncpasswd -f >/home/marina/.vnc/passwd && \
    chown -R marina:marina /home/marina/.vnc && \
    chmod -R go= /home/marina/.vnc
RUN useradd --create-home --shell /bin/bash --password "*" --user-group work && \
    mkdir /home/work/.vnc && \
    echo 9m4oxjdU | vncpasswd -f >/home/work/.vnc/passwd && \
    chown -R work:work /home/work/.vnc && \
    chmod -R go= /home/work/.vnc
COPY --chown=marina:marina xstartup /home/marina/.vnc/xstartup
COPY --chown=work:work xstartup /home/work/.vnc/xstartup
COPY --chown=work:work ssh/id_ed25519.pub /home/work/.ssh/authorized_keys
RUN chmod -R go= /home/work/.ssh
COPY vnc.service vnc-work.service /etc/systemd/system/
RUN systemctl enable ssh vnc vnc-work
COPY start.sh /start.sh
CMD ["/start.sh", "/bin/systemd"]
