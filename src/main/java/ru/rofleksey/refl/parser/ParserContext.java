package ru.rofleksey.refl.parser;

import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;
import ru.rofleksey.refl.parser.error.InvalidLexemError;
import ru.rofleksey.refl.parser.error.ParserError;
import ru.rofleksey.refl.parser.error.UnexpectedEofError;

import java.util.LinkedList;
import java.util.List;

class ParserContext {
    private final LinkedList<Lexem> lexems;

    public ParserContext(List<Lexem> lexems) {
        this.lexems = new LinkedList<>(lexems);
    }

    public Lexem peek() throws ParserError {
        if (lexems.isEmpty()) {
            throw new UnexpectedEofError();
        }
        return lexems.getFirst();
    }

    public void take() throws ParserError {
        if (lexems.isEmpty()) {
            throw new UnexpectedEofError();
        }
        lexems.removeFirst();
    }

    public boolean lookUp(LexemType type, int ahead) {
        if (lexems.size() < ahead + 1) {
            return false;
        }
        var lexem = lexems.get(ahead);
        return lexem.type().equals(type);
    }

    public void consume(LexemType type) throws ParserError {
        if (lexems.isEmpty()) {
            throw new UnexpectedEofError();
        }
        var lexem = lexems.removeFirst();
        if (!lexem.type().equals(type)) {
            throw new InvalidLexemError(type, lexem.type());
        }
    }

    public boolean isEmpty() {
        return lexems.isEmpty();
    }

    public int remainingSize() {
        return lexems.size();
    }
}
