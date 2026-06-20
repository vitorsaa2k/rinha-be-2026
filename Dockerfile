FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY models/ models/
COPY public/ public/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o /out/build-index ./cmd/build-index

RUN mkdir -p /out/resources && \
    /out/build-index /src/public/references.json.gz /out/resources/index.bin

FROM golang:1.26-alpine
WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/resources/index.bin /app/resources/index.bin
COPY --from=builder /src/public/ /app/public/

EXPOSE 9999
CMD ["./api"]
