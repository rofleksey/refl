package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public final class GetMemberNode implements Node {
    private final Node child;
    private final Value key;

    public GetMemberNode(Node child, Value key) {
        this.child = child;
        this.key = key;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).getVar(key);
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx);
    }

    @Override
    public Value getSetterKey(ReflContext ctx) {
        return key;
    }

    @Override
    public String toString() {
        return "(" + child.toString() + "." + key.toString() + ")";
    }
}
