FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ВАЖНО: путь без дефиса, соответствует cmd/latexservice/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/latexservice ./cmd/latexservice

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    texlive-latex-base latexmk ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/latexservice /app/latexservice
COPY templates ./templates

EXPOSE 8081

ENTRYPOINT ["/app/latexservice"]
