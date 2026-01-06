# Refl - scriptable interpretable language

## Requirements

* everything is an object (key/value pairs), except nil
* THERE ARE NO ";" IN THE LANGUAGE SYNTAX
* no arrays - we use objects with numerical keys (first index is 0)
* types: nil, number (double), object, function, string
* no booleans - 0 and nil are "false", anything else is "true", there are no false/true literals, we use 0 and 1
* numbers are passed by value
* everything else is passed by reference, e.g. assigning an object to another variable does not copy it
* strings are always inside "", support escaping (\")
* raw string literals can be put inside ``, no support for escaping inside it
* strings are immutable
* math operators: -, +, *, /, %, - (unary negate)
* assignment operators: = (only 1)
* comparison operators: <, >, <=, >=, ==, !=
* logical operators: ||, &&, ! (unary not)
* same operator precedence as in Javascript, a:foo has the same precedence as a.foo

## Syntax:

```
# Line comments only, no block comments

var x = 1 # define local variable
x = 1     # define global variable (or reassign already defined one)

# Objects
{}          # define an empty object
{a:1,"b":2} # object literal, only string literals as keys
{1,a,"b"}   # sugar for creating an "array", e.g. {0:1,1:a,2:"b"}

x[0]     # field access by key 0
x["key"] # field access by key "key"
x.key    # same as x["key"]
x[key]   # use variable "key" as a key, can be an expression

# Functions
foo = fun (a, b) { return a + b } # all functions are anonymous, always have a ( ) and a { }

foo(1,2)   # function call
x.foo(1,2) # get a field "foo" in x and call it
x:foo(1,2) # same as x.foo(x,1,2)

foo()     # unused arguments will be nil
foo(1, 2) # extra arguments will be discarded

foo = fun (a, f) { return f(a) } # functions are a first-class type, support closures

# Control flow

if x { # no brackets for conditions in if / while / for
  # if body
} elif y {
  # elif body
} else {
  # else body
}

while x {
  # while body
}

for key, value in obj { # scoped variables, always shadow
  # for body
}

break    # exits the inner-most while/for
continue # skips a single iteration of the inner-most while/for
```

## Type conversions
```
x # uninitialized variables are always nil, nil is a reserved word

nil + 1  # non-logical operators applied to nil always panic, 
nil == 1 # except for == and !=

c = (a = b) # syntax error, multiple assignment is not possible

"5" + 2 = "52"  # in string + X operations, second argument is always converted to string
"5" * 3 = "555" # string replication
5 + "2"         # panic, need to convert "2" to number first

{} + 5   # objects panic on math operations and comparison operations, except for == and !=
{} == {} # false, objects are only equal if they are a reference to the same object
```

## Standard library
```
$    # special variable, is an object containing all global variables
args # special variable, is set inside every function call, is an object containing all arguments

exit(x) # raises a panic (similar to golang)
type(x) # returns string that is a type of x
len(x) # returns a length of an "array"-like object, O(n)
str(x) # converts anything to string
number(x) # converts anything to number, returns nil on error
clone(x) # creates a deep copy of an object, functions are copied by reference, returns nil on error
```