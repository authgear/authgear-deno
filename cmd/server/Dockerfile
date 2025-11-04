FROM quay.io/theauthgear/golang:1.23.6-noble AS stage1
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM denoland/deno:bin-1.41.3 AS stage2

FROM ubuntu:noble
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic-dev \
    libmagic-mgc \
    ca-certificates \
    media-types \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=stage1 /src/authgear-deno /usr/local/bin/
COPY --from=stage2 /deno /usr/local/bin/
COPY ./docker-entrypoint.sh /usr/local/bin/
# Ask Deno to use the system certificate store instead.
# The system certificate store supports mounting extra certificates to /usr/local/share/ca-certificates/*.crt
# https://deno.com/blog/v1.13#use-system-certificate-store-for-tls
ENV DENO_TLS_CA_STORE=system
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 8090
CMD ["authgear-deno"]
