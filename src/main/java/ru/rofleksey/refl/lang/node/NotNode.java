package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public final class NotNode implements Node {
    private final Node child;

    public NotNode(Node child) {
        this.child = child;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).not();
    }

    @Override
    public String toString() {
        return "(!" + child.toString() + ")";
    }
}
