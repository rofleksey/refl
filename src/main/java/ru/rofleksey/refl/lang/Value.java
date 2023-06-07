package ru.rofleksey.refl.lang;


import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.StringValue;

import java.util.List;

public interface Value {
    
    Value add(Value other) throws EvalError;
    
    Value subtract(Value other) throws EvalError;
    
    Value multiply(Value other) throws EvalError;
    
    Value divide(Value other) throws EvalError;
    
    Value and(Value other) throws EvalError;
    
    Value or(Value other) throws EvalError;
    
    Value compare(Value other) throws EvalError;
    
    Value not() throws EvalError;
    
    Value call(ReflContext ctx, List<Value> args) throws EvalError;
    boolean isTruthy();
    
    StringValue asString();
    
    NumberValue asNumber();
    
    String getType();
}
