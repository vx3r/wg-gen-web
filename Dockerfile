FROM alpine

WORKDIR /app

ADD wg-gen-web-linux-amd64 .
ADD .env .
ADD ui/dist ui/dist

RUN chmod +x ./wg-gen-web-linux-amd64
RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["/app/wg-gen-web-linux-amd64"]
