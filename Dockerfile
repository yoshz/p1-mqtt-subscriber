FROM golang AS BUILD
WORKDIR /src
ADD go.mod go.sum main.go /src/
RUN go get
RUN go build -ldflags "-linkmode external -extldflags -static" -o p1-mqtt-subscriber

FROM scratch
COPY --from=build /src/p1-mqtt-subscriber /
CMD ["/p1-mqtt-subscriber"]
