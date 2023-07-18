FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY .env ./

RUN go build -o /docker-go-discord

CMD ["/docker-go-discord"]
