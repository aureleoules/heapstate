FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/aureleoules/heapstate
COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/heapstate


FROM scratch

COPY --from=builder /go/bin/heapstate /go/bin/heapstate

CMD ["./go/bin/heapstate"]

