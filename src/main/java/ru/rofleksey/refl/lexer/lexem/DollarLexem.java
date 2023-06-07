package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class DollarLexem implements Lexem {
    public static final DollarLexem INSTANCE = new DollarLexem();

    @Override
    public  LexemType type() {
        return LexemType.DOLLAR;
    }
}
