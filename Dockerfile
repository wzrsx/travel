FROM golang:1.22

WORKDIR /webApp

COPY . .

RUN go build -o app_route ./cmd/app_route

CMD ["./app_route"]