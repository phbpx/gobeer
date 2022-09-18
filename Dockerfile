# Build the Go Binary.
FROM golang:1.19 as builder
ENV CGO_ENABLED 0
ARG BUILD_REF

# Copy the source code into the container.
COPY . /service

# Build the service binary.
WORKDIR /service/cmd/gobeer-api
RUN go build

# Run the Go Binary in Alpine.
FROM alpine:3.16
WORKDIR /app

# Create a non-root user.
RUN addgroup -g 1000 -S appuser && \
  adduser -u 1000 -h /service -G appuser -S appuser

COPY --from=builder /service/cmd/gobeer-api/gobeer-api /service/gobeer-api
WORKDIR /service

USER appuser
CMD ["./gobeer-api"]
