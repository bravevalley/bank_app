# Build Stage - Dev
FROM  golang:1.20.4-alpine3.17 AS Buidah
WORKDIR /app

COPY . .

RUN go build -o server main.go


# Final Stage - Prod
FROM alpine:latest
WORKDIR /app
COPY --from=Buidah /app/server .
COPY --from=Buidah /app/app.env .


EXPOSE 8080
CMD ["/app/server"]