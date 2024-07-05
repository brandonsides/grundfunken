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
            Exp(x - 1).fmap(func(x) 2 * x),
    
    Box = func(x) {
        Value: x,
        fmap: func(f) Box(f(x)),
        flatten: func() x,
        str: utils.concatAll(["Box(", toString(x), ")"]),
        andThen: bind(this)
    },

    IO = func(x) {
        str: "IO(Some(" + toString(x) + "))",
        flatten: func() x,
        fmap: func(f) IO(f(x)),
        andThen: bind(this)
    }

    putStr = func(x) {
        str: "IO(None)",
        
    }

    takeStr = func() {
        runIO: func() {
            let x = readLine()
            in Some(x)
        }
    },
in [
    Some(1).andThen(GetHalf).andThen(Exp).str,
    Some(2).andThen(GetHalf).andThen(Exp).str,
    Some(-2).andThen(GetHalf).andThen(Exp).str,

    Box(1).andThen(TimesTwo).andThen(AddOne).str
]