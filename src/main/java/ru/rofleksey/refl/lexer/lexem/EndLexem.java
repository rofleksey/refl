package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class EndLexem implements Lexem {
    public static final EndLexem INSTANCE = new EndLexem();

    @Override
    public  LexemType type() {
        return LexemType.END;
    }
}
