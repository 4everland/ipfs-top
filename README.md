# IPFS-Top
IPFS Top is a node program compatible with IPFS RPC and Gateway. It supports any abstract Block Store implementation. It contains the following parts:

RPC&Gateway: This is the part compatible with IPFS RPC and Gateway
Indexer: This is an index layer that maps the relationship between the ID of any storage layer and the IPFS CID.
Block Store: This is an abstract implementation of the original Block Store of IPFS. It will support RocksDB, LevelDB, AWS S3, ARWeave, etc.

> This branch is about s3 file bridge to arweave

#### Build
```
go build -o bridge tasks/cmd/main.go
```

#### Config

```

s3: // s3 bucket config
  endpoint: xxx 
  region: xxx
  bucket: xxx
  accessKey: xxx
  secretKey: xxx
arseeding:
  gateway_addr: xxx //arseeding rpc 
  mnemonic: xxx     //the eth mnemonic  is used to sign data
```

#### Run

sync one file
```
# ./bridge once -h
send one s3 object to arseeding

Usage:
   once [flags]

Flags:
  -h, --help         help for once
  -k, --key string   s3 object key

Global Flags:
  -c, --config string   config file (default "config.yaml")
  -p, --pass string     mnemonic password
```
batch sync
```
# ./bridge batch -h 
batch send s3 objects to arseeding

Usage:
   batch [flags]

Flags:
  -s, --batch-file string   batch file path
  -h, --help                help for batch

Global Flags:
  -c, --config string   config file (default "config.yaml")
  -p, --pass string     mnemonic password
```
batch file example:
```
example/test1
example/test2
example2/test1
```
