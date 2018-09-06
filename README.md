# go-pubg
PUBG API Wrapper

## documentation
Documentation can be found at [here](https://godoc.org/github.com/redorb/go-pubg).
This is a hard fork of https://github.com/LtSnuggie/pubg to remove the callback style of requests.

## testing
In order to run the tests locally you will need to add a conf.json file in the root folder

In the file you will need to add the following:
>{
>   "key"       : "[API key]",
>}

## dependencies
Though this client only has one [dependency](https://github.com/google/go-querystring) the project uses the new go module system. More information on go modules can be found [here](https://github.com/golang/go/wiki/Modules).
