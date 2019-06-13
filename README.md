# props-transaction-processor

![Props Token](https://propsproject.com/static/images/main-logo.png)
 
this propject was generated with goprops-template

## Building docker image 
In root directory of project run
```bash
make docker-image
```

Now bring up sawtooth network
```bash
cd examples/sample-networkd
 docker-compose -f sample-network.yaml up --force-recreate
```

## Generate Protobuffer bindings
Run command in root level of project 

- Go 
```bash
mkdir -p ./core/proto/pending_props_pb
protoc -I ./protos ./protos/earning.proto ./protos/payload.proto --go_out=./core/proto/pending_props_pb
```
- JS (Node)
```bash
protoc -I ./protos ./protos/earning.proto ./protos/payload.proto --js_out=import_style=commonjs,binary:OUTPUTDIR
```

- JS (Browser) using [protobuf.js](https://github.com/dcodeIO/ProtoBuf.js) 
```bash
pbjs -t static-module -w commonjs -o OUTPUTDIR ./protos/earning.proto ./protos/payload.proto
```


## Using dev-cli 

### Pre-reqs

Install dependencies with these commands
```bash
cd dev-cli
yarn
```

### Usage

Issue an earning, a browser window should open to the transaction status
```bash
node index.js pending-props -a 100 -r 0x42EB768f2244C8811C63729A21A3569731535f06 
```

If that transaction was successful, lets query for it using the state address 
```bash
 node index.js query-state -a e0a87c93dc0ccef35a577fd05ac533eb2f6e601917d26c3ba8be75f4ab14f9d39370a0
```