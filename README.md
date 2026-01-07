# Refl

A simple scripting language interpreter written in Go.

## Overview

Refl is an interpreted language, featuring:

- **Simple Types**: `number`, `string`, `object`, `function`, `error`, `nil`
- **Dynamic Typing**: No type declarations needed
- **First-Class Functions**: Functions are values that support closures
- **No arrays**: Objects keys can be of any value except `nil`
- **No Booleans**: `0`, `nil` and `""` are "false", anything else is "true"
- **Event Loop**: Built-in asynchronous programming with promises and events
- **Coroutines**: Lightweight concurrency with the refl() function

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
for i, val in range(0, 10) {
    io.println(i, val)
}

# Create and format errors
var err = errors.new("Something went wrong")
var formatted = errors.fmt("Error: $ at $", "failure", time:now())

# Check errors
if errors.is(result) {
    io.println("Got an error: " + result)
}

# Create coroutines with refl(), returns a promise
var p = refl(fun() {
    time.sleep(1000)
    return "done"
}).then(fun(res) {
    io.println("Success:", res)
}).catch(fun(e) {
    io.println("Failure:", e)
}).finally(fun() {
    io.println("Coroutine completed")
})

# Promises are cancellable
p.cancel() 

# Register event handlers, they can be triggered from go
events.register("click", fun(evt, x, y) {
    io.println("Clicked at", x, y)
})

# Schedule tasks
events.schedule(fun() {
    io.println("Scheduled task")
}, time:now() + 5000)

# Panic (unrecoverable error)
errors.panic("Fatal error")

# Eval code
eval(`"6"+7`) # "67"
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
    "refl/runtime/objects"
)

// Example usage
func main() {
	p := parser.New()
	program, err := p.Parse(`io.println("Hello world!")`)
	if err != nil {
		panic(err)
	}

	env := runtime.NewEnvironment(nil)
	env.Define("my_func", objects.NewWrapperFunction(func(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
      return objects.NewString("Hello world!"), nil
	})) // can define custom global variables

	evaluator := eval.New(context.Background(), program, env)
	evaluator.FireEvent("my_event", []runtime.Object{objects.NewString("my_data")}) // can fire events from go
	
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
* `refl` - spawns a coroutine
* `eval` - evaluates a refl code
