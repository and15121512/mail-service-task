FROM golang:latest as builder

RUN mkdir -p /task
ADD . /task
WORKDIR /task

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /task ./cmd/main.go

FROM scratch
COPY --from=builder /task /task
COPY --from=builder /etc/ssl/certs /etc/ssl/certs/
WORKDIR /task

CMD ["./main"]
EXPOSE 3000 4000
