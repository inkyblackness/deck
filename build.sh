#!/bin/sh

pushd $(dirname "${0}") > /dev/null
DECK_BASE=$(pwd -L)

export GOPATH=$DECK_BASE

echo Cleaning output directories...
rm -rf bin
rm -rf dist
rm -rf pkg

mkdir -p $DECK_BASE/dist/linux/inkyblackness-deck
mkdir -p $DECK_BASE/dist/win/inkyblackness-deck


echo Building executables...

function buildNative() {
   local name=$1

   echo "Building " $name
   go build -o $DECK_BASE/dist/linux/inkyblackness-deck/$name .
   GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o $DECK_BASE/dist/win/inkyblackness-deck/$name.exe .
}

for name in "construct" "chunkie" "hacker" "shocked-client"
do
   cd $DECK_BASE/src/github.com/inkyblackness/$name
   buildNative $name
done

echo Copying resources...

for os in "linux" "win"
do
   packageDir=$DECK_BASE/dist/$os/inkyblackness-deck

   cp $DECK_BASE/LICENSE $packageDir
   cp -R $DECK_BASE/resources/* $packageDir
done


echo Creating packages...

cd $DECK_BASE/dist/linux
tar -cvzf $DECK_BASE/dist/inkyblackness-deck.linux64.tgz ./inkyblackness-deck

cd $DECK_BASE/dist/win
zip -r $DECK_BASE/dist/inkyblackness-deck.win64.zip ./inkyblackness-deck

popd > /dev/null
