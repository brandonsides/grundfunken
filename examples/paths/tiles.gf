let types = {
    EMPTY: 0,
    WALL: 1,
    START: 2,
    END: 3
} in {
    types: types,

    toString: func(tile)
        if tile is types.EMPTY then
            " "
        else if tile is types.WALL then
            "#"
        else if tile is types.START then
            "S"
        else if tile is types.END then
            "E"
        else
            "?"
}