#!/bin/bash
if [ "$1" == "dyn" ]; then
  go build
else
  CGO_ENABLED=0 go build -a -ldflags -d
fi

mv kogia /kogia
