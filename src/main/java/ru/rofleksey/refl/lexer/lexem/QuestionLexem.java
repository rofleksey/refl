package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class QuestionLexem implements Lexem {
    public static final QuestionLexem INSTANCE = new QuestionLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.QUESTION;
    }
}
