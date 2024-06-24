let a = 3,
    b = 4,
    funcIf = func (cond, res1, res2)
        if cond then
            res1 + a
        else
            res2 + b,
    a = 5,
    b = 6
in [
    funcIf(true, a, b),
    funcIf(false, a, b)
]