### Back-End
FROM --platform=${BUILDPLATFORM} golang:alpine AS build-back
ENV CGO_ENABLED=0
WORKDIR /app
ARG COMMIT="N/A"
RUN --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=cache,target=/go/pkg \
    go mod download
ARG TARGETARCH TARGETOS
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build -o /out/wg-gen-web -ldflags "-w -s -X 'github.com/vx3r/wg-gen-web/version.Version=${COMMIT::7}'" ./cmd/wg-gen-web

### Front-End
FROM --platform=${BUILDPLATFORM} node:18-alpine AS build-front
WORKDIR /app
COPY ui/package.json ui/package-lock.json ./
RUN npm ci --no-fund
COPY ui/ ./
RUN npm run build

FROM alpine
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
COPY .env .
COPY --from=build-back /out/wg-gen-web .
COPY --from=build-front /app/dist ./ui/dist
RUN chmod +x ./wg-gen-web
EXPOSE 8080

CMD ["/app/wg-gen-web"]
