# Builder
FROM whatwewant/builder-go:v1.20-1 as builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build \
  -trimpath \
  -ldflags '-w -s -buildid=' \
  -v -o gzproxy

# Server
FROM whatwewant/go:v1.20-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-zoox/gzproxy"

ARG VERSION=latest

ENV MODE=production

COPY --from=builder /build/gzproxy /bin

ENV VERSION=${VERSION}

COPY ./entrypoint.sh /

CMD /entrypoint.sh
