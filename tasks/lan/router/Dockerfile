FROM python:alpine
RUN apk update && apk add --no-cache iptables && pip install aiohttp quart uvicorn
WORKDIR /app
CMD ["/app/start.sh", "uvicorn", "--host", "10.13.0.1", "--port", "80", "main:app"]
