FROM golang:alpine AS build-back
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache git
COPY . .
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && go build -ldflags "-X main.VersionGitCommit=$GIT_COMMIT" go build -o wg-gen-web-linux

FROM node:10-alpine AS build-front
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache git
COPY ui/package*.json ./
RUN npm install
COPY ui/ ./
COPY .git .
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