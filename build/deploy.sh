#!/usr/bin/env bash
echo "remove application..."
rm bin/dora-logging
echo "build application..."
go build -o bin/dora-logging cmd/main.go
echo "upFie..."
rsync -avzP -r --delete \
    configs \
    bin/dora-logging \
    doraemon@110.35.75.45:/home/doraemon/truonglv/Dora-Logging/