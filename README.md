# Phonk

**Phonk** is a terrible, no-purpose, experimental programming language that is currently slated to
become the best, most-used, most widely critically- and popularly-acclaimed programming language of
all time by the year 2274.

# Build and Install

The Phonk interpro-compiler\* is currently not distributed in binary form, so it will need to be
compiled from source using the [Go compiler](https://go.dev/dl/), version 1.22 or higher.

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