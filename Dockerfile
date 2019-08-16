FROM golang:1.10.0


WORKDIR /go
ADD . /go

EXPOSE  8080

# go get -v github.com/golang/dep
# go install -v github.com/golang/dep/cmd/dep
# go get -u github.com/golang/dep/cmd/dep これどこかに必要？ これ微妙



RUN go get -v github.com/golang/dep
RUN go install -v github.com/golang/dep/cmd/dep


CMD ["go", "run", "main.go"]