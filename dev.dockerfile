ARG COMMIT="N/A"

FROM golang AS build-back
WORKDIR /app
ARG COMMIT
COPY . .
RUN go build -o wg-gen-web-linux -gcflags="all=-N -l" -ldflags="-X 'github.com/vx3r/wg-gen-web/version.Version=${COMMIT}'" github.com/vx3r/wg-gen-web/cmd/wg-gen-web
RUN go get github.com/go-delve/delve/cmd/dlv

FROM node:lts AS build-front
WORKDIR /app
COPY ui/package*.json ./
RUN npm install
COPY ui/ ./
RUN npm run build

FROM debian
WORKDIR /app
COPY --from=build-back /app/wg-gen-web-linux .
COPY --from=build-back /go/bin/dlv .
COPY --from=build-front /app/dist ./ui/dist
COPY .env .
RUN chmod +x ./wg-gen-web-linux
RUN apt-get update && apt-get install -y ca-certificates
EXPOSE 8080

CMD ["/app/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/wg-gen-web-linux"]