FROM golang:1.12 as build
WORKDIR /go/src/github.com/shawntoffel/ddns
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "nobody:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/shawntoffel/ddns/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=100:101 /go/src/github.com/shawntoffel/ddns/bin/ddns .
USER nobody
ENTRYPOINT ["./ddns"]
