# ---- Build stage ---
FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'bisecur/version.Version=${{  github.ref_name }}' -X 'bisecur/version.BuildDate=$(date)'" -o /halsecur

FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /halsecur /halsecur
RUN chmod +x /halsecur

ENTRYPOINT ["/halsecur"]

# Build the Docker image with:
# docker build -t halsecur:latest .

# Run the Docker container with:
# docker run --rm -it -v /path/to/your/config:/config halsecur:latest --config /config/config.yaml