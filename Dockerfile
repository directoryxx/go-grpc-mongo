FROM golang:1.16.7

# Install VSCode Dependencies Suggestion
RUN go get github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get github.com/ramya-rao-a/go-outline
RUN go get github.com/cweill/gotests/gotests
RUN go get github.com/fatih/gomodifytags
RUN go get github.com/josharian/impl
RUN go get github.com/haya14busa/goplay/cmd/goplay
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/go-delve/delve/cmd/dlv@master
RUN go get honnef.co/go/tools/cmd/staticcheck
RUN go get golang.org/x/tools/gopls

RUN mkdir /app
WORKDIR /app

COPY . /app

#RUN apt update & apt install -y clang-format protobuf-compiler

RUN go get -u google.golang.org/grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
RUN go get github.com/golang/protobuf/protoc-gen-go

RUN cd /tmp && wget https://github.com/ktr0731/evans/releases/download/0.10.0/evans_linux_arm.tar.gz && tar -xvf evans_linux_arm.tar.gz && mv evans /usr/local/bin

# RUN go mod download
CMD /bin/sh
