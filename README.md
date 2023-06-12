<h1 align="center">refl</h1>
<p align="center">
      Simple embeddable scripting language for Java 11+ and Android<br><br>
</p>

--------

[![Java 11+](https://img.shields.io/badge/java-11-4c7e9f.svg)](http://java.oracle.com)
[![License](https://img.shields.io/badge/license-MIT-4c7e9f.svg)](https://raw.githubusercontent.com/rofleksey/refl/main/LICENSE.txt)
[![Maven Central](https://img.shields.io/maven-central/v/ru.rofleksey.refl/refl)](https://central.sonatype.com/artifact/ru.rofleksey.refl/refl)

## Installation

#### Gradle

```groovy
implementation 'ru.rofleksey.refl:refl:0.0.6'
```

#### Maven

```xml
<dependency>
    <groupId>ru.rofleksey.refl</groupId>
  <artifactId>refl</artifactId>
  <version>0.0.6</version>
</dependency>
```

## Usage

Declare functions and variables via context and run script.

```java
import ru.rofleksey.refl.lang.Refl;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ParseError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) throws EvalError, ParseError {
        var ctx = ReflContext.empty();

        ctx.setVar("x", new NumberValue(5));
        ctx.setVar("print", new FunctionValue("print") {
            @Override
            public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) {
                var prefix = namedArgs.get("prefix").toString();
                args.forEach(it -> System.out.println(prefix + it.toString()));
                return NilValue.INSTANCE;
            }
        });

        var executor = Refl.compile("while x > 0 \n print(x, prefix ~ '>') \n x-- \n end");
        var result = executor.execute(ReflContext.empty());
    }
}
```

## Language

* operators: +, -, *, /, %, =, +=, -=, *=, /=, %=, ==, !=, <, <=, >, >=, &, |, !, ??, <<
* keywords: if, while, end, fun, scope
* types: number (double), string, nil, function, object
* built-in functions:
  * exit (stops execution, returns first argument or nil)
  * wait (waits for `ReflContext.notifyCtx()` call)
  * sleep (calls Thread.sleep with first argument)
  * random, floor, ceil, round
  * string, number (conversion)
* logical operators return 1 or 0

## Syntax

```
# variable definitions
var = 1

# conditions
if condition1
  body1
  body2
elif condition2
  body1
  body2
else
  body1
  body2
end

# loops
while condition
  body1
  body2
end

# functions
fun fact
  result = 1
  i = 2

  while i <= it
    result *= i
    i++
  end

  result
end

fun concat
  if args.length == 0
    << nil
  end

  if args.length == 1
    << args[0]
  end

  result = it
  separator = args.separator ?? ''

  i = 1
  while i < args.length
    result += separator + args[i]
    i++
  end

  result
end

# scope (singleton struct)
scope Math
  PI = 3.14

  fun min
    result = it
    i = 1

    while i < args.length
      if args[i] < result
        result = args[i]
      fi
      i++
    end

    result
  end
  ...
end

pi = Math.PI
z = Math.min(x, y)
```