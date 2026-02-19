# Standard Dockerfile for local development and manual builds
# For automated releases, see Dockerfile.goreleaser (used by goreleaser)
FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o cisshgo cissh.go

FROM scratch
COPY --from=builder /app/cisshgo /cisshgo
COPY --from=builder /app/transcripts /transcripts
ENTRYPOINT ["/cisshgo"]