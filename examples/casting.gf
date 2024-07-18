let
    halfOf = func(i int) int | bool
        if i % 2 is 0 then
            i/2
        else
            false
in
    2 + halfOf(3) as int