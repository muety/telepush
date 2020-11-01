FROM golang:alpine

WORKDIR /data
COPY . .
RUN GO111MODULE=on go build -o output .

FROM alpine

WORKDIR /data
COPY --from=0 /data/output /webhook2telegram

EXPOSE 8080
ENTRYPOINT ["/webhook2telegram"]
