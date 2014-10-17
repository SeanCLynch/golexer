32-bit Integer, Octal and Hexadecimal Lexical Analyser
by Sean C Lynch

-----

Execution Procedure:

To use normally, one should type: "./golexer num", where 'num' is any desired input. 

The result will tell you if the emited token is valid, what it looks like (const), what its base is, what its decimal value is, and if there was any errors.

If you would like to run tests, please type "./golexer tests". This will show a list of "valid" and "invalid" tests, and give a more detailed procedure of how the machine works. 

(I've never submitted an assignment in GO (golang.org) before. If there are any problems running it, or really any problems at all, please contact me as soon as possible.)


-----

Should match: "((+|-)?[0-9]+)|([0-7]+[bB])|([0-9a-fA-F]+[hH])"

Notes:
Octals have base b
Hexadecimals have base h
Max (decimal) values are: 4294967295
Tests: see above.
