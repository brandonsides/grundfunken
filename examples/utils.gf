let
    // takes a list and returns everything after the first element
    tail = func(l)
        l[1:],
    
    // takes a list and returns all elements that meet the given condition
    filter = func(l, f)
        if len(l) is 0 then
            []
        else
            let
                this = l[0],
                rest = filter(tail(l), f)
            in
                if f(this) then
                    prepend(this, rest)
                else
                    rest,
in {
    tail: tail,

}