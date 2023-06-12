package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.StringValue;

public final class GetVarNode implements Node {
    private final StringValue name;

    public GetVarNode(StringValue name) {
        this.name = name;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return ctx.getVar(name);
    }

    @Override
    public Value getLeftSide(ReflContext ctx) {
        return ctx;
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return name;
    }

    @Override
    public String toString() {
        return name.toString();
    }
}
