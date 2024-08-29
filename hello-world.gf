let
    // returns the nth fibonacci number, or unit if n is negative
    fib = func(n int) int | unit
        if n < 0 then
            unit
        else if n < 2 then
            n
        else
            // 
            match res on [fib(n - 1), fib(n - 2)]
            case [int]
                res[0] + res[1]
            case any
                unit
in
    fib(x) for x in range(-10, 10)