[
    [
        let x = 4 in x for x in [1, 2, 3],
        (let x = 4 in x) for x in [1, 2, 3]
    ],
    [
        let x = 5 in let x = 4 in x + x,
        let x = 5 in (let x = 4 in x) + x
    ],
    [
        if true then 1 else 2 + 1,
        (if true then 1 else 2) + 1
    ]
]