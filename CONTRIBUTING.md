# Contribution

## Local Development with a "client" project

If you want to directly work on both cops-hq and your IaC project, without having the need to publish cops-hq changes to upstream, 
a `go.work` file can be used (in the root of your project):

``` 
go 1.18

use (
   .
    ../../../Conplement/CoreOps/cops-hq
)
```

The path should simply point to a local directory where you stored the cops-hq sources. Go will do the rest.

## Dependencies

### Adding new dependencies

To add a new dependency, simply specify the dependency directly in the source code, for example:

``` go
import (
    "github.com/sirupsen/logrus"
)
```

then, execute

```shell
go mod tidy
```

Required dependency, and all transient dependencies, will be recorded in the go.mod and go.sum files.

## Tests

- Please use testify/assert assertions!
- To check coverage, use (from the project root):

````shell 
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
````

## Error handling

We use a dual approach of handling errors in cops-hq: errors are returned to the caller, but we also support panics 
in case the global settings PanicOnAnyError is set. This means that any public member of cops-hq, exposing errors as 
return values, should be wrapped in something similar to this:

```go
	if err != nil && error_handling.PanicOnAnyError {
		logurs.Fatal(err)
		panic(err)
	}

	return err
```

Please check code for references on this behaviour, as any public member implements it. 