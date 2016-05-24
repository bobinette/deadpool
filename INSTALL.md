# Deadpool

## Setup
The servers are written in Go. The communication between server and client uses protocol buffers and gRPC.

If you use Vagrant and Ansible, have a look at the `deploy` folder, otherwise simply follow the instructions below.

### Go
The server is written in `Go`, hence you need to install it. First, install the GVM (you need git installed):
```bash
curl -s -S -L https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer | bash
```
Then, download and install `go1.6.1` (you need bison):
```bash
gvm install go1.6.1 --binary
```

I recommend the following project structure, to compel with Go import guidelines:
```
go (or whatever you want here) <-- your $GOPATH
|__ src
    |__ github.com
        |__ bobinette
            |__ deadpool
```
Once you are done, you can add the following to your `.bashrc`:
```bash
source "/home/vagrant/.gvm/scripts/gvm" # Normally added by the gvm
gvm use go1.6.1
export GOPATH=/path/to/go # Setup the GOPATH as per the previous schema
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin # To use the binaries
```

### Protocol buffer
Start by installing protocol buffer and protoc. You can find all the information [here](https://github.com/google/protobuf), but basically:

Install protoc
```bash
wget https://github.com/google/protobuf/releases/download/v3.0.0-beta-2/protoc-3.0.0-beta-2-linux-x86_64.zip
unzip protoc-3.0.0-beta-2-linux-x86_64.zip -d $HOME/protoc
```

Then put
```bash
export PATH=$PATH:$HOME/protoc
```
in your bash profile.

Install protobuf and protoc for go:
```bash
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

### Install dependencies
```bash
go get -u google.golang.org/grpc
```
or
```
make install
```
