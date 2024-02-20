FROM node:18-bookworm AS ui-builder
WORKDIR /media-vault-spa
COPY website/package.json .
RUN yarn install
COPY website/index.html website/tsconfig.* website/vite.config.ts ./
COPY website/src/ src/
COPY website/public/ public/
RUN cd src/js/serviceWorker && npx webpack --mode production --config webpack.config.cjs
RUN yarn build-only

FROM golang:1.22.0-bookworm AS backend-builder
WORKDIR /go/src/github.com/rishabhkailey/media-vault
COPY go.* ./
RUN go mod download
COPY cmd cmd/
COPY internal internal/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o media-vault cmd/media-vault/main.go
RUN chmod +x media-vault

FROM debian:bookworm-20240211
RUN apt update \
  && apt install -y ca-certificates \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /etc/media-vault
COPY --from=backend-builder /go/src/github.com/rishabhkailey/media-vault/media-vault .
COPY --from=ui-builder /media-vault-spa/dist/ website/dist/
ENTRYPOINT [ "./media-vault" ]