<h1 align="center">refl</h1>
<p align="center">
      Simple embeddable scripting language for Java 11+<br><br>
</p>

--------

[![Java 11+](https://img.shields.io/badge/java-11-4c7e9f.svg)](http://java.oracle.com)
[![License](https://img.shields.io/badge/license-MIT-4c7e9f.svg)](https://raw.githubusercontent.com/rofleksey/refl/main/LICENSE.txt)
[![Maven Central](https://img.shields.io/maven-central/v/ru.rofleksey.refl/refl)](https://central.sonatype.com/artifact/ru.rofleksey.refl/refl)

## Installation

#### Gradle

```groovy
dependencies {
  implementation 'ru.rofleksey.refl:refl:0.0.2'
}
```

#### Maven

```xml
<dependency>
    <groupId>ru.rofleksey.refl</groupId>
    <artifactId>refl</artifactId>
    <version>0.0.1</version>
</dependency>
```

## Usage

Declare functions and variables via context and run script.

```java
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.ReflValue;
import ru.rofleksey.refl.lexer.LexerError;
import ru.rofleksey.refl.parser.error.ParserError;

import java.util.List;
import java.util.Map;

public class Main {
  public static void main(String[] args) throws ParserError, LexerError, EvalError {
    var ctx = new ReflContext();
    var compiler = new ReflCompiler();

    ctx.setVar("x", new NumberValue(5));
    ctx.setVar("print", new FunctionValue("print") {
      @Override
      public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) {
        var prefix = namedArgs.get("prefix").toString();
        args.forEach(it -> System.out.println(prefix + it.toString()));
        return ReflValue.INSTANCE;
      }
    });

    var executor = compiler.compile("while x > 0: print(x, prefix = '>'); x = x - 1; end;");
    var result = executor.execute(ctx);
  }
}
```

## Language

* operators: +, -, *, /, =, ==, <, >, &, |, !
* keywords: if, while, end
* types: number (double), string, refl (same as null), function (can be declared only via `ReflContext`)
* built-in functions:
  * exit (stops execution, returns first argument or refl)
  * wait (waits for `ReflContext.notifyCtx()` call)
  * sleep (calls Thread.sleep with first argument)
  * random, floor, ceil, round
  * string, number (conversion)

* all variables are global
* no support for function declarations inside code
* logical operators return 1 or 0
* function arguments can't be inlined expressions

Grammar:

```
declList -> declList decl ; | decl ;
decl -> if and : declList end | while and : declList end | var = and | and
and -> orExp | and & orExp
orExp -> notExp | orExp or notExp
notExp -> rel | not rel
rel -> add | add < add | add == add | add > add
add -> mul | add + mul | add - mul
mul -> unary | mul * unary | mul / unary
unary -> term | - term
term -> const | var | ( and ) | call
call -> var ( args )
args -> argsList | ϵ
argsList -> argsList , var | argsList , const | argsList , namedArg | var | const | namedArg
namedArg -> var = var | var = s
```
