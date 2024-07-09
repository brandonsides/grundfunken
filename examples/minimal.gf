let
    f = func(m int) (int | bool)
        if m % 2 is 0 then
            m
        else
            false
in
    f(2)