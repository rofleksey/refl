package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public final class GetVarNode implements Node {
    private final String name;

    public GetVarNode(String name) {
        this.name = name;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        return ctx.getVar(name);
    }

    @Override
    public String toString() {
        return name;
    }
}
