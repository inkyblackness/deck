#!/bin/sh

pushd $(dirname "${0}") > /dev/null
DECK_BASE=$(pwd -L)
export GOPATH=$DECK_BASE

COMPONENTS="construct chunkie hacker shocked-client"

rm -rf $DECK_BASE/src

for name in $COMPONENTS
do
   go get -d github.com/inkyblackness/$name
done


for name in $(find $DECK_BASE/src -iname ".git")
do
   rm -rf $name
done

popd > /dev/null
