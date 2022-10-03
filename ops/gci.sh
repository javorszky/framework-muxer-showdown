#!/bin/sh

gci write -s Standard -s Default -s "Prefix(github.com/suborbital)"  -s "Prefix(github.com/suborbital/framework-muxer-showdown)" --NoInlineComments --NoPrefixComments "$(find ./ -type f -name '*.go')"
