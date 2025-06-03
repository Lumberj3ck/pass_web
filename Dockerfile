FROM golang:1.24
WORKDIR /usr/local/app

COPY go.mod go.sum ./
RUN go mod download
RUN apt update && apt install -y pass

EXPOSE 8080
RUN go install github.com/air-verse/air@latest

CMD ["air", "run", "main.go"]
