@echo off

set SEVENZIP="C:\Program Files\7-Zip\7z.exe"
set DECK_BASE="%~dp0"

pushd "%DECK_BASE%"
set OLDGOPATH=%GOPATH%
set GOPATH=%DECK_BASE%

echo Cleaning output directories...
rmdir /s /q bin >NUL
rmdir /s /q pkg >NUL
rmdir /s /q inkyblackness-deck >NUL
del inkyblackness-deck.*


echo Building executables...
cd %DECK_BASE%src\github.com\inkyblackness\construct
go test ./...
go install

cd %DECK_BASE%src\github.com\inkyblackness\chunkie
go test ./...
go install

cd %DECK_BASE%src\github.com\inkyblackness\hacker
go test ./...
go install

cd %DECK_BASE%src\github.com\inkyblackness\shocked-server
go test ./...
go install


echo Copying resources...

mkdir %DECK_BASE%bin\client
xcopy %DECK_BASE%src\github.com\inkyblackness\shocked-client\www %DECK_BASE%\bin\client /s /e

copy %DECK_BASE%LICENSE %DECK_BASE%\bin
xcopy %DECK_BASE%resources %DECK_BASE%\bin /s /e


echo Creating package...
move %DECK_BASE%bin %DECK_BASE%inkyblackness-deck

cd %DECK_BASE%
%SEVENZIP% a -r -tzip -bd inkyblackness-deck.zip inkyblackness-deck\*


set GOPATH=%OLDGOPATH%
set OLDGOPATH=
popd
