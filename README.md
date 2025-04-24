# GoPluginPOC
Simple Go Plugin to add a header

## Build Add Header Plugin 

go build -buildmode=plugin -o addheader.so ./addheader/addheader.go

copy the generated .so file to plugin-host/ directory.

## Build Plugin Host

go build -o plugin-host

## Run Plugin Host

./plugin-host 

# Approach 2 - Hashicorp go plugin

# Build Plugins and Host

go mod tidy
go build -o bin/addheader ./plugins/addheader
go build -o bin/removeheader ./plugins/removeheader
go build -o bin/host ./host

cd bin
./host

