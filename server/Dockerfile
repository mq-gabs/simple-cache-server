FROM golang:1.23.3 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server

FROM scratch
COPY --from=build /bin/server /
EXPOSE 9012
CMD [ "/server" ]