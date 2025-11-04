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
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 8090
CMD ["authgear-deno"]
