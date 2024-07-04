let
    beginningHands = let
        withLeftHand = func(hands, newLeftHand) {
            split: this.split,
            receiveAttack: this.receiveAttack,
            attack: this.attack,
            left: newLeftHand,
            right: this.right
        },

        withRightHand = func(newRightHand) {
            split: this.split,
            receiveAttack: this.receiveAttack,
            attack: this.attack,
            left: this.left,
            right: newRightHand
        }
    in {
        split: func()
            if this.left.numFingers is not 0 and this.right.numFingers is not 0 then
                // must have no fingers left on one hand
                false
            else
                let newNumFingers = if this.left.numFingers is 0 then
                            this.right.numFingers / 2
                        else
                            this.left.numFingers / 2
                in {
                    split: this.split,
                    left: {
                        numFingers: newNumFingers
                        attack: this.left.attack
                    },
                    right: {
                        numFingers: newNumFingers
                        attack: this.right.attack
                    }
                },

        receiveAttack: func(fromHand, targetHand)
            if targetHand is HAND_LEFT then {
                left: this.left.receiveAttack(fromHand),
                right: this.right
            } else {
                left: this.left,
                right: this.right.receiveAttack(fromHand)
            }

        attack: func(withHand, targetHand)
            if withHand is HAND_LEFT then
                this.left.attack(target)
            else
                this.right.attack(target)

        left: {
            numFingers: 1,
            attack: func(hand) {
                numFingers: (hand.numFingers + this.numFingers) % 5,
                attack: hand.attack
            },
            receiveAttack: func(fromHand) {
                numFingers: (this.numFingers + fromHand.numFingers) % 5,
                attack: this.attack
            }
        },

        right: {
            numFingers: 1,
            attack: func(hand) {
                numFingers: (hand.numFingers + this.numFingers) % 5,
                attack: hand.attack
            }
        }
    },

    TURN_PLAYER1 = true,
    TURN_PLAYER2 = false,

    HAND_LEFT = true,
    HAND_RIGHT = false,
    
    beginningGame = {
        turn: TURN_PLAYER1,
        player1: beginningHands,
        player2: beginningHands,
        attack: func(withHand, target) {
            turn: not this.turn,
            player1: if this.turn is TURN_PLAYER1 then
                    this.player1
                else {
                    split: this.player1.split,
                    left: if target is HAND_RIGHT then
                            this.player2.attack(withHand, this.player1.left)

                    right: this.player1.right.attack(withHand)
                },
        }
    } in 