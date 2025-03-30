FROM golang:1.23.3

WORKDIR /src

COPY . .
RUN go build -o /bin/server

EXPOSE 9012

CMD ["/bin/server"]
