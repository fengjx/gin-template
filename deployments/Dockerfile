FROM node:24-alpine AS admin-builder
WORKDIR /app/admin
COPY admin/package.json admin/package-lock.json ./
RUN npm ci
COPY admin ./
RUN npm run build

FROM golang:1.25-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=admin-builder /app/admin/dist ./admin/dist
RUN CGO_ENABLED=1 go build -o /gin-template ./cmd/server

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates sqlite
COPY --from=go-builder /gin-template /app/gin-template
COPY configs /app/configs
EXPOSE 3000
CMD ["/app/gin-template", "serve", "--env", "prod", "--config", "/app/configs/config.prod.yaml"]
