# Contribution

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