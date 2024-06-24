let
    add = func(a) func(b) a + b,
    add5 = add(5) // func(b) 5 + b
in
    add5(3) // 8