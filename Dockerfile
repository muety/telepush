FROM golang:alpine

WORKDIR /data
COPY . .
RUN GO111MODULE=on go build -o output .

FROM alpine

WORKDIR /data
COPY --from=0 /data/output /telegram-middleman-bot

EXPOSE 8080
ENTRYPOINT ["/telegram-middleman-bot"]
