# idmapper

simple web app using IDMappers

## IDMappers

IDMapper is in-memory cache for mapping IDs to Names. Can be used to cache lists of IDs and Names (eg country list, languages list). More about IDMappers [here](https://github.com/danielkraic/idmapper/tree/master/idmapper) 

### IDMappers reloading

IDMappers are reloaded automaticaly in background using [scheduler](https://github.com/danielkraic/idmapper/tree/master/scheduler)

## API

```
GET /health
GET /version
GET /metrics

GET /v1/country/{countrycode}
GET /v1/currency/{currencycode}
GET /v1/language/{languagecode}
```

Example:

```bash
curl localhost:8080/v1/country/sk
curl localhost:8080/v1/country/sk
```

## Developement

### Requirements

* golang 1.12
* make
* git

### Building

```bash
# build webapp
make build
# run tests
make test
# run tests and build webapp 
make all
# build docker image
make docker
```

### Usage

```bash
./idmapperapp -h
```
```
Usage of ./idmapperapp:
  -a, --addr string     HTTP service address. (default "0.0.0.0:80")
  -c, --config string   path to config file
      --config-check    check configuration
  -p, --print-config    print configuration
```

Example:

```bash
./idmapperapp -c config-example.yaml -a localhost:8080
```

### Documentation

https://godoc.org/github.com/danielkraic/idmapper

### Configuration

There are two ways of providing configuration: Using config file and using environmental variables. Both ways can be combined.

#### Using config file

See [config-example.yaml](config-example.yaml) for available options.

```bash
./idmappersapp --config config-example.yaml
```

#### Using environmental variables

Prefix for all environmental variables is `IDMAPPER_`.

Examples:
```
IDMAPPER_ADDR="localhost:8080" ./idmapperapp --print-config --config-check
IDMAPPER_LOGGER_JSON=true ./idmapperapp --print-config --config-check
IDMAPPER_REDIS_PASSWORD="secretpass" ./idmapperapp --print-config --config-check

```

### Handling dependencies

Save all dependecies to `vendor` directory

```bash
go mod vendor
```