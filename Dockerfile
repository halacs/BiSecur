FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ADD dist/halsecur /halsecur
RUN chmod +x /halsecur

ENTRYPOINT ["/halsecur"]

# Build the Docker image with:
# docker build -t halsecur:latest .

# Run the Docker container with:
# docker run --rm -it -v /path/to/your/config:/config halsecur:latest --config /config/config.yaml