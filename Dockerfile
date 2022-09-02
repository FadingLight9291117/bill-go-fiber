FROM golang:alpine

WORKDIR /usr/src/app

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
COPY go.mod  ./
RUN go mod download && go mod verify

COPY . .
COPY .env .
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]