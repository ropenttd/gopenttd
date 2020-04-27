# gopenttd

[![](https://godoc.org/github.com/ropenttd/gopenttd?status.svg)](https://godoc.org/github.com/ropenttd/gopenttd)

_gopenttd_ is a simple Golang library for querying OpenTTD game servers.

## Usage

There's a command line utility called [_openttd\_scrape_](cmd/openttd_scrape) which produces nice JSON objects for you to parse externally. See the documentation there for more information.

You can run it with something like the following:
```
go get github.com/ropenttd/gopenttd
go run github.com/ropenttd/gopenttd/cmd/openttd_scrape
```

### API

Here's a brief example:
```go
import "github.com/ropenttd/gopenttd/pkg/gopenttd"

result, err := gopenttd.ScanServer("s1.ttdredd.it", 3979)
```

Please see the [godoc](https://godoc.org/github.com/ropenttd/gopenttd) for further information on the API.
