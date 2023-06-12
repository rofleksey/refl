package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;

public final class UnaryPlusNode implements Node {
    private final Node child;

    public UnaryPlusNode(Node child) {
        this.child = child;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).asNumber();
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws NotReferencableError {
        throw new NotReferencableError(toString());
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return null;
    }

    @Override
    public String toString() {
        return "(+" + child.toString() + ")";
    }
}
