package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.value.ObjectValue;
import ru.rofleksey.refl.lang.value.StringValue;

public final class ScopeDeclarationNode implements Node {
    private final StringValue key;
    private final Node body;

    public ScopeDeclarationNode(StringValue key, Node body) {
        this.key = key;
        this.body = body;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var scope = new ObjectValue();
        var newCtx = ReflContext.startScope(scope);
        var result = body.evaluate(newCtx);

        ctx.setVar(key, scope);
        return result;
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws EvalError {
        throw new NoObjectContextError(toString());
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return null;
    }

    @Override
    public String toString() {
        return "scope " + key.toString();
    }
}
