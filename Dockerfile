FROM 		golang:1.12-alpine as builder
WORKDIR		/go
ADD			app.go .
RUN			go build -o app


FROM 		alpine:latest
WORKDIR 	/go/
COPY 		--from=builder /go/app .
ENTRYPOINT  /go/app
