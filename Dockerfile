# Build Stage - Dev
FROM  golang:1.20.4-alpine3.17 AS Buidah
WORKDIR /app

COPY . .

RUN go build -o server main.go
RUN $ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz  | tar xvz


# Final Stage - Prod
FROM alpine:latest
WORKDIR /app
COPY --from=Buidah /app/server .
COPY --from=Buidah /app/app.env .
COPY --from=Buidah /app/migrate /usr/bin/migrate
COPY /db/migration ./migrations
COPY ./entry.sh .



EXPOSE 8080
CMD ["/app/server"]
ENTRYPOINT ["./entry.sh"]