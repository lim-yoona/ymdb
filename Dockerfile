FROM golang:1.21

WORKDIR /root/ymdb

RUN go env -w GOPROXY=https://goproxy.cn,direct
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ymdb ./main.go

CMD ["/usr/local/bin/ymdb"]