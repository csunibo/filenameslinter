FROM golang:alpine AS go-builder
COPY . /build
WORKDIR /build
RUN go build -ldflags "-s -w" -o /build/filenameslinter

FROM alpine
COPY --from=go-builder /build/filenameslinter /usr/bin/filenameslinter
RUN chmod +x /usr/bin/filenameslinter

ENTRYPOINT ["/usr/bin/filenameslinter"]
