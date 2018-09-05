FROM golang:1.11 as builder

WORKDIR /go/src/github.com/uphy/elastalert-mail-gateway
#COPY go.mod go.sum ./
#RUN go mod download
COPY vendor ./vendor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /elastalert-mail-gateway

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /elastalert-mail-gateway /bin/
WORKDIR /etc/elastalert-mail-gateway
COPY example/declarative.yml config.yml
ENTRYPOINT [ "/bin/elastalert-mail-gateway", "/etc/elastalert-mail-gateway/config.yml" ]