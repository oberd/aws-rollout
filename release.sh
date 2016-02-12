#!/bin/bash

rm -rf dist && mkdir dist
tar -zcvf dist/aws-rollout-linux-x86_64.tar.gz aws-rollout
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)
github-release "oberd/aws-rollout" "$latest_tag" "$(git rev-parse --abbrev-ref HEAD)" "" 'dist/*'