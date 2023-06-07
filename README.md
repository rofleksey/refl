<h1 align="center">refl</h1>
<p align="center">
      Simple embeddable scripting language for Java 11+<br><br>
</p>

--------

## How To Use

Declare functions and variables via context and run script.

```java
var ctx = new ReflContext();

ctx.setVar("x", new NumberValue(5));
ctx.setVar("print", new FunctionValue("print") {
  @Override
  public Value call(ReflContext ctx, List<Value> args) {
    args.forEach(it -> {
    	System.out.println(it.toString());
    });
    return ReflValue.INSTANCE;
  }
});

var result = Refl.eval(ctx, "while x > 0: print(x); x = x - 1; end;");

```

## Language

* operators: +, -, *, /, =, ==, <, >, &, |, !
* keywords: if, while, end
* types: number (double), string, refl (same as null), function (can be declared only via `ReflContext`)
* no support for global variables
* no support for function declarations inside code
* no std
* logical operators return 1 or 0

Grammar:

```
declList -> declList decl ; | decl ;
decl -> if and : declList end | while and : declList end | var = and | and
and -> orExp | and & orExp
orExp -> notExp | orExp or notExp
notExp -> rel | ! rel
rel -> add | add < add
add -> mul | add + mul | add - mul
mul -> unary | mul * unary
unary -> term | - term
term -> const | var | ( and ) | call
call -> var ( args )
args -> argsList | ϵ
argsList -> argsList , var | argsList , const | var | const
```
