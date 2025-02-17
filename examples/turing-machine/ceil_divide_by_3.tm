# This turing machine calculates ceil(n/3)
# it uses unary system for simplicity

# Idea:
# For input 1111111 we take leftmost 1, change it to X, go to the last 1, change it to B, change previous 1 to B and go back.
# After first step we have X1111BB.
# We do it again and we have XX1BBBB.
# Now we change last 1 to X, but after that there is no more 1 so we end up with:
# XXXBBBB
# We change every X to 1, and this is our final output:
# 111BBBB
# indeed, ceil(7/3) = 3
# In actual algorithm first X is gonna be different value, G for example, as we need a guard so
# we don't go out of tape.

# States
qStart # change 1 to G and go to the right until finding B, switch to qGoEnd
qGoEnd # go to the first B, switch to qRemoveFromEndFirst
qRemoveFromEndFirst # change 1 to B, switch to qRemoveFromEndSecond
qRemoveFromEndSecond # change 1 to B, switch to qFindNext
qFindNext # go back until G or X is reached, switch to qNext
qNext # same as qStart, but use X
qBack # go to G and change it to 1, switch to qFinish
qFinish # change X to 1 until B is reached, then switch to qAcc
qAcc
;

# Initial State
qStart;

# Accepting states
qAcc;

# Symbols
1 G X;

# Transitions
# qStart
(qStart, B) > (qAcc, 1, R) # case of empty input
(qStart, 1) > (qGoEnd, G, R) 
# qGoEnd
(qGoEnd, 1) > (qGoEnd, 1, R)
(qGoEnd, X) > (qGoEnd, X, R)
(qGoEnd, B) > (qRemoveFromEndFirst, B, L)
# qRemoveFromEndFirst
(qRemoveFromEndFirst, 1) > (qRemoveFromEndSecond, B, L)
(qRemoveFromEndFirst, X) > (qBack, X, L)
(qRemoveFromEndFirst, G) > (qAcc, 1, R)
# qRemoveFromEndSecond
(qRemoveFromEndSecond, 1) > (qFindNext, B, L)
(qRemoveFromEndSecond, X) > (qBack, X, L)
(qRemoveFromEndSecond, G) > (qAcc, 1, R)
# qFindNext
(qFindNext, 1) > (qFindNext, 1, L)
(qFindNext, X) > (qNext, X, R)
(qFindNext, G) > (qNext, G, R) 
# qNext
(qNext, 1) > (qGoEnd, X, R)
(qNext, B) > (qBack, B, L)
# qBack
(qBack, X) > (qBack, X, L)
(qBack, G) > (qFinish, 1, R)
# qFinish
(qFinish, X) > (qFinish, 1, R)
(qFinish, B) > (qAcc, B, L)
;

# Initial tape
1 1 1 1 1 1 1 ; # 7, so result should be 3 (111 in unary)