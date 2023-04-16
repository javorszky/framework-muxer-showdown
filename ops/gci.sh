#!/bin/sh

gci write -s Standard -s Default -s "Prefix(github.com/javorszky)"  -s "Prefix(github.com/javorszky/framework-muxer-showdown)" --skip-generated --custom-order $(find . -not \( -path ./vendor -prune \) -type f -name '*.go')
