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
    # iterate over iterator, object or string
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
	program, err := p.Parse(`io.println("Hello world!")`)
	if err != nil {
		panic(err)
	}

	env := runtime.NewEnvironment(nil)

	evaluator := eval.New(context.Background(), program, env)
	result, err := evaluator.Run()
	if err != nil {
		panic(err)
	}
		
    fmt.Println(result.String()) // "20"
}
```

## Standard Library

Refl includes several built-in modules:

* `math` - Mathematical functions (`abs`, `floor`, `random`, etc.)
* `strings` - String manipulation (`upper`, `split`, `contains`, etc.)
* `time` - Time functions (`now`, `parse`, `format`, `sleep`, etc.)
* `io` - Input/output functions (`print`, `println`, `printf`, e.t.c.)
* `events` - Event loop functions (`schedule`, `register`, e.t.c.)
* `errors` - Creating errors: (`new`, `is`, `fmt`, e.t.c.)

Global functions:
* `range` - creates iterators over integers, same as in python
* `type` - outputs type of the argument
* `str` - converts the argument to string
* `len` - outputs the length of an indexable argument
* `clone` - creates a deep copy of the argument, functions are copied by reference
* `eval` - evaluates a refl code
