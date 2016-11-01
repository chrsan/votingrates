#!/bin/bash

OS=(darwin linux windows)

for x in ${OS[@]}; do
    OUT="votingrates-${x}"
    if [ ${x} == 'windows' ]; then
        OUT="${OUT}.exe"
    fi

    echo "Building for ${x}"
    GOOS=${x} GOARCH=amd64 go build -o $OUT
done
