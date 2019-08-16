FROM golang:1.10.0


WORKDIR /go/src/github.com/sugimotosyo/hellogo
ADD . /go/src/github.com/sugimotosyo/hellogo
ENV GOPATH=/go


EXPOSE  8080


RUN go get -v github.com/golang/dep
RUN go install -v github.com/golang/dep/cmd/dep


RUN dep ensure

ENV export -p 




CMD ["go", "run", "main.go"]