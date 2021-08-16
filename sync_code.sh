#!/bin/sh
# rsync code from mac to linux.

work_path=$(dirname $0)
cd $work_path

src=~/go/src/agent
dst=root@wenfh2020_sgx.com:/home/gopath/src
echo "$src --> $dst"

# only rsync *.go files.
rsync -ravz --exclude=".git/" --include="*.go" --include="*/" --exclude="*" $src $dst
