FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /drift ./cmd/drift

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /drift /usr/local/bin/drift
RUN mkdir -p /data
EXPOSE 8080
ENTRYPOINT ["drift"]
CMD ["serve", "--state", "/data/terraform.tfstate", "--port", "8080"]
