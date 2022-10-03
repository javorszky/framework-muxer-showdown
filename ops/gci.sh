#!/bin/sh

gci write -s Standard -s Default -s "Prefix(github.com/suborbital)"  -s "Prefix(github.com/suborbital/framework-muxer-showdown)" --skip-generated --custom-order $(find . -type f -name '*.go')
