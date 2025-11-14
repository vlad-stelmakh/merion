Linux, using apt or apt-get, for example:
```
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```
MacOS, using Homebrew:
```
brew install protobuf
protoc --version  # Ensure compiler version is 3+
```

Windows, using Winget
```
>d
> protoc --version # Ensure compiler version is 3+
```

Generate pb:
```
protoc --go_out=. --go-grpc_out=. calculator.proto
```