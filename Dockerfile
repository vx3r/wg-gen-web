### Back-End
FROM golang:alpine AS go-base
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM go-base AS build-back
ARG COMMIT="N/A"
COPY . .
RUN go install -ldflags "-w -s -X 'github.com/vx3r/wg-gen-web/version.Version=${COMMIT::7}'" ./cmd/...

### Front-End
FROM node:18-alpine AS node-base
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
RUN npm ci --no-fund

FROM node-base AS build-front
COPY ui/ .
RUN npm run build

### Final
FROM alpine AS final-base
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
COPY .env .
COPY --from=build-back /go/bin/wg-gen-web .
COPY --from=build-front /app/dist ./ui/dist
RUN chmod +x ./wg-gen-web
EXPOSE 8080

CMD ["/app/wg-gen-web"]
