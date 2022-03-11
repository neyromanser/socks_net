#!/bin/sh
# apk add go-cross

mkdir -p release/

if $1 == "dev"; then
  #export GOX=$HOME/go/bin/gox
  export GOX=gox
  export CGO_ENABLED=0
  export GOX_LINUX_AMD64_LDFLAGS="-extldflags -static -s -w"
  export GOX_LINUX_386_LDFLAGS="-extldflags -static -s -w"

  $GOX -os="linux windows" -ldflags="-s -w" -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" ./client/. 2>release/error.log
else
  go build -o release/client_dev ./client/.
fi