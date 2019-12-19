#!/usr/bin/env bash

version=$1

if [[ -z "$version" ]]; then
    echo -e "Enter version"
    exit 1
fi

docker commit loggo mrandersen7/logc-go:${version}; docker push mrandersen7/logc-go