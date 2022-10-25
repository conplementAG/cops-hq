# Upgrading

To update the cops-hq in your code to a new version, you need to perform a minor or a major upgrade.

## Package versioning basics

In Go, each major library version changes the module path as well (in your go.mod and in your import() statements in *.go files). For example, 
for major version v0 and v1, you will have the module path github.com/conplementag/cops-hq. For versions v2 and above, you will need to explicitly
specify github.com/conplementag/cops-hq/v2 (or v3, etc.).

## Minor upgrade

For a minor upgrade, simply run `go get <<cops hq module you are already using>>`. For example, if you use github.com/conplementag/cops-hq/v2 in your code,
and you want to upgrade from v2.0.0 to v2.0.1, you could simply run  `go get github.com/conplementag/cops-hq/v2`. This will automatically
update your go.mod and go.sum files, and the references in import() statements will stay the same.

## Major upgrade

For a major upgrade, you need to first manually replace all the module imports in your code. For example, if upgrading from v1 to v2, you need to change this

```go
import (
    "github.com/conplementag/cops-hq"
)
```

to

```go
import (
    "github.com/conplementag/cops-hq/v2"
)
```

Then, run `go get github.com/conplementag/cops-hq/v2` and `go mod tidy` from the root of your project. 