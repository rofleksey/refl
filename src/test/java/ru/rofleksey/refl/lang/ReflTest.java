package ru.rofleksey.refl.lang;

import org.junit.jupiter.api.Test;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ParseError;
import ru.rofleksey.refl.lang.error.VarUndefinedError;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class ReflTest {
    @Test
    public void testIntLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("2");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(2.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testFloatLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("3.14");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(3.14, resultValue.asNumber().getValue());
    }

    @Test
    public void testStringLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("'hello world'");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals("hello world", resultValue.toString());
    }

    @Test
    public void testEscapedStringLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("\"\\\"\"");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals("\"", resultValue.toString());
    }

    @Test
    public void testEmptyStringLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("\"\"");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals("", resultValue.toString());
    }

    @Test
    public void testNilLiteral() throws EvalError, ParseError {
        var compiled = Refl.compile("nil");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(NilValue.INSTANCE, resultValue);
    }

    @Test
    public void test2plus2() throws EvalError, ParseError {
        var compiled = Refl.compile("2+2");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(4.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testAddSub() throws EvalError, ParseError {
        var compiled = Refl.compile("2 - 2 + 3 - 3 + 4 - 4 + 5");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(5.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testMulDiv() throws EvalError, ParseError {
        var compiled = Refl.compile("5 * 6 / 2 % 8");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(7.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testBrackets() throws EvalError, ParseError {
        var compiled = Refl.compile("2 - 3 * 4 + 11");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testMath() throws EvalError, ParseError {
        var compiled = Refl.compile("2 - 3 * 4 + 11");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testAnd() throws EvalError, ParseError {
        var compiled = Refl.compile("1 & 0 & 1 & 1");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(0.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testOr() throws EvalError, ParseError {
        var compiled = Refl.compile("1 | 0 | 1 | 1");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testNot() throws EvalError, ParseError {
        var compiled = Refl.compile("!0");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testWeirdInput() {
        assertThrows(ParseError.class, () -> {
            Refl.compile("~@invalid");
        });
    }

    @Test
    public void testMultipleStatements() throws EvalError, ParseError {
        var compiled = Refl.compile("0\n1");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testAddVar() throws EvalError, ParseError {
        var compiled = Refl.compile("x+y");
        var ctx = ReflContext.empty();
        ctx.setVar("x", new NumberValue(2));
        ctx.setVar("y", new NumberValue(3));
        var resultValue = compiled.execute(ctx);
        assertEquals(5.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testIfSimple() throws EvalError, ParseError {
        var compiled = Refl.compile("if 2 > x \n 10 \n end");
        var ctx = ReflContext.empty();
        ctx.setVar("x", new NumberValue(1));
        var resultValue = compiled.execute(ctx);
        assertEquals(10.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testElif() throws EvalError, ParseError {
        var compiled = Refl.compile("if 0 \n 0 \n elif 1 \n 10 \n end");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(10.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testElvis() throws EvalError, ParseError {
        var compiled = Refl.compile("nil ?? 1");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testElse() throws EvalError, ParseError {
        var compiled = Refl.compile("if 0 \n 0 \n elif 0 \n 1 \n else \n 10 \n end");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(10.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testWhile() throws EvalError, ParseError {
        var compiled = Refl.compile("x = 10 \n count = 0 \n while x > 0 \n if x % 3 == 0 \n count++ \n end \n x-- \n end \n count");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(3.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testSimpleFunc() throws EvalError, ParseError {
        var compiled = Refl.compile("fun fact \n result = 1 \n i = 2 \n while i <= it \n result *= i \n i++ \n end \n result \n end \n fact(6)");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(720.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testNamedArgs() throws EvalError, ParseError {
        var compiled = Refl.compile("fun echo \n args.text + '!' \n end \n echo(text ~ 'hello world')");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals("hello world!", resultValue.toString());
    }

    @Test
    public void testMultipleArgs() throws EvalError, ParseError {
        var compiled = Refl.compile("fun add \n (args[0] + args[1]) * args.mult \n end \n add(4, 5, mult ~ 2)");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(18.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testFuncContext() throws EvalError, ParseError {
        var compiled = Refl.compile("fun fact \n result = 1 \n i = 2 \n while i <= it \n result *= i \n i++ \n end \n result \n end \n fact(6)");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(720.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testRecursiveFunction() throws EvalError, ParseError {
        var compiled = Refl.compile("fun fact \n if it <= 1 \n << 1 \n end \n it * fact(it - 1) \n end \n fact(6)");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(720.0, resultValue.asNumber().getValue());
    }

    @Test
    public void testIfGlobalScope() throws EvalError, ParseError {
        var compiled = Refl.compile("x = 1 \n if 1 \n x = 2 \n end \n x");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(2, resultValue.asNumber().getValue());
    }

    @Test
    public void testIfLocalScope() throws ParseError {
        var compiled = Refl.compile("if 1 \n x = 2 \n end \n x");
        assertThrows(VarUndefinedError.class, () -> {
            compiled.execute(ReflContext.empty());
        });
    }

    @Test
    public void testFuncGlobalScope() throws EvalError, ParseError {
        var compiled = Refl.compile("x = 1 \n fun test \n x = 2 \n end \n test() \n x");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(2, resultValue.asNumber().getValue());
    }

    @Test
    public void testFuncLocalScope() throws ParseError {
        var compiled = Refl.compile("fun test \n x = 2 \n end \n test() \n x");
        assertThrows(VarUndefinedError.class, () -> {
            compiled.execute(ReflContext.empty());
        });
    }

    @Test
    public void testFuncLocalIfScope() throws ParseError, EvalError {
        var compiled = Refl.compile("if 1 \n x = 1 \n fun test \n x = 2 \n end \n test() \n x \n end");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(1, resultValue.asNumber().getValue());
    }

    @Test
    public void testFuncGlobalIfScope() throws ParseError, EvalError {
        var compiled = Refl.compile("y = 1 \n if 1 \n fun test \n y = 2 \n end \n test() \n end \n y");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(2, resultValue.asNumber().getValue());
    }

    @Test
    public void testScope() throws ParseError, EvalError {
        var compiled = Refl.compile("pi = 3.14 \n scope Math \n fun getPi \n pi \n end \n end \n Math.getPi()");
        var resultValue = compiled.execute(ReflContext.empty());
        assertEquals(3.14, resultValue.asNumber().getValue());
    }
}