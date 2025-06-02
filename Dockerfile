FROM golang:1.24
WORKDIR /usr/local/app

COPY go.mod go.sum ./
RUN go mod download
RUN apt update && apt install -y pass

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
