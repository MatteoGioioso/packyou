#!/usr/bin/env sh

rm -rf dist/ &&

go build ../pku &&

./pku \
  -e services/functionOne/index.js \
  -r /Users/madeo/Development/Tools/packyou/__test__/testProject/project-root \
  -o dist &&

rm pku
