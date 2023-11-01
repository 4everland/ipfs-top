# Run with Docker Compose
```
$ cd deployment/docker/all-in-one/
$ docker-compose up -d
```
## Examples:

### RPC
```
$ echo 'Hello ipfs-top!' > hello.txt
$ curl -X POST '127.0.0.1:8086/api/v0/add?raw-leaves=true&cid-version=1' --form 'file=@"./hello.txt"'
{"Hash":"bafkreifx3lv5up5e2ic4ybdmkxce6w77x2autjuvk7yyoeaes7kbliaed4","Name":"hello.txt","Size":"16"}
```
### GATEWAY
```
$ curl 127.0.0.1:8080/ipfs/bafkreifx3lv5up5e2ic4ybdmkxce6w77x2autjuvk7yyoeaes7kbliaed4
Hello ipfs-top!
```
