#!/bin/bash

parent=`dirname $PWD`
parent=`dirname $parent`

#go-bindata -o="../../static_data.go" -pkg="patrol" -prefix=$parent ../../static/...
rm -f patrol
go build -tags "EMBED_STATIC"
OUT=$?
if [ $OUT -eq 0 ];then
    ./patrol -secret_key=baGCbYmpdRxeSZ2rJYS4D7kxgQAzq5u2dMpYoRKdoNIJEZxv0U6utKWapRx06MO3 -logtostderr=true -v=2 $@
fi

