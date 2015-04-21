#!/bin/bash
export PGHOST="postgres://patrol:patrol@localhost/patrol"
export CACHEDSN="redis://127.0.0.1:6379/1"
export QUEUEDSN="redis://127.0.0.1:6379/2"
$GOPATH/bin/goconvey -depth=1