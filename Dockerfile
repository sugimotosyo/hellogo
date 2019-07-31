FROM golang:1.10.0

# RUN go get github.com/labstack/echo/...

WORKDIR /go
ADD . /go

EXPOSE  8080

CMD ["go", "run", "main.go"]