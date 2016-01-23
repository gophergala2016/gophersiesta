# GOPHERSIESTA

## Installation

```
go get github.com/gin-gonic/gin
```

## Run

```
cd server

go run server.go
```

## Usage

```
curl http://localhost:4747/conf/app1
curl http://localhost:4747/conf/app1/values
curl http://localhost:4747/conf/app1/placeholders?labels=prod

```

## TODO

+ Post value for placeholders given some labels
+ Render the templace conf applying the saved values given some labels