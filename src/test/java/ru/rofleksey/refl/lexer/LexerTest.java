package ru.rofleksey.refl.lexer;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import ru.rofleksey.refl.lexer.lexem.*;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class LexerTest {
    private static Lexer lexer;

    @BeforeAll
    public static void setup() {
        lexer = new Lexer();
    }

    @Test
    public void testEmpty() throws LexerError {
        var result = lexer.process("");
        assertEquals(List.of(EofLexem.INSTANCE), result);
    }

    @Test
    public void testBlank() throws LexerError {
        var result = lexer.process(" \n ");
        assertEquals(List.of(EofLexem.INSTANCE), result);
    }

    @Test
    public void testSimple() throws LexerError {
        var result = lexer.process("2+2");
        assertEquals(List.of(
                new NumberLexem(2),
                PlusLexem.INSTANCE,
                new NumberLexem(2),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testSpace() throws LexerError {
        var result = lexer.process("  2 \n \r + \r\n \n\r \t 2 ");
        assertEquals(List.of(
                new NumberLexem(2),
                PlusLexem.INSTANCE,
                new NumberLexem(2),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testVars() throws LexerError {
        var result = lexer.process("x1 = '123'\nabc = 123 y=refl");
        assertEquals(List.of(
                new VarLexem("x1"),
                AssignLexem.INSTANCE,
                new StringLexem("123"),
                new VarLexem("abc"),
                AssignLexem.INSTANCE,
                new NumberLexem(123),
                new VarLexem("y"),
                AssignLexem.INSTANCE,
                ReflLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testDoubleValue() throws LexerError {
        var result = lexer.process("0.25");
        assertEquals(List.of(
                new NumberLexem(0.25),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testUnclosedQuote() {
        var e = assertThrows(LexerError.class, () -> {
            lexer.process("x = ' test");
        });
        assertEquals("Closing quote not found, begins at 4", e.getMessage());
    }

    @Test
    public void testUnexpected() {
        var e = assertThrows(LexerError.class, () -> {
            lexer.process("x = @{");
        });
        assertEquals("Unexpected symbol '@' at 4", e.getMessage());
    }

    @Test
    public void testStringEscape() throws LexerError {
        var result = lexer.process("x = '12\\\'3\\\\5'");
        assertEquals(List.of(
                new VarLexem("x"),
                AssignLexem.INSTANCE,
                new StringLexem("12'3\\5"),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testMath() throws LexerError {
        var result = lexer.process("x+y*z/(a-c)");
        assertEquals(List.of(
                new VarLexem("x"),
                PlusLexem.INSTANCE,
                new VarLexem("y"),
                MultiplyLexem.INSTANCE,
                new VarLexem("z"),
                DivideLexem.INSTANCE,
                BracketOpenLexem.INSTANCE,
                new VarLexem("a"),
                MinusLexem.INSTANCE,
                new VarLexem("c"),
                BracketCloseLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testLogic() throws LexerError {
        var result = lexer.process("x&y|z&(!a&!b)");
        assertEquals(List.of(
                new VarLexem("x"),
                AndLexem.INSTANCE,
                new VarLexem("y"),
                OrLexem.INSTANCE,
                new VarLexem("z"),
                AndLexem.INSTANCE,
                BracketOpenLexem.INSTANCE,
                NotLexem.INSTANCE,
                new VarLexem("a"),
                AndLexem.INSTANCE,
                NotLexem.INSTANCE,
                new VarLexem("b"),
                BracketCloseLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testComparison() throws LexerError {
        var result = lexer.process("a < b & c > d & e == f");
        assertEquals(List.of(
                new VarLexem("a"),
                LtLexem.INSTANCE,
                new VarLexem("b"),
                AndLexem.INSTANCE,
                new VarLexem("c"),
                GtLexem.INSTANCE,
                new VarLexem("d"),
                AndLexem.INSTANCE,
                new VarLexem("e"),
                EqLexem.INSTANCE,
                new VarLexem("f"),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testConditions() throws LexerError {
        var result = lexer.process("if a : b c end");
        assertEquals(List.of(
                IfLexem.INSTANCE,
                new VarLexem("a"),
                ColonLexem.INSTANCE,
                new VarLexem("b"),
                new VarLexem("c"),
                EndLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testCycle() throws LexerError {
        var result = lexer.process("while a : b end");
        assertEquals(List.of(
                WhileLexem.INSTANCE,
                new VarLexem("a"),
                ColonLexem.INSTANCE,
                new VarLexem("b"),
                EndLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testFuncCall() throws LexerError {
        var result = lexer.process("foo (a, 5);");
        assertEquals(List.of(
                new VarLexem("foo"),
                BracketOpenLexem.INSTANCE,
                new VarLexem("a"),
                CommaLexem.INSTANCE,
                new NumberLexem(5),
                BracketCloseLexem.INSTANCE,
                SemicolonLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }
}