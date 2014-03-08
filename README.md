automata.go
===========

An overly convoluted implementation of 2D cellular automata using golang channels

This is an exercise in using channels in Go.

Things I'm still curious about:
- Can I get rid of `clk`?
- How do I ensure my go routines are running on more than one core? (I'm skeptical that they are doing so currently)

General idea
------------
![Structures](https://raw.github.com/uberj/automata.go/master/diagram.jpg)
(Click to see image rotated by pi)

F1Cell behavior:
- Block reading from clk
- Write state to out (blocking?)
- Write to next\_state (blocking)
- Read from appropriate p cells to calculate new state
- loop

PCell behavior:
- Read next state (blocking)
- Write correct number of reads to state channel
- loop



Examples
```
$ go run automata.go -help
  -N=32: Number of nodes
  -gen=30: Number of generations to iterate
  -rule=30: Rule to be used for calculating state: i.e. rule 90 http://en.wikipedia.org/wiki/Rule_90
  -seed=65536: A seed value (powers of 2 are nice)
```

(Default) rule 30
```
$ go run automata.go
Rule: 30
Seed (decimal): 65536
Generations: 30
State array: 0 1 1 1 1 0 0 0

Set seed to: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0


State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 1 1 0 1 1 1 1 0 1 1 1 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 1 1 0 0 1 0 0 0 0 1 0 0 1 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 1 1 0 1 1 1 1 0 0 1 1 1 1 1 1 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 1 1 0 0 1 0 0 0 1 1 1 0 0 0 0 0 1 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 1 1 0 1 1 1 1 0 1 1 0 0 1 0 0 0 1 1 1 0 0 0 0 0 0
State array: 0 0 0 0 0 0 1 1 0 0 1 0 0 0 0 1 0 1 1 1 1 0 1 1 0 0 1 0 0 0 0 0
State array: 0 0 0 0 0 1 1 0 1 1 1 1 0 0 1 1 0 1 0 0 0 0 1 0 1 1 1 1 0 0 0 0
State array: 0 0 0 0 1 1 0 0 1 0 0 0 1 1 1 0 0 1 1 0 0 1 1 0 1 0 0 0 1 0 0 0
State array: 0 0 0 1 1 0 1 1 1 1 0 1 1 0 0 1 1 1 0 1 1 1 0 0 1 1 0 1 1 1 0 0
State array: 0 0 1 1 0 0 1 0 0 0 0 1 0 1 1 1 0 0 0 1 0 0 1 1 1 0 0 1 0 0 1 0
State array: 0 1 1 0 1 1 1 1 0 0 1 1 0 1 0 0 1 0 1 1 1 1 1 0 0 1 1 1 1 1 1 1
State array: 1 1 0 0 1 0 0 0 1 1 1 0 0 1 1 1 1 0 1 0 0 0 0 1 1 1 0 0 0 0 0 0
State array: 0 0 1 1 1 1 0 1 1 0 0 1 1 1 0 0 0 0 1 1 0 0 1 1 0 0 1 0 0 0 0 0
State array: 0 1 1 0 0 0 0 1 0 1 1 1 0 0 1 0 0 1 1 0 1 1 1 0 1 1 1 1 0 0 0 0
State array: 1 1 0 1 0 0 1 1 0 1 0 0 1 1 1 1 1 1 0 0 1 0 0 0 1 0 0 0 1 0 0 0
State array: 0 0 0 1 1 1 1 0 0 1 1 1 1 0 0 0 0 0 1 1 1 1 0 1 1 1 0 1 1 1 0 0
State array: 0 0 1 1 0 0 0 1 1 1 0 0 0 1 0 0 0 1 1 0 0 0 0 1 0 0 0 1 0 0 1 0
State array: 0 1 1 0 1 0 1 1 0 0 1 0 1 1 1 0 1 1 0 1 0 0 1 1 1 0 1 1 1 1 1 1
State array: 1 1 0 0 1 0 1 0 1 1 1 0 1 0 0 0 1 0 0 1 1 1 1 0 0 0 1 0 0 0 0 0
State array: 0 0 1 1 1 0 1 0 1 0 0 0 1 1 0 1 1 1 1 1 0 0 0 1 0 1 1 1 0 0 0 0
State array: 0 1 1 0 0 0 1 0 1 1 0 1 1 0 0 1 0 0 0 0 1 0 1 1 0 1 0 0 1 0 0 0
State array: 1 1 0 1 0 1 1 0 1 0 0 1 0 1 1 1 1 0 0 1 1 0 1 0 0 1 1 1 1 1 0 0
State array: 0 0 0 1 0 1 0 0 1 1 1 1 0 1 0 0 0 1 1 1 0 0 1 1 1 1 0 0 0 0 1 0
State array: 0 0 1 1 0 1 1 1 1 0 0 0 0 1 1 0 1 1 0 0 1 1 1 0 0 0 1 0 0 1 1 1
State array: 0 1 1 0 0 1 0 0 0 1 0 0 1 1 0 0 1 0 1 1 1 0 0 1 0 1 1 1 1 1 0 0
```

rule 110
```
$ go run automata.go -rule=110
Rule: 110
Seed (decimal): 65536
Generations: 30
State array: 0 1 1 1 0 1 1 0

Set seed to: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0


State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 0 1 1 1 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 0 1 1 0 1 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 0 1 1 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 1 1 1 0 0 0 0 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 1 1 0 1 0 0 0 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 1 1 1 1 1 0 0 1 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 1 1 0 0 0 1 0 1 1 1 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 1 1 1 0 0 1 1 1 1 0 1 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 1 1 0 1 0 1 1 0 0 1 1 1 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 1 1 1 1 1 1 1 1 0 1 1 0 0 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 1 0 0 0 0 0 0 1 1 1 1 0 0 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 0 1 1 0 0 1 0 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 0 1 1 1 0 1 1 1 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 1 1 0 1 1 1 0 0 1 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 1 1 1 1 1 0 1 0 1 1 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 1 1 0 0 0 1 1 1 1 1 1 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 1 1 1 0 0 1 1 0 0 0 0 0 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 1 1 0 1 0 1 1 1 0 0 0 0 0 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 1 1 1 1 1 1 0 1 0 0 0 0 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 1 0 0 0 0 1 1 1 0 0 0 1 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 0 1 1 0 1 0 0 1 1 1 0 0 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 0 1 1 1 1 1 0 1 1 0 1 0 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
State array: 0 0 1 1 0 0 0 1 1 1 1 1 1 1 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
```
