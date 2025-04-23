# GoPluginPOC
Simple Go Plugin to add a header

## Build Add Header Plugin 

go build -buildmode=plugin -o addheader.so ./addheader/addheader.go

copy the generated .so file to plugin-host/ directory.

## Build Plugin Host

go build -o plugin-host

## Run Plugin Host

./plugin-host 