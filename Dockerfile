# Use the official Golang image as the build base
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

# Set the Go working directory
WORKDIR /app

# Copy go.mod and go.sum files for dependency management optimization
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Compile the Go program
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /app/itv .

# Runtime image
FROM alpine:latest

# Adding the utilities tools to create the container
RUN apk --no-cache --update add tzdata

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

# Set the working directory in the container
WORKDIR /app

# Copy the contents
COPY --from=builder /app/itv /app/itv

# Expose port 8123
EXPOSE 8123

# Add entrypoint script
ADD docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

# Define entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]