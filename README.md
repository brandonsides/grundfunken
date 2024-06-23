# Phonk

**Phonk** is a terrible, no-purpose, experimental programming language that is currently slated to
become the best, most-used, most widely critically- and popularly-acclaimed programming language of
all time by the year 2274.

# Build and Install

The Phonk interpro-compiler\* is currently not distributed in binary form, so it will need to be
compiled from source using the [Go compiler](https://go.dev/dl/), version 1.22 or higher.

# The Basics

The most important unit of code in Phonk is the *expression*; an expression is simply a semantic
unit that can be *evaluated*; that is, resolved to a *value*.

All Phonk programs are ultimately expressions, constructed from smaller expressions connected
by the language's syntactic constructs.  Likewise, any expression is a valid Phonk program; the
simplest valid Phonk programs, then, are literals of the fundamental types, say `string`:

```
"hello world"
```

The above is a valid Phonk program that evaluates to the string "hello world"; saving it as
`hello-world.ph` and running it, we see that the interpreter outputs the string as the program's
result:

```
% ./phonk -input hello-world.ph
Result: hello world
```

# Types

The current fundamental types in Phonk are *integers*, *booleans*, *strings*, *arrays*, *objects*,
and *functions*.  Here are example literals of the various types:

- integer: `1`, `45`, `-16`
- boolean: `true`, `false`
- string: `"hello world"`, `"\"hello world\""`
- array: `[]`, `[1, 2, 3]` `[1, false, "hello"]`
- object: `{}`, `{hello: "hello", world: "world"}`
- functions: `func(x) x`, `func(a, b) a + b`

# Variables

The language has three ways of introducing a new variable: `let`, `for`, and `func`.

## Let

A `let` expression is the simplest way of binding a variable in Phonk.  It consists of a
*binding identifier*, a `let` clause, and an `in` clause, in the following form:

```swift
let x = 4 in x
```

In the above program, the binding identifier is `x`, the `let` clause is `4`, and the `in` clause is
`x`.  The program evaluates to `4`.  The value of the `let` clause - namely, `4` - is bound to
the identifier `x`, and then the `in` clause is evaluated given that binding.  In this case, the
`in` clause is just `x`, which evaluates to `4`.

The binding is only in scope when the `in` clause is being evaluated; it is "forgotten" after that
point.  Thus, the following program is invalid:

```swift
let 
    x = let
            y = 3
        in
            y + 1
in
    y + x
```

Running this, we get an error:

```
error at line 8, column 5: cannot evaluate unbound identifier
    y + x
    ^-here
```

In this case, we have two separate `let` expressions: the *inner* and *outer* expression.  The outer
expression has `x` as its binding identifier and `y + x` as its `in` clause.  Its `let` clause is
the inner `let` expression, which has `y` as its binding identifier, `3` as its `let` clause, and
`y + 1` as its `in` clause.

The inner `let` expression can be evaluated just fine.  The value `3` is bound to the identifier
`y`, and the `in` clause - `y + 1` - evaluates to `4`; thus, the entire inner `let` expression
evaluates to `4`.  Now that we are done evaluating this `in` clause, the binding `y = 3` falls out
of scope, and `y` becomes unbound.

The problem arises when evaluating the outer `let` expression.  The result of the inner
`let` expression - `4` - is bound to the outer binding identifier `x`.  Then, we attempt to evaluate
the outer expression's `in` clause.  However, this clause makes reference to the identifier `y`,
which was bound in the *inner* `let` expression.  Since this binding has fallen out of scope, the
identifier `y` is no longer bound to a value, resulting in the "unbound identifier" error we see
above.

A single `let` expression can support multiple bindings, and hence multiple `let` clauses:

```swift
    let x = 3,
        y = 4
    in
        x + y // 7
```

Each binding will be present when evaluating subsequent `let` clauses even in the same let
expression:

```swift
    let x = 3,
        y = x + 1 // 4
    in
        x + y // 7
```

However, a `let` clause cannot reference a binding that will be bound *later* in the expression:

```swift
    let x = y - 1, // error at line 1, column 9: cannot evaluate unbound identifier
        y = 4
    in
        x + y
```

## For

A `for` expression consists of a `for` clause, a binding identifier, and an `in` clause can be used to manipulate an array
value-by-value.  It looks like this:

```swift
(i + 1) for i in [1,2,3]
```

This program evaluates to `[2 3 4]` (as a temporary quirk, arrays are represented differently
in a program's output than in the language itself).  Each element of the original array (given in
the `in` clause) is bound to the binding identifier `i`, and the `for` clause is evaluated; the returned
array contains the results of each evaluation.  So, in this case, the for clause is evaluated for
each value of the original array:

 - `1` is bound to `i`; `i+1` is evaluated as `2`
 - `2` is bound to `i`; `i+1` is evaluated as `3`
 - `3` is bound to `i`; `i+1` is evaluated as `4`

 The returned array contains all these results in the same order as the original array; hence,
 `[2, 3, 4]`.

As another temporary quirk, the `for` keyword has the highest precedence in the language, so it's
usually necessary to wrap the `for` clause in parentheses.  If they are omitted, unexpected results
can occur; if we try to evaluate the following program:

```swift
i + 1 for i in [1,2,3]
```

We get an error:

```
./phonk -input loop.ph
error at line 1, column 1: cannot evaluate unbound identifier
i + 1 for i in [1,2,3]
^-here
```

Because the `for` keyword takes precedence over `+`, this program is equivalent to

```swift
i + (1 for i in [1,2,3])
```

In this case, the `for` clause is just `1`.  The `for` expression only binds its binding identifer



# Roadmap

The 230-year roadmap for Phonk includes the following language features:

- purely functional: the Phonk language syntax will strictly require code to abide by the pure
functional paradigm.

- Dependently-typed: Phonk will support dependent types, allowing types to depend on runtime values.
Syntactically, types will behave like any other kind of object, and it will be possible to
parametrize them not only by other types (as one would do when creating a "generic" type, as has
become standard in modern programming languages), but also by runtime values.  Indeed, Phonk erodes
the boundary between types and values to the point that a developer may be hesitant to think
primarily in those terms: `type` is, itself, just another type an object can have.

For example, `List(elem type, len uint) type` will be a type constructor that takes a type
representing the type of the list's elements, and a uint representing the number of elements in the
list; it will return a new type.  A value of type `List(string, 7)` will be guaranteed to satisfy
the type returned by this constructor; namely, a list containing exactly seven values of type
`string`.

The `elem` parameter may be familiar to many programmers from languages that support generics; even
an integer type parameter may be familiar to users of some languages like C++.  However, the
distinctive feature of dependent types is that the value of `len` need not be fixed at compile-time;
for example, a function `FillList(elemType type, elem elemType, len uint)` might return a
`List(elemType, len)`, thus guaranteeing that the returned list will contain exactly `len` elements,
even though the value of `len` will not be known until runtime.

This will allow types to express nearly arbitrary facts about the code, allowing many aspects of a
program's correctness to be checked at compile-time.  Because reasoning about dependent types can be
complex and time-consuming for a compiler, the language will support *proof annotations* allowing
a programmer to supply manual proofs that a value or type necessarily satisfies another type.

- Paraconsistent: The Phonk interpro-compiler will support paraconsistency in type reasoning,
enabling types to be fully first-class objects while avoiding (or rather, embracing) certain
paradoxes such as:

    - Russell's paradox as applied to types; i.e., the language will be able to support a type of
    all types which do not instantiate themselves.  The answer to the question of whether this type
    instantiates itself will be, as far as the compiler is concerned, both yes *and* no.

    - The Burali-Forti paradox; i.e., the language will be able to support:
        - a type `ordinal` representing ordinal numbers
        - a function `order(s wellOrder) ordinal`, where wellOrder is a type representing a
        [well-order](https://en.wikipedia.org/wiki/Well-order), which returns the ordinal
        - a function `ordered(t type, cmp func(a t, b t) bool)` which returns a well-ordering of the
        elements of a given type, given a comparison function over elements of that type.

By 2550, Phonk is slated to support the following features:

- Compiled and Bootstrapped: Once the core language features are present, a compiler backend
will be written, thus transforming the interpro-compiler into, well, an actual compiler.  With a
minimal compiler in place, written in Go as the current interpro-compiler is, an equivalent compiler
will be written in Phonk and compiled, providing a self-compiling Phonk compiler that can be used to
continue developing Phonk *in* Phonk.

- Looser functionalism enforcement: The language will continue to encourage purely functional code
wherever possible, while allowing explicit violations of this principle where required.  Where
side-effects are necessary, the language will encourage them to be hidden from user code and
isolated to a single object or package.  Code that causes side-effects will be syntactically
conspicuous.

- Pseudo-imperativity: The language will support imperative syntax while retaining a functional
semantics where possible using native implicit monads representing "global state".  It will be
possible to make the same kinds of guarantees about the global state as any other value.

- Natively concurrent: Phonk will support concurrent loops and subroutine calls natively for
code that abides by the pure-functional paradigm.  More minimal compile-time concurrency features
will be available to code that manipulates state.  Code that is marked as state-manipulating will
have the ability to more directly manage concurrency, allowing optimization where the language's
native concurrency features are insufficient.

^* an *interpro-compiler* is a technical industry term for a interpreter that will
hopefully one day be a compiler