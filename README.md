# Refl

A simple scripting language interpreter written in Go.

## Overview

Refl is an interpreted language where everything is an object (except `nil`). The language features:

- **Simple Types**: `nil`, `number`, `string`, `object`, `function`
- **Objects everywhere**: All values are objects with key/value pairs
- **No Booleans**: `0`, `nil` and `""` are "false", anything else is "true"
- **Dynamic Typing**: No type declarations needed
- **First-Class Functions**: Functions are values that support closures

## Syntax Examples

```javascript
// Variables
var x = 1
x = 2  // global assignment

// Objects
var obj = {a: 1, "b": 2}
var arr = {1, "two", 3}  // array syntax sugar

// Functions
var add = fun(a, b) { return a + b }
var result = add(1, 2)

// Control flow
if x > 0 {
    // do something
} elif x == 0 {
    // something else
} else {
    // default
}

// Loops
while condition {
    // loop body
}

for key, value in obj {
    // iterate over object or string
}
```