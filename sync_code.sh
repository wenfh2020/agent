#!/bin/sh
# copy src from mac to linux.
SRC=~/go/src/agent
LAN_DST=root@192.168.0.200:/home/gopath/src
REMOTE_DST=root@wenfh2020_sgx.com:/home/gopath/src
DST=$LAN_DST

work_path=$(dirname $0)
cd $work_path

[ $1x == 'remote'x ] && DST=$REMOTE_DST

echo "send files to: $DST"

rsync -avz --exclude="*.o" \
    --exclude=".git" \
    --exclude=".vscode" \
    --exclude="*.so" \
    --exclude="*.a" \
    --exclude="*.log" \
    --exclude="config.yml" \
    --exclude="sync_code.sh" \
    --exclude="main" \
    --exclude="client" \
    $SRC $DST
