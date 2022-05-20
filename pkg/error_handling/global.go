package error_handling

// PanicOnAnyError will issue panic if any of cops-hq methods returns an error. Beware that using this method will stop
// the program execution, and you will not be able to inspect the error in your code (although panic will log
// everything to stdout/stderr, and the error will be written into logs too). This mode might be interesting for code
// equivalent to Bash scripts running with 'set -e'. Setting panic mode only affects the commands / methods executed
// after calling this method, and the mode can also be reverted by setting to false.
var PanicOnAnyError = false
