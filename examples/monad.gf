let
    utils = import("utils.gf"),

    bind = func(monad) func(f)
        monad.fmap(f).flatten(),

    None = {
        isSome: false,
        flatten: func() this,
        fmap: func(f) this,
        str: "None",
        andThen: bind(this)
    },

    Some = func(x) {
            isSome: true,
            Value: x,
            flatten: func() x,
            fmap: func(f) Some(f(x)),
            str: utils.concatAll(["Some(", toString(x), ")"]),
            andThen: bind(this)
        },

    GetHalf = func(x)
        if x % 2 is 0 then
            Some(x / 2)
        else
            None,

    Exp = func(x)
        if x < 0 then
            None
        else if x is 0 then
            Some(1)
        else
            Exp(x - 1).fmap(func(x) 2 * x)
in [
    Some(1).andThen(GetHalf).andThen(Exp).str,
    Some(2).andThen(GetHalf).andThen(Exp).str,
    Some(-2).andThen(GetHalf).andThen(Exp).str
]