#!/bin/sh

pushd $(dirname "${0}") > /dev/null
DECK_BASE=$(pwd -L)

export GOPATH=$DECK_BASE

echo Cleaning output directories...
rm -rf bin
rm -rf pkg
rm -rf inkyblackness-deck
rm -rf inkyblackness-deck.*


echo Building executables...

cd $DECK_BASE/src/github.com/inkyblackness/construct
go test ./...
go install

cd $DECK_BASE/src/github.com/inkyblackness/chunkie
go test ./...
go install

cd $DECK_BASE/src/github.com/inkyblackness/hacker
go test ./...
go install

cd $DECK_BASE/src/github.com/inkyblackness/shocked-server
go test ./...
go install


echo Copying resources...

mkdir -p $DECK_BASE/bin/client
cp -R $DECK_BASE/src/github.com/inkyblackness/shocked-client/www/* $DECK_BASE/bin/client

cp $DECK_BASE/LICENSE $DECK_BASE/bin
cp -R $DECK_BASE/resources/* $DECK_BASE/bin


echo Creating package...
cd $DECK_BASE
mv $DECK_BASE/bin $DECK_BASE/inkyblackness-deck
tar -cvzf $DECK_BASE/inkyblackness-deck.tgz ./inkyblackness-deck

popd > /dev/null
