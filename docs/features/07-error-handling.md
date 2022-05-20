# Error handling

This is small feature of cops-hq, but one which might be handy for simple IaC code. Well, Go is known for explicit error 
handling, which might get too verbose in some cases. For example, canonical Go way would be to check the returned errors
of every method, every time. Unchecked errors are an anti-pattern in Go, so error checking code might get a bit repetitive. 

If most infra as code apps, we simply want to fail and stop the program on any error which occurs. For this
purpose, you can use the global `error_handling.PanicOnAnyError` variable. This setting will propagate through any 
HQ functionality you use.

Note: using the panic on any error mode will make your code difficult to parallelize, because this is a global (shared) variable. 
In case you build your IaC as async (parallel) code, don't use this setting. 