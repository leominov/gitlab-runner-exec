FROM golang:1.12.9-alpine3.10 as builder
WORKDIR /go/src/github.com/leominov/gitlab-runner-exec
COPY . .
ENV GO111MODULE=on
RUN apk --no-cache add git make && make build

FROM alpine:3.10
COPY --from=builder /go/src/github.com/leominov/gitlab-runner-exec/gitlab-runner-exec /go/bin/gitlab-runner-exec
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN apk --no-cache add git
ENTRYPOINT ["/go/bin/gitlab-runner-exec"]
