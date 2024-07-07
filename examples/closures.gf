let a = 3,
    b = 4,
    funcIf = func (cond bool, res1, res2)
        if cond then
            res1
        else
            res2,
    a = 5,
    b = 6
in [
    funcIf(true, a, b),
    funcIf(false, a, b)
]