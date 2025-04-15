FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(echo $TARGETPLATFORM | cut -d/ -f2) go build -o hello-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/hello-service .
EXPOSE 3000
CMD ["./hello-service"]
