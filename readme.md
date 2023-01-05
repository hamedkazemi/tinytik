
> ### TinyTik API - Very Simple Url Shortener written in Golang

# How it works
```
.
├── main.go
├── common
│   ├── utils.go       // small tools function
│   └── db.go          // DB connect manager
│   └── logrus.go      // logrus log manager instance
│   └── kafka.go       // kafka conenction manager
│   └── config.go      // configuration manager
└── modules
    └── moduleName
        ├── controller.go   // controllers functions used by router
        ├── models.go       // data models define & DB operation
        ├── serializers.go  // response computing & format
        ├── routers.go      // router binding
        ├── middlewares.go  // put the before & after logic of handle request
        └── validators.go   // form/json checker
```

# Getting started

## Install the Golang
https://golang.org/doc/install
## Environment Config
make sure your ~/.*shrc have those variable:
```
➜  echo $GOPATH
/Users/[user]/test/
➜  echo $GOROOT
/usr/local/go/
➜  echo $PATH
...:/usr/local/go/bin:/Users/zitwang/test//bin:/usr/local/go//bin
```
## Install Go Modules
I used Native Go Modules to manage the packages.

Using fresh we have auto-reloading functionality.
https://github.com/pilu/fresh
```
go get -u ./...
```

## Development Start

First you need specify the configuration in config.toml file:

```
[database]
server = "127.0.0.1" # Mysql Database Server
port = "3306" # Mysql Database Port
database = "take10dashboard" # Mysql Database Name
user = "root" # Mysql Database User
password = "root" # Mysql Database Password
debug = true # Mysql Database Debug

[app]
name = "Take10Dashboard" # Application Name
port = "0.0.0.0:8080" # Application Port And IP
proxy = "" # Outgoing Proxy
environment = "development" # development,production

[redis]
connectionString = "127.0.0.1:6379"

[kafka]
ip = "" # Kafka IP ( ignored if empty, eg. 0.0.0.0 )
port = "" # Kafka Port ( ignored if empty, eg. 9092 )
topic = "" # Kafka Default Topic To Publish ( ignored if empty, eg. Test )
```

```
➜  fresh
```

## Testing
```
➜  go test -v ./... -cover
```

## To-Do
- Clean Structure
- CI/CD
- Write Tests
