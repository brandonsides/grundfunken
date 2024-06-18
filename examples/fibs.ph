let
    // functions
    tail = func(l)
        slice(l, 1, -1),

    all = func(f, l)
        if equals(len(l), 0) then
            true
        else
            let this = at(l, 0) in
                if f(this) then
                    all(f, tail(l))
                else
                    false,

    takeWhile = func(f, l)
        if equals(len(l), 0) then
            []
        else
            let this = at(l, 0) in
                if f(this) then
                    prepend(this, takeWhile(f, tail(l)))
                else
                    [],

    isPrime = func(n)
        if lessThan(n, 2) then
            false
        else if equals(n, 2) then
            true
        else all(
            func(x) not(equals(0, mod(n, x))),
            takeWhile(
                func(x) not(greaterThan(x * x, n)),
                range(2, n - 1)
            )
        ),

    filter = func(f, l)
        if equals(len(l), 0) then
            []
        else
            let
                this = at(l, 0),
                rest = filter(f, tail(l))
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,

    fibHelper = func(n, a, b)
        if equals(n, 0) then
            a
        else
            fibHelper(n - 1, b, a + b),

    fib = func(n) fibHelper(n, 1, 1),

    fibUpTo = func(lim) fib(x) for x in range(0, lim)
in let
    // variables
    lim = 30,
    fibs = fibUpTo(lim)
in [
        "fibs:", fibs,
        "prime fibs:", filter(isPrime, fibs)
]