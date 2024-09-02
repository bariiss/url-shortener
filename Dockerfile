FROM ghcr.io/bariiss/golang-upx:1.23.0-bookworm AS build

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

# Create a minimal image
FROM scratch
LABEL org.opencontainers.image.source="https://github.com/bariiss/url-shortener"

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
