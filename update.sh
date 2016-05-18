#!/bin/sh

pushd $(dirname "${0}") > /dev/null
DECK_BASE=$(pwd -L)
export GOPATH=$DECK_BASE

COMPONENTS="construct chunkie hacker shocked-server shocked-client"

rm -rf $DECK_BASE/src

for name in $COMPONENTS
do
   go get -d github.com/inkyblackness/$name
done
go get -d github.com/inkyblackness/shocked-client/app/shocked-client-console


for name in $(find $DECK_BASE/src -iname ".git")
do
   rm -rf $name
done

popd > /dev/null
