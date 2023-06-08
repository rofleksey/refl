package ru.rofleksey.refl.parser;


import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.ReflValue;
import ru.rofleksey.refl.lexer.Lexer;

import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ParserTest {
    private static Parser parser;
    private static Lexer lexer;

    @BeforeAll
    public static void setup() {
        parser = new Parser();
        lexer = new Lexer();
    }

    @Test
    public void test2Plus2() throws Throwable {
        var lexems = lexer.process("2 + 2;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(4, evalResult);
    }

    @Test
    public void testAddSequential() throws Throwable {
        var lexems = lexer.process("1 + 2 - 3;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(0, evalResult);
    }

    @Test
    public void testAddWithMult() throws Throwable {
        var lexems = lexer.process("1 + 2 * 3 - 4;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(3, evalResult);
    }

    @Test
    public void testBrackets() throws Throwable {
        var lexems = lexer.process("2 * (3 + 4);");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(14, evalResult);
    }

    @Test
    public void testMath() throws Throwable {
        var lexems = lexer.process("1 -1   + 2   - 2   +  4 - 4 +    6;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(6, evalResult);
    }

    @Test
    public void testMultiply() throws Throwable {
        var lexems = lexer.process("1 * 2 * 3 + 4 * 5 * 6 / 2 * 3;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(186, evalResult);
    }

    @Test
    public void testAnd() throws Throwable {
        var lexems = lexer.process("1 & 0 & 1 & 1;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(0, evalResult);
    }

    @Test
    public void testOr() throws Throwable {
        var lexems = lexer.process("0 | 0 | 1 | 0;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var evalResult = result.get(0).evaluate(new ReflContext()).asNumber().getValue();
        assertEquals(1, evalResult);
    }

    @Test
    public void testAddVar() throws Throwable {
        var lexems = lexer.process("x + y;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var ctx = new ReflContext();
        ctx.setVar("x", new NumberValue(2));
        ctx.setVar("y", new NumberValue(3));

        var evalResult = result.get(0).evaluate(ctx).asNumber().getValue();
        assertEquals(5, evalResult);
    }

    @Test
    public void testIfSimple() throws Throwable {
        var lexems = lexer.process("if 2 > x : 10; end;");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var ctx = new ReflContext();
        ctx.setVar("x", new NumberValue(1));

        var evalResult = result.get(0).evaluate(ctx).asNumber().getValue();
        assertEquals(10, evalResult);
    }

    @Test
    public void testCycle() throws Throwable {
        var lexems = lexer.process("x = 10; count = 0; while x > 0: if x / 2 == x / 2 * 2: count = count + 1; end; end; count;");
        var result = parser.parse(lexems);
        assertEquals(4, result.size());

        var ctx = new ReflContext();

        Value evalResult = ReflValue.INSTANCE;
        for (var res : result) {
            evalResult = res.evaluate(ctx);
        }

        assertEquals(0, evalResult.asNumber().getValue());
    }

    @Test
    public void testIf() throws Throwable {
        var lexems = lexer.process("x = 1; y = 2; if x < y : z = 10; end; if x > y : z = 100; end;");
        var result = parser.parse(lexems);
        assertEquals(4, result.size());

        var ctx = new ReflContext();
        result.get(0).evaluate(ctx);
        result.get(1).evaluate(ctx);
        result.get(2).evaluate(ctx);

        var evalResult = result.get(3).evaluate(ctx).asNumber().getValue();
        assertEquals(0, evalResult);
        assertEquals(10, ctx.getVar("z").asNumber().getValue());
    }

    @Test
    public void testExternalFunc() throws Throwable {
        var lexems = lexer.process("print('hello world');");
        var result = parser.parse(lexems);
        assertEquals(1, result.size());

        var reference = new AtomicReference<String>();

        var ctx = new ReflContext();
        ctx.setVar("print", new FunctionValue("print") {
            @Override
            public  Value call(ReflContext ctx, List<Value> args) throws EvalError {
                reference.set(args.get(0).toString());
                return ReflValue.INSTANCE;
            }
        });

        var evalResult = result.get(0).evaluate(ctx).asNumber().getValue();
        assertEquals(0, evalResult);
        assertEquals("hello world", reference.get());
    }

}