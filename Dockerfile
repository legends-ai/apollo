FROM golang:1.7.1-wheezy

# Copy files. we don't need glide etc with this
COPY . /go/src/github.com/asunaio/apollo
WORKDIR /go/src/github.com/asunaio/apollo

# Build binary
RUN rm -f apollo
RUN go build .
