package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public final class ArrayCallNode implements Node {
    private final Node child;
    private final Node arg;

    public ArrayCallNode(Node child, Node arg) {
        this.child = child;
        this.arg = arg;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx).getVar(arg.evaluate(ctx));
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws EvalError {
        return child.evaluate(ctx);
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return arg.evaluate(ctx);
    }

    @Override
    public String toString() {
        return "(" + child.toString() + ")[" + arg.toString() + "]";
    }
}
