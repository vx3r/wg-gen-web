ARG COMMIT="N/A"

FROM golang:alpine AS build-back
WORKDIR /app
ARG COMMIT
COPY . .
RUN go build -o wg-gen-web-linux -ldflags="-X 'gitlab.127-0-0-1.fr/vx3r/wg-gen-web/version.Version=${COMMIT}'" gitlab.127-0-0-1.fr/vx3r/wg-gen-web/cmd/wg-gen-web

FROM node:10-alpine AS build-front
WORKDIR /app
COPY ui/package*.json ./
RUN npm install
COPY ui/ ./
RUN npm run build

FROM alpine
WORKDIR /app
COPY --from=build-back /app/wg-gen-web-linux .
COPY --from=build-front /app/dist ./ui/dist
COPY .env .
RUN chmod +x ./wg-gen-web-linux
RUN apk add --no-cache ca-certificates
EXPOSE 8080

CMD ["/app/wg-gen-web-linux"]