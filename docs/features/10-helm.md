# Helm recipe

Dealing with helm for an automated app deployment process in code is not that difficult, it's mostly a one-liner at the end.
But the variable passing to this one-liner could be very challenging, especially if there are many variables to set,
and you need to provide them in development/test time or in the CD process.
The helm cli provides mainly tho ways to set variables. Either you set your variables one by one or you provide a
variables override files stored on disc.
In this recipe we choose to use the variables file pattern internally to set variables. As we want to keep the client code
as simple as possible, this variable file write logic is done within this recipe. The client can simply 
provide a map of strings to set variables.

## basic usage

```go 

// Login/Connect to your kubernetes cluster, where you want to deploy to.

// First create a new helm instance and set the global parameters.
h := helm.New(s.executor, "my-namespace", "my-chartname", filepath.Join(copshq.ProjectBasePath, "helm"))

// Next step is normally to set variables for your helm deployment.
// This is always a simple map of strings, but supporting any simple or complex object. Nested structures are
// also supported. 
// You can also simply load a whole structure of your viper configuration and convert this directly to a map.
// The given map structure must match partially your values structure in your helm default values.yaml to override
// these values.
// For example if your helm root values.yaml looks like this
// ...
// app:
//     key_1: 0
//     key_2: value_two
//     nested_under_app:
//        key_1_of_nested: 1
//        key_2_of_nested: nested_value
// ...
// you need to provide a nested map structure to override the default values as follows:
vars := map[string]interface{}{
		"app": map[string]interface{}{
			"key_1": 1,
			"key_2": "value_two_override",
			"nested_under_app": map[string]interface{}{
				"key_1_of_nested": 2,
				"key_2_of_nested": "nested_value_override"}}}
				
// alternatively you can auto-populate your map of strings from your config (e.g. sops) with the help of viper
vars = viper.GetStringMap("application")
				
err = h.SetVariables(vars)
... do something with the error

// If you skip the SetVariables call, the default values in the helm's values.yaml will be used

// then, we can start the real deployment.
err = h.Deploy()
... do something with the error

```

## Variable auto-populate from config file

In the example above, a code snippet is included on how to auto-populate the variables from a config section. You could also combine
both auto-population and your own keys, since this is just simple map manipulation in Go!