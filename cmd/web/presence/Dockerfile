# Build
FROM golang:1.24.2-alpine AS builder
WORKDIR /build
COPY . .
RUN go -C cmd/web/presence build -o main 

# Application
FROM alpine:latest
RUN addgroup -S tarantula && adduser -S tarantula -G tarantula

#USER root
RUN mkdir /etc/tarantula 
COPY cmd/web/presence/presence-conf.json /etc/tarantula/presence-conf.json

USER tarantula
#VOLUME /home/tarantula/presence
WORKDIR /home/tarantula/cmd/web/presence


COPY --from=builder /build/cmd/web/presence/main .

EXPOSE 8080
EXPOSE 5000

CMD ["./main"]