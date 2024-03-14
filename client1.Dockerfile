FROM golang:1.22 as builder 

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o bin ./cmd 

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin . 

RUN echo ls

CMD ["./bin"] 

EXPOSE 8080