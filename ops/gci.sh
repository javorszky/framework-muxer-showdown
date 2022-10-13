#!/bin/sh

gci write -s Standard -s Default -s "Prefix(github.com/suborbital)"  -s "Prefix(github.com/suborbital/framework-muxer-showdown)" --skip-generated --custom-order $(find . -not \( -path ./vendor -prune \) -type f -name '*.go')
