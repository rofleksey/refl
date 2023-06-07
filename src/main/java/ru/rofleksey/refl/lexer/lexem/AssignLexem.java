package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class AssignLexem implements Lexem {
    public static final AssignLexem INSTANCE = new AssignLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.ASSIGN;
    }
}
