#!/bin/bash
export TESTING=TRUE
export DB_DSN=postgres://patrol:patrol@localhost/patrol
export ERGOQ_DSN=redis://localhost:6379
export CACHER_DSN=redis://localhost:6379/1?prefix=patrol
export SECRET_KEY=6nD98ZbRHR4MPFCtEj85ZliUNakJkvUZQY5TLUWstlzg7ALH1u7zTcl4IVQYOmpL
#go test -v
goconvey . -secret_key=baGCbYmpdRxeSZ2rJYS4D7kxgQAzq5u2dMpYoRKdoNIJEZxv0U6utKWapRx06MO3