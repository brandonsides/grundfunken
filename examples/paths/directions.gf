let
    // directions
    LEFT = 0,
    UP = 1,
    RIGHT = 2,
    DOWN = 3
in {
    LEFT: LEFT,
    UP: UP,
    RIGHT: RIGHT,
    DOWN: DOWN,
    
    dirString: func(dir)
        if dir is DIR_LEFT then
            "L"
        else if dir is DIR_UP then
            "U"
        else if dir is DIR_RIGHT then
            "R"
        else if dir is DIR_DOWN then
            "D"
        else
            "?"
}