FROM golang:1.22-alpine as dev

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY ./ ./
RUN go mod download
ENV GOPATH= 



# RUN go build -o /bench

FROM alpine:latest as production
COPY --from=dev /app/docs ./docs
COPY --from=dev /bench ./bench
WORKDIR /app
ENV PORT=8080
ENTRYPOINT [ "/bench" ]

