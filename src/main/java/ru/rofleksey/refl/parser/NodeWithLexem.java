package ru.rofleksey.refl.parser;

import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lexer.LexemType;

public final class NodeWithLexem {
    private final Node node;
    private final LexemType lexemType;

    NodeWithLexem(Node node, LexemType lexemType) {
        this.node = node;
        this.lexemType = lexemType;
    }

    public Node node() {
        return node;
    }

    public LexemType lexemType() {
        return lexemType;
    }
}
