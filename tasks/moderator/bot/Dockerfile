FROM golang:bookworm

RUN DEBIAN_FRONTEND=noninteractive && \
  apt-get update && \
  apt-get -y install --no-install-recommends libsqlite3-dev sqlite3 && \
  apt-get clean && rm -rf /var/lib/apt/lists/*

COPY . /app
WORKDIR /app
RUN go build -o adminbot ./cmd/adminbot

EXPOSE 8080

CMD [ \
  "/app/adminbot", \
  "-config", "/app/config.json", \
  "-min-score", "1337", \
  "-period-voting", "15m", \
  "-min-candidates", "5", \
  "-min-activity", "100", \
  "-min-complexity", "3", \
  "-max-complexity", "10", \
  "-voting-duration", "20s" \
  ]
