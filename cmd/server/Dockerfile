FROM golang:1.19.2-bullseye as stage1
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM denoland/deno:bin-1.27.0 as stage2

FROM debian:bullseye-20220125-slim
ENV RUNNER_SCRIPT=/app/runner.ts
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic-dev \
    libmagic-mgc \
    ca-certificates \
    mime-support \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY --from=stage1 /src/authgear-deno /usr/local/bin/
COPY --from=stage1 /src/pkg/deno/runner.ts /app/
COPY --from=stage2 /deno /usr/local/bin/
EXPOSE 8090
CMD ["authgear-deno"]
