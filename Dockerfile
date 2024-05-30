FROM golang:1.22-alpine

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/gopls@latest
COPY ./ ./
RUN go mod download

ENTRYPOINT [ "go","test","-bench",".","--benchmem", "-cpu=16"]