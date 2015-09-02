FROM alpine:3.1

COPY . /go/yml2vcard

ENV GOPATH=/go/yml2vcard/Godeps/_workspace

RUN apk --update add go && \
    go build -o /bin/yml2vcard /go/yml2vcard/yml2vcard.go && \
    rm -rf /go && \
    apk del go

ENTRYPOINT ["/bin/yml2vcard"]
