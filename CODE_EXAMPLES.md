# Refl code examples

## Fibonacci
```
fib = fun (n) {
    if n == 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    
    return fib(n-1) + fib(n-2)
}

fib(10)
```

## Counter
```
var Counter = {
    new: fun(self, initial) {
        var inst = clone(self)
        inst.value = initial
        return inst
    },
    inc: fun(self) {
        self.value = self.value + 1
        return self.value
    },
    get: fun(self) {
        return self.value
    }
}

var c = Counter:new(5)
c:inc()
c:get()
```

## Map reduction
```
map = fun (arr, fn) {
    var result = {}
    for i, val in arr {
        var idx = number(i)
        result[i] = fn(val)
    }
    return result
}

var arr = {1, 2, 3, 4}
var doubled = map(arr, fun(x) { x * 2 })
doubled[2]
```