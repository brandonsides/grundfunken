let
    curry(f, x) = func(y) f(x, y),
    add = func(a, b) a + b,
    add5 = curry(add, 5)
in
    add5(3) // 8