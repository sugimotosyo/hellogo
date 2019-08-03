FROM golang:1.10.0


WORKDIR /go
ADD . /go

EXPOSE  8080

CMD ["go", "run", "main.go"]