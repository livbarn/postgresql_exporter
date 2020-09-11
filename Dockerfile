FROM golang:latest as builder
WORKDIR /app
COPY . .
# RUN git clone https://github.com/livbarn/postgresql_exporter.git
RUN go install && GOOS=linux GOARCH=amd64 go build -o postgre_exporter .


FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /app/. /app
# ENV DATA_SOURCE_NAME='host=192.168.224.12 user=dbuser dbname=chat_local password=123kkk sslmode=disable'

ENTRYPOINT [ "/app/postgre_exporter" ]
