![](assets/logo.png | width=250)

[![](https://godoc.org/github.com/ropenttd/gopenttd?status.svg)](https://godoc.org/github.com/ropenttd/gopenttd)

**gopenttd** is a simple Golang library for querying OpenTTD game servers.

## Command Line Usage

There's a command line utility called [_openttd\_scrape_](cmd/openttd_scrape) which produces nice JSON objects for you to parse externally. See the documentation there for more information.

You can run it with something like the following:
```
go get github.com/ropenttd/gopenttd
go run github.com/ropenttd/gopenttd/cmd/openttd_scrape
```

## API

There are two APIs:
* the client protocol, which is a UDP-based polling protocol (i.e one shot) and can communicate with any server
* the Admin protocol, which is a TCP based protocol with significantly more capability, but requires that you have the admin password of the server you are connecting to.

Please see the [godoc](https://godoc.org/github.com/ropenttd/gopenttd) for further information on the API.

### Client Protocol

Here's a brief example:
```go
import "github.com/ropenttd/gopenttd/pkg/gopenttd"

result, err := gopenttd.ScanServer("s1.ttdredd.it", 3979)
```

### Admin Protocol

The Admin Protocol is a connection based protocol that you communicate to using a combination of a Write command and a channel reader for responses.

There is a helper "ScanServerAdm" function that acts very similarly to the ScanServer function, except it returns significantly more data.

```go
import "github.com/ropenttd/gopenttd/pkg/gopenttd"

result, err := gopenttd.ScanServerAdm("s1.ttdredd.it", 3977, "password")
```

Please see the [godoc](https://godoc.org/github.com/ropenttd/gopenttd) for help using the rest of the API.