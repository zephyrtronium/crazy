Calculate the parameters of the n-segment ziggurat for any monotonically
decreasing probability distribution.

The program uses mpmath for quadrature of the pdf's tail as well as to provide
a wide variety of mathematical functions for the use of the pdf. Calling the
program is as simple as providing the pdf, e.g.

	python3 ziggurat.py "exp(-0.5*x*x)" >normal.out

You can also provide `-h` to get a list of options.

The parameters will be calculated with periodic updates printed to stderr and
the results printed to stdout in Go syntax.

It is assumed that the pdf is supported on the interval `[0, inf)` and that the
only variable is `x`.

See <http://www.jstatsoft.org/v05/i08/paper>.