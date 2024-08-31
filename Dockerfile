# ----------------- First stage -----------------
# Go version: go1.23
FROM golang:1.23-bookworm AS build

# Set the upx version
ARG upx_version=4.2.4

# Install upx and cleanup
RUN apt-get update && apt-get install -y --no-install-recommends xz-utils && \
  arch=$(uname -m) && \
  case "$arch" in \
    x86_64) TARGETARCH=amd64 ;; \
    aarch64) TARGETARCH=arm64 ;; \
    *) echo "Unsupported architecture: $arch" && exit 1 ;; \
  esac && \
  curl -Ls https://github.com/upx/upx/releases/download/v${upx_version}/upx-${upx_version}-${TARGETARCH}_linux.tar.xz -o - | tar xvJf - -C /tmp && \
  cp /tmp/upx-${upx_version}-${TARGETARCH}_linux/upx /usr/local/bin/ && \
  chmod +x /usr/local/bin/upx && \
  apt-get remove -y xz-utils && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=linux go build -o url-shortener -a -ldflags="-s -w" -installsuffix cgo

# Compress the binary
RUN upx --ultra-brute -qq url-shortener && upx -t url-shortener
# ----------------- End of the first stage -----------------

# ----------------- Second stage -----------------
# Create a minimal image
FROM scratch

WORKDIR /app

# Copy the binary from the build container
COPY --from=build /app/url-shortener /app/url-shortener

# Copy the static files
COPY --from=build /app/static /app/static

# Copy the templates
COPY --from=build /app/templates /app/templates

# Copy environment file
COPY .env /app/.env

# Copy the configuration file
ENTRYPOINT ["/app/url-shortener"]
# ----------------- End of the second stage -----------------
