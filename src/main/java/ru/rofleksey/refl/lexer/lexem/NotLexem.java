package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class NotLexem implements Lexem {
    public static final NotLexem INSTANCE = new NotLexem();

    @Override
    public  LexemType type() {
        return LexemType.NOT;
    }
}
