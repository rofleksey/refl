package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;
import ru.rofleksey.refl.lang.value.NilValue;

public final class WhileNode implements Node {
    private final Node condition;
    private final Node body;

    public WhileNode(Node condition, Node body) {
        this.condition = condition;
        this.body = body;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        Value result = NilValue.INSTANCE;

        while (condition.evaluate(ctx).isTruthy()) {
            var newCtx = ctx.createChild();
            result = body.evaluate(newCtx);
        }

        return result;
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
        return "while " + condition.toString() + "\n" + body.toString() + "\n";
    }
}
