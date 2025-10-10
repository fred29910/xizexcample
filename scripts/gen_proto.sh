#!/bin/bash
protoc --go_out=. --go_opt=module=xizexcample \
    api/proto/game.proto
