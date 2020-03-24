FROM golang:latest

WORKDIR /app
COPY . .

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o heapstate" --command=./heapstate
EXPOSE 9000