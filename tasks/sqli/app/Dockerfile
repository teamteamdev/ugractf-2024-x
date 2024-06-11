FROM alpine
RUN apk update && apk add python3 py3-pip && pip install quart --break-system-packages
WORKDIR /app
CMD ["/usr/bin/python3", "main.py"]