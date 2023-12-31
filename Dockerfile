FROM golang:alpine AS go-builder
COPY . /build
WORKDIR /build/cmd
RUN go build -ldflags "-s -w" -o /build/filenameslinter

FROM alpine
COPY --from=go-builder /build/filenameslinter /usr/bin/filenameslinter

ENTRYPOINT ["/usr/bin/filenameslinter"]
