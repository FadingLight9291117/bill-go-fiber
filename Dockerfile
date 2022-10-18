FROM golang:alpine

WORKDIR /usr/src/app

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change

COPY . .
COPY .env .

RUN go mod verify

RUN go build -v -o /usr/local/bin/app .

CMD ["app"]
