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

## approach2-extended(varied dependency check)

### sync modules
go work sync

### Build plugins (each uses its own uuid version)
go build -o plugins/addheader/addheader ./plugins/addheader
go build -o plugins/removeheader/removeheader ./plugins/removeheader

add the generated plugins to host/plugins directory

### Build host (uses uuid v1.0.0)
go build -o host ./host

### Run host
./host

## Yaegi Scription Approach

go mod tidy

go run host/main.go

# Custom Plugin Loader

## 1. Build plugins
cd plugins/addheader
go build -o ../../bin/addheader
cd ../removeheader
go build -o ../../bin/removeheader
cd ../addremoveheader
go build -o ../../bin/addremoveheader

## 2. Build host
cd ../../host
go build -o host

## 3. Run
cd ../../host
./host

