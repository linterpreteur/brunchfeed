#!/bin/bash

dotenv () {
  echo "Reading .env variables..." &&
    eval "$(cat .env)" &&
    [ -n "$BRUNCHFEED_TEMPLATE_FILE" ] &&
    BRUNCHFEED_TEMPLATE="$(cat $BRUNCHFEED_TEMPLATE_FILE)"
}

gopath () {
  echo "Verifying GOPATH..." &&
    [ -z $GOPATH ] && export GOPATH=$(go env GOPATH)
}

build () {
  (cd $GOPATH/src/github.com/linterpreteur/brunchfeed &&

    echo "Fetching Brunch RSS feed..." &&
    go run main.go fetch \
      -id "$BRUNCHFEED_BRUNCH_ID" \
      -src "$BRUNCHFEED_DATA_PATH" &&

    echo "Building posts data..." &&
    go run main.go build \
      -src "$BRUNCHFEED_DATA_PATH" \
      -dest "$BRUNCHFEED_POST_PATH" \
      -template "$BRUNCHFEED_TEMPLATE" &&

    echo "Generating pages..." &&
      eval "$BRUNCHFEED_GENERATE")
}

dotenv
gopath
build
