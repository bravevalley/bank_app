FROM  golang:1.20.4-alpine3.17
WORKDIR /app

COPY . .

RUN go build -o server main.go

EXPOSE 8080
CMD ["/app/server"]