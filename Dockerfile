FROM golang:1.23 AS builder
WORKDIR /app
COPY server.go ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server ./server.go

FROM gcr.io/distroless/static-debian12 AS runner
WORKDIR /app
COPY --from=builder /app/server ./server
USER nonroot:nonroot
COPY --chown=nonroot:nonroot docs/ ./
ENTRYPOINT ["/app/server"]
CMD ["-p", "8080"]