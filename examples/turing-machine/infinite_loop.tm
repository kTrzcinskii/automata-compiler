# two states, we will be jumping from q0 to q1 and from q1 to q0 without any accepting state
q0 q1;
q0;
; # no accepting state
1;
# jumping left and right
(q0, 1) > (q1, 1, R)
(q1, 1) > (q0, 1, L);
1 1;
