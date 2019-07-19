FROM 		golang:1.12-alpine as builder

ENV     GO111MODULE=on

RUN     apk add git

WORKDIR	/kubernetes

#       Copies go.mod and go.sum if go.sum exist
COPY    go.* ./

RUN     go mod download

COPY		app.go .
COPY		job.go .
COPY		secret.go .

RUN			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app


FROM 		scratch
COPY 		--from=builder /kubernetes/app /app
ENTRYPOINT  ["/app"]
