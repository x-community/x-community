# Introduction
Mail service used to send mail to users in community platform.

## Protobuf
```
protoc --proto_path=proto --micro_out=proto --go_out=proto proto/mail.proto
```

## Configuration
### config.toml
```toml
address = ":8080"

[db]
host = "127.0.0.1"
database = "default"
username = "root"
password = ""
show_log = true
```

## Usage
```
go run .
```