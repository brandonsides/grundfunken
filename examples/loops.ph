let
    list = [1, 2, 3, 4, 5, 6, 7, 8],
    factors = [10, 20, 30, 40]
in
    ((x * factor) for x in list) for factor in factors