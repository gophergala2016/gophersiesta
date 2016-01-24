# GOPHERSIESTA

A manager/service for configurations files and properties. GopherSiesta is composed of two parts, a server and an example command line client to communicate with the server's API. The goal of GopherSiesta is to make it easier to manage the configurations of all your services. GopherSiesta client will run prior to your application and fetch the corresponding configuration for your service.

![alt tag](assets/gopherswrench.jpg)

## Installation

As the project uses [godep](https://github.com/tools/godep) to make builds reproducibly

```
go get github.com/gophergala2016/gophersiesta
```

## Run

```
cd cmd/start-we-server

go run main.go
```

## API

### Get template for :appname
```
http://gophersiesta.herokuapp.com/conf/:appname
```
Retrieve the full template file for the application.

*Example*
```
GET http://gophersiesta.herokuapp.com/conf/app1

application:
    name: "App1"
    version: 0.0.1

datasource:
    url: ${DATASOURCE_URL:jdbc:mysql://localhost:3306/shcema?profileSQL=true} # has default value
    username: ${DATASOURCE_USERNAME} # has no default value. If no value is passed should producer error if validated
    password: ${DATASOURCE_PASSWORD}
```


### Retrieve list of placeholders
Get the list of all possible variables of the template.

```
http://gophersiesta.herokuapp.com/conf/:appname/placeholders
```

*Example*

```
GET http://gophersiesta.herokuapp.com/conf/app1/placeholders
{
  "placeholders": [
    {
      "placeholder": "DATASOURCE_URL",
      "property_value": "${DATASOURCE_URL:jdbc:mysql://localhost:3306/shcema?profileSQL=true}",
      "property_name": "datasource.url"
    },
    {
      "placeholder": "DATASOURCE_USERNAME",
      "property_value": "${DATASOURCE_USERNAME}",
      "property_name": "datasource.username"
    },
    {
      "placeholder": "DATASOURCE_PASSWORD",
      "property_value": "${DATASOURCE_PASSWORD}",
      "property_name": "datasource.password"
    }
  ]
}
```

### Retrieve current values of placeholders for :appname  
Get the values that are going to be used to generate the template. Labels override previous values. 

```
http://gophersiesta.herokuapp.com/conf/:appname/values?labels=:label1,:label2
```


*Example*
```
GET http://gophersiesta.herokuapp.com/conf/app1/values
{
  "values": [
    {
      "value": "FOOBAR",
      "name": "datasource.password"
    },
    {
      "value": "GOPHER",
      "name": "datasource.username"
    }
  ]
}

GET http://gophersiesta.herokuapp.com/conf/app1/values?labels=dev
{
  "datasource.username": "GOPHER-dev",
  "datasource.password": "LOREM"
}

GET http://gophersiesta.herokuapp.com/conf/app1/values?labels=prod
{
  "datasource.username": "GOPHER-prod",
  "datasource.url": "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true",
  "datasource.password": "IPSUM"
}

GET http://gophersiesta.herokuapp.com/conf/app1/values?labels=dev,prod
{
  "datasource.username": "GOPHER-prod",
  "datasource.url": "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true",
  "datasource.password": "IPSUM"
}

GET http://gophersiesta.herokuapp.com/conf/app1/values?labels=prod,dev
{
  "datasource.username": "GOPHER-dev",
  "datasource.url": "jdbc:mysql://proddatabaseserver:3306/shcema?profileSQL=true",
  "datasource.password": "LOREM"
}

```

### Retrieve list of labels
Get the list of all possible variables of the template.
```
http://gophersiesta.herokuapp.com/conf/:appname/labels
```

*Example*

```
GET http://gophersiesta.herokuapp.com/conf/app1/labels
{
  "labels": [
    "default",
    "dev",
    "prod"
  ]
}
```


## TODO

+ Render the template conf applying the saved values given some labels


The Gopher character is based on the Go mascot designed by Renée French and copyrighted under the Creative Commons Attribution 3.0 license.
