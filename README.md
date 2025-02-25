# Automata Compiler

Automata Compiler is a CLI program written in Go that allows you to simulate computations for different types of automata, including Deterministic Finite Automata (DFA), Pushdown Automata (PA), and Turing Machines (TM).

## Requirements

This application requires:

- Go `>= 1.24`
- Make `>= 4.0`

To run this, fist compile it:
```bash
make
```

Then, navigate to the build folder and run:

```bash
./automata-compiler AUTOMATON_TYPE INPUT_FILE [flags]
```
where:
- `AUTOMATON_TYPE` specifies the type of automaton:
   - DFA (deterministic finite automaton), 
   - PA (pushdown automaton), 
   - TM (turing machine)
- `INPUT_FILE` is the path to the file containing the automaton's source code. You can find example input files in the `examples` folder

## Supported Automata

### Turing Machine (Standard Model)

A **Turing Machine** uses a tape and a state to determine its next move. The model implemented in this application assumes that the tape is **one-way infinite**, meaning you can move indefinitely to the right. However, moving left beyond the starting position results in an error.

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
- Each symbol must constist of one or more alphanumeric characters.
- `move` can be either:
  - `L` (left)
  - `R` (right)
  - These symbols (`L` and `R`) are **reserved** and cannot be used for anything else.
- The Turing Machine starts in the `initial state`, pointing at the **first element** of the `initial tape`.
- Each section must be **terminated by a semicolon** (`;`).
- `B` is a **reserved symbol** representing a blank space on the tape.

#### Examples

You can find example Turing Machine programs in the [examples/turing-machine](examples/turing-machine) directory.

### Deterministic Finite Automaton

A **Deterministic Finite Automaton (DFA)** determines its next move based on its current state and the input symbol. At each step, the automaton transitions to a new state and reads the next input symbol. The computation ends once the entire input has been processed.

If no transition is defined for a given state and symbol, the program terminates with an error.

#### Input Format

```
q0 q1 ... qn; [states]
qs; [initial state]
qf1 qf2 ... qfk; [accepting states]
a1 a2 ... an; [symbols]

(q, s) > (new_q)
(q, s) > (new_q)
(q, s) > (new_q)
...;

a1 a1 a3 a8 ...; [input]
```

#### Rules and Conventions
- Each state must start with the letter `q`, followed by one or more alphanumeric characters.
- Each symbol must constist of one or more alphanumeric characters.
- Each section must be **terminated by a semicolon** (`;`).

#### Examples

You can find example DFA programs in the [examples/deterministic-finite-automaton](examples/deterministic-finite-automaton) directory.

# Pushdown Automaton (PA)

A **Pushdown Automaton (PA)** determines its next move based on its current state, the input symbol, and the symbol at the top of the stack. At each step, the automaton transitions to a new state and may push an arbitrary number of symbols onto the stack.

At the start of the computation, the stack contains a single symbol (`}`), which serves as the stack start symbol. The user-provided input is concatenated with `{`, which acts as the input end symbol.

It is important to note that if the automaton attempts a transition when the stack is empty, an error is returned. The only valid moment for the stack to be empty is at the end of the computation.

#### Input Format

```
q0 q1 ... qn; [states]
qs; [initial state]
qf1 qf2 ... qfk; [accepting states]
a1 a2 ... an; [symbols]

(q, s_i, s_s) > (new_q, s_s1, s_s2, ...)
(q, s_i, s_s) > (new_q, s_s1, s_s2, ...)
(q, s_i, s_s) > (new_q, s_s1, s_s2, ...)
...;

a1 a1 a3 a8 ...; [input]
```

#### Rules and Conventions

- `{` is a **reserved symbol** representing the end of input. It **cannot** be used in the symbol declaration section but **must** be used in transitions.
- `}` is a **reserved symbol** representing the start of the stack. It **cannot** be used in the symbol declaration section but **must** be used in transitions.
- Each state must start with the letter `q`, followed by one or more alphanumeric characters.
- Each symbol must consist of one or more alphanumeric characters.
- Each section must be **terminated by a semicolon** (`;`).
- In the transitions section:
  - `s_i` represents the input symbol.
  - `s_s` represents the symbol from the top of the stack.
  - `s_s1`, `s_s2`, ... are symbols to be pushed onto the stack.
  - Symbols are pushed in the order they are provided, meaning `s_s2` will be closer to the top of the stack than `s_s1`.

#### Examples

You can find example DFA programs in the [examples/pushdown-automaton](examples/pushdown-automaton) directory.