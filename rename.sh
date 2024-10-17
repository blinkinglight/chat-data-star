#!/usr/bin/env sh

export CUR="github.com/ituoga/go-start" 
export NEW="$1"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;
