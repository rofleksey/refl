package ru.rofleksey.refl.lexer;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import ru.rofleksey.refl.lexer.lexem.*;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

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
                new VarLiteral("x1"),
                AssignLexem.INSTANCE,
                new StringLexem("123"),
                new VarLiteral("abc"),
                AssignLexem.INSTANCE,
                new NumberLexem(123),
                new VarLiteral("y"),
                AssignLexem.INSTANCE,
                ReflLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testLargeNumber() {
        var e = assertThrows(LexerError.class, () -> {
            lexer.process("x = 12345678901234567890");
        });
        assertEquals("Failed to parse number '12345678901234567890' at 4", e.getMessage());
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
                new VarLiteral("x"),
                AssignLexem.INSTANCE,
                new StringLexem("12'3\\5"),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testMath() throws LexerError {
        var result = lexer.process("x+y*z/(a-c)");
        assertEquals(List.of(
                new VarLiteral("x"),
                PlusLexem.INSTANCE,
                new VarLiteral("y"),
                MultiplyLexem.INSTANCE,
                new VarLiteral("z"),
                DivideLexem.INSTANCE,
                BracketOpenLexem.INSTANCE,
                new VarLiteral("a"),
                MinusLexem.INSTANCE,
                new VarLiteral("c"),
                BracketCloseLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testLogic() throws LexerError {
        var result = lexer.process("x&y|z&(!a&!b)");
        assertEquals(List.of(
                new VarLiteral("x"),
                AndLexem.INSTANCE,
                new VarLiteral("y"),
                OrLexem.INSTANCE,
                new VarLiteral("z"),
                AndLexem.INSTANCE,
                BracketOpenLexem.INSTANCE,
                NotLexem.INSTANCE,
                new VarLiteral("a"),
                AndLexem.INSTANCE,
                NotLexem.INSTANCE,
                new VarLiteral("b"),
                BracketCloseLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testComparison() throws LexerError {
        var result = lexer.process("a < b & c > d & e == f");
        assertEquals(List.of(
                new VarLiteral("a"),
                LtLexem.INSTANCE,
                new VarLiteral("b"),
                AndLexem.INSTANCE,
                new VarLiteral("c"),
                GtLexem.INSTANCE,
                new VarLiteral("d"),
                AndLexem.INSTANCE,
                new VarLiteral("e"),
                EqLexem.INSTANCE,
                new VarLiteral("f"),
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testConditions() throws LexerError {
        var result = lexer.process("if a : b c end");
        assertEquals(List.of(
                IfLexem.INSTANCE,
                new VarLiteral("a"),
                ColonLexem.INSTANCE,
                new VarLiteral("b"),
                new VarLiteral("c"),
                EndLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }

    @Test
    public void testCycle() throws LexerError {
        var result = lexer.process("while a : b end");
        assertEquals(List.of(
                WhileLexem.INSTANCE,
                new VarLiteral("a"),
                ColonLexem.INSTANCE,
                new VarLiteral("b"),
                EndLexem.INSTANCE,
                EofLexem.INSTANCE
        ), result);
    }
}