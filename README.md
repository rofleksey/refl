# Refl

A simple scripting language interpreter written in Go.

## Overview

Refl is an interpreted language, featuring:

- **Simple Types**: `number`, `string`, `object`, `function`, `iterator`, `error`, `nil`
- **Dynamic Typing**: No type declarations needed
- **First-Class Functions**: Functions are values that support closures
- **No arrays**: Objects keys can be of any value except `nil`
- **No Booleans**: `0`, `nil` and `""` are "false", anything else is "true"

## Syntax Examples

```javascript
# Variables
var x = 1
x = 2  # global assignment

# No arrays, only objects
var obj = {a: 1, "b": 2}
var arr = {1, "two", 3}  # "array" syntax sugar, 0-indexed

# All functions are anonymous
var add = fun(a, b) { return a + b }
var result = add(1, 2)

# Control flow
if x > 0 {
    # do something
} elif x == 0 {
    # something else
} else {
    # default
}

# Loops
while condition {
    # loop body
}

# Iteration
for key, value in obj {
    # iterate over object or string
}

"6" + 7 # "67"
```

## Using from Go

You can embed Refl in your Go application:

```go
package main 

import (
    "context"
    "fmt"
    "refl/parser"
    "refl/runtime"
    "refl/runtime/eval"
)

// Example usage
func main() {
	p := parser.New()
	program, err := p.Parse(`var x = 10; x * 2`)
	if err != nil {
		panic(err)
	}

	env := runtime.NewEnvironment(nil)
	
    result, err := eval.Eval(context.Background(), program, env)()
    if err != nil {
        panic(err)
    }
		
    fmt.Println(result.String()) // "20"
}
```

## Standard Library

Refl includes several built-in modules:

* math - Mathematical functions (abs, floor, random, etc.)
* strings - String manipulation (upper, split, contains, etc.)
* io - Input/output functions (print, println, printf)
* Global functions: type(), str(), len(), clone(), eval(), panic()
* Error handling: newerr(), iserr(), errfmt()