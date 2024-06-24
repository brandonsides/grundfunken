let
    compose = func(f, g) func(x) f(g(x)),

    listify = func(x) [x],

    add1 = func(x) x + 1
in
    compose(listify, add1)(3)