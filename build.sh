#!/bin/bash

dotenv () {
  echo "Reading .env variables..." &&
    eval "$(cat .env)" &&
    test "$BRUNCHFEED_TEMPLATE_FILE" &&
    BRUNCHFEED_TEMPLATE="$(cat $BRUNCHFEED_TEMPLATE_FILE)" || true
}

build () {
  (echo "Fetching Brunch RSS feed..." &&
      brunchfeed fetch \
        -id "$BRUNCHFEED_BRUNCH_ID" \
        -src "$BRUNCHFEED_DATA_PATH" &&

    echo "Building posts data..." &&
      brunchfeed build \
        -src "$BRUNCHFEED_DATA_PATH" \
        -dest "$BRUNCHFEED_POST_PATH" \
        -template "$BRUNCHFEED_TEMPLATE" &&

    echo "Generating pages..." &&
      eval "$BRUNCHFEED_GENERATE")
}


dotenv && build
