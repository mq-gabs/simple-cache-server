FROM golang:1.23.3 AS build
WORKDIR /app
COPY . .
RUN go build -o server


FROM build
COPY /app/server /server
EXPOSE 9012
CMD ["./server"]