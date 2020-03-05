#!/usr/bin/env bash
echo "remove application..."
rm bin/dora-logging
echo "build application..."
#go build -o bin/dora-logging cmd/main.go
env GOOS=linux GOARCH=amd64 go build -o bin/dora-logging cmd/main.go
echo "upFie..."
rsync -avzP -r --delete \
    configs \
    bin/dora-logging \
    sontc@110.35.75.40:/home/sontc/truonglv/Dora-Logging/