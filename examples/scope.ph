let
    nums = [1, 2, 3, 4],
    a = 1,
    b = a
in
    (let b = a+b in b) for a in nums