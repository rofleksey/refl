package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class BracketCloseLexem implements Lexem {
    public static final BracketCloseLexem INSTANCE = new BracketCloseLexem();

    @Override
    public LexemType type() {
        return LexemType.BRACKET_CLOSE;
    }
}
