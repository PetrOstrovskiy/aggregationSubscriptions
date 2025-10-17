FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o aggregationSubscriprions

EXPOSE 8080

CMD ["sh", "-c", "sleep 5 && ./aggregationSubscriprions"]