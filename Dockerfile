FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    heapstate

WORKDIR $GOPATH/src/github.com/aureleoules/heapstate
COPY . .

RUN go mod download
RUN go mod verify

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-w -s" -o /go/bin/heapstate

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/heapstate /go/bin/heapstate

# USER heapstate:heapstate

ENTRYPOINT ["/go/bin/heapstate"]
EXPOSE 80 443