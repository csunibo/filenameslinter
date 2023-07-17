FROM golang:alpine AS go-builder
COPY . /build
WORKDIR /build
RUN go build -ldflags "-s -w" -o /build/filenameslinter

FROM scratch
COPY --from=go-builder /build/filenameslinter /usr/bin/filenameslinter

ENTRYPOINT ["/usr/bin/filenameslinter"]
