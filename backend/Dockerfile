FROM golang:1.21-alpine

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

COPY wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh

EXPOSE 8080

CMD ["./wait-for-postgres.sh", "./main"]