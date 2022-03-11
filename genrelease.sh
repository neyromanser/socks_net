#!/bin/sh
# install https://github.com/mitchellh/gox

mkdir -p release/

if [ "$1" = "dev" ]; then
  go build -o release/ ./client/.
  go build -o release/ ./server/.
else
  #export GOX=$HOME/go/bin/gox
  export GOX=gox
  export CGO_ENABLED=0
  export GOX_LINUX_AMD64_LDFLAGS="-extldflags -static -s -w"
  export GOX_LINUX_386_LDFLAGS="-extldflags -static -s -w"
  export G_OS="linux windows"
  export G_ARCH="386 amd64 arm arm64"

  $GOX -os="$G_OS" -arch="$G_ARCH" -ldflags="-s -w" -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" ./client 2>release/error.log
  $GOX -os="$G_OS" -arch="$G_ARCH" -ldflags="-s -w" -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" ./server 2>>release/error.log
fi