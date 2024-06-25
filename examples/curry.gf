let
    curry = func(f, x) func(y) f(x, y),
    add = func(a, b) a + z
in
    curry(add, 5)(3)