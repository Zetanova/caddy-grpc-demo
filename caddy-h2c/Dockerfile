FROM golang:1.14-alpine3.11 AS builder

ENV CADDY_COMMIT b61e1a200635bac21b9c9992266ebe66b19910be

WORKDIR /go/src

RUN apk add --no-cache \
    git \
    ca-certificates

ENV CADDY_SOURCE_VERSION=h2c

RUN git clone -b $CADDY_SOURCE_VERSION https://github.com/caddyserver/caddy.git --single-branch

WORKDIR /go/src/caddy/cmd/caddy

RUN go build

#RUN go get -d ./...
#COPY caddy-builder.sh /usr/bin/caddy-builder
#WORKDIR /src/custom-caddy/cmd/caddy

RUN mv ./caddy /go/bin/



FROM alpine:3.11

RUN apk add --no-cache ca-certificates

# https://github.com/caddyserver/dist/commits
ENV CADDY_DIST_COMMIT 80870b227ded910971ecace4a0c136bf0ef46342

RUN set -eux; \
	mkdir -p \
		/config/caddy \
		/data/caddy \
		/etc/caddy \
		/usr/share/caddy \
	; \
	wget -O /etc/caddy/Caddyfile "https://github.com/caddyserver/dist/raw/$CADDY_DIST_COMMIT/config/Caddyfile"; \
	wget -O /usr/share/caddy/index.html "https://github.com/caddyserver/dist/raw/$CADDY_DIST_COMMIT/welcome/index.html"

COPY --from=builder /go/bin/* /usr/bin/

RUN chmod +x /usr/bin/caddy; \
	caddy version

# See https://caddyserver.com/docs/conventions#file-locations for details
ENV XDG_CONFIG_HOME=/config
ENV XDG_DATA_HOME=/data

COPY ./caddy.json /etc/caddy/caddy.json

VOLUME /config
VOLUME /data

EXPOSE 80
EXPOSE 443
EXPOSE 2019

CMD ["caddy", "run", "--config", "/etc/caddy/caddy.json"]