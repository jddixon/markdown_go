#!/usr/bin/env bash

echo 
echo Checking our output against standard
cmd/xgoMarkdown/xgoMarkdown -v -i commonmark -o tmp spec.txt
diff commonmark/spec.html tmp/spec.html
