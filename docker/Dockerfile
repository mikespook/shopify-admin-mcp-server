FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /shopify-admin-mcp-server ./cmd/shopify-admin-mcp-server


FROM alpine:3

RUN adduser -D -u 1000 appuser

COPY --from=builder /shopify-admin-mcp-server /usr/local/bin/shopify-admin-mcp-server

USER appuser
WORKDIR /home/appuser

EXPOSE 9093

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
	CMD wget -qO- http://localhost:9093/health || exit 1

ENTRYPOINT ["shopify-admin-mcp-server", "-addr", "0.0.0.0:9093"]
