package ru.rofleksey.refl.lang;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.StringValue;

import java.util.List;

public interface Value {
    @NotNull
    Value add(Value other) throws EvalError;
    @NotNull
    Value subtract(Value other) throws EvalError;
    @NotNull
    Value multiply(Value other) throws EvalError;
    @NotNull
    Value divide(Value other) throws EvalError;
    @NotNull
    Value and(Value other) throws EvalError;
    @NotNull
    Value or(Value other) throws EvalError;
    @NotNull
    Value compare(Value other) throws EvalError;
    @NotNull
    Value not() throws EvalError;
    @NotNull
    Value call(ReflContext ctx, List<Value> args) throws EvalError;
    boolean isTruthy();
    @NotNull
    StringValue asString();
    @NotNull
    NumberValue asNumber();
    @NotNull
    String getType();
}
