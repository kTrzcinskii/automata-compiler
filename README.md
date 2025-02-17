# Automata Compiler

_Work in progress_

## Supported Automata

### Turing Machine (Standard Model)

A Turing Machine uses a tape and a state to determine its next move. The model implemented in this application assumes that the tape is **one-way infinite**, meaning you can move indefinitely to the right. However, moving left beyond the starting position results in an error.

If no transition is defined for a given state and symbol, the program terminates with an error. The final tape output is trimmed of any trailing `B` (blank symbols) except for one. For example, if the final tape is `S1|S1|S2|B|B|B`, the program returns `S1|S1|S2|B`.

#### Input Format

```
q0 q1 ... qn; [states]
qs; [initial state]
qf1 qf2 ... qfk; [accepting states]
a1 a2 ... an; [symbols]

(q, s) > (new_q, new_s, move)
(q, s) > (new_q, new_s, move)
(q, s) > (new_q, new_s, move)
...;

a1 a1 a3 a8 ...; [initial tape]
```

#### Rules and Conventions
- Each state must start with the letter `q`, followed by one or more alphanumeric characters.
- `move` can be either:
  - `L` (left)
  - `R` (right)
  - These symbols (`L` and `R`) are **reserved** and cannot be used for anything else.
- The Turing Machine starts in the `initial state`, pointing at the **first element** of the `initial tape`.
- Each section must be **terminated by a semicolon** (`;`).
- `B` is a **reserved symbol** representing a blank space on the tape.

#### Examples

You can find example Turing Machine programs in the [examples/turing-machine](examples/turing-machine) directory.

