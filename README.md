# Automata Compiler

## Input format

### Turing Machine (standard model)
```
q0 q1 ... qn; [states]
qs; [initial state]
qf1 qf2 ... qfk; [accepting states]
a1 a2 ... an; [input symbols]
x1 x2 ... xn; [tape symbols]

(q, s) -> (new_q, new_s, move)
(q, s) -> (new_q, new_s, move)
(q, s) -> (new_q, new_s, move)
...;

a1a1a3a8...; [initial tape]

```

where:
- Elements in each line are separated by a single whitespace.
- Each state must start with the letter `q` followed by one or more alphanumeric characters.
- The line with tape symbols must not include both input symbols and the blank symbol.
- After the last line of `tape symbols`, there is one empty line, and then each line is a transition-function element.
- If for a given pair `(q, s)`, there is no transition defined, then computation terminates with an error.
- `move` is either `L` (left) or `R` (right) (so `L` and `R` are reserved symbols and cannot be used for other names).
- After transitions, there is one empty line, and then there is the `initial tape` state.
- The Turing machine starts with the `initial state` pointing to the first element of the `initial tape`.
- If a transition makes the Turing machine go outside of the tape, then computation terminates with an error.
- There should be a semicolon after each section.
- `B` is a reserved symbol (meaning blank symbol).
