#!/bin/bash

eval "$(jq -r '@sh "NAME=\(.name)"')"

jq -n --arg input "$NAME" '{"output_value":$input}'