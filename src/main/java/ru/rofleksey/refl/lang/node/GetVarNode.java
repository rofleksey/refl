package ru.rofleksey.refl.lang.node;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public class GetVarNode implements Node {
    private final String name;

    public GetVarNode(String name) {
        this.name = name;
    }


    @Override
    public @NotNull Value evaluate(ReflContext ctx) throws EvalError {
        return ctx.getVar(name);
    }
}
