# Sops recipe

Using Sops is directly in-built into the Hq object itself. For example, for loading Sops based configs, refer to (05-configuration)[05-configuration.md].

This recipe is for some extensions points and helper functions around Sops workflows. 

## Regenerating MAC values of config files

If you merge config files from multiple branches, you will often get the "invalid MAC" error from Sops. Solution for this is to regenerate the MAC
values for the config files, and this Sops recipe provides a function to do just that on a directory you provide (all .yaml files in the directory 
will be opened and analyzed).

```go 

import (
    copshq "github.com/conplementag/cops-hq/v2/pkg/hq"
    sopshq "github.com/conplementag/cops-hq/v2/pkg/recipes/sops"
)

// ...

    // for example, we have two directories for our configs:
    logrus.Info("[Config] 📁 🔧 Fixing the MAC values... 🔧 ")
    
    sops := sopshq.New(hq.GetExecutor())
    err := sops.RegenerateMacValues(filepath.Join(hq.ProjectBasePath, "config", "first-directory"))
    // ...
	
    err = sops.RegenerateMacValues(filepath.Join(hq.ProjectBasePath, "config", "second-directory"))
    // ...
    
    logrus.Info("[Config] 📁 Done.")

// ...

```
