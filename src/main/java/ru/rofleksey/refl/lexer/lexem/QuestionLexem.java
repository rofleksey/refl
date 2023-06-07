package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class QuestionLexem implements Lexem {
    public static final QuestionLexem INSTANCE = new QuestionLexem();

    @Override
    public  LexemType type() {
        return LexemType.QUESTION;
    }
}
