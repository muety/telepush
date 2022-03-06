FROM golang:alpine

WORKDIR /data
COPY . .
RUN go build -o telepush .

FROM alpine

WORKDIR /app
COPY --from=0 /data/telepush telepush
COPY --from=0 /data/views views
COPY --from=0 /data/version.txt version.txt
COPY --from=0 /data/docker/entrypoint.sh entrypoint.sh

ENV APP_MODE "webhook"

VOLUME /srv/data
EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
