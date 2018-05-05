FROM library/golang:alpine

RUN apk add --no-cache git mercurial

RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp

COPY elcep/*.go ./
RUN go build -o logmonitor

FROM alpine

RUN mkdir -p /go/src/github.com
COPY --from=0 /go/logmonitor /
COPY --from=0 /go/src/github.com /go/src/github.com
COPY /conf/queries.cfg /conf/queries.cfg

ENTRYPOINT ["./logmonitor"]
CMD [""]
