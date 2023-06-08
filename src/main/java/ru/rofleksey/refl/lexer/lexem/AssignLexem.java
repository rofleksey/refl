package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class AssignLexem implements Lexem {
    public static final AssignLexem INSTANCE = new AssignLexem();

    @Override
    public LexemType type() {
        return LexemType.ASSIGN;
    }
}
