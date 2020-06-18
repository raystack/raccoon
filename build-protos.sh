#!/bin/bash

GO_REPO_PATH="/temp-build"
git clone git@source.golabs.io:mobile/clickstream-go-proto.git "${GO_REPO_PATH}"
go build 
