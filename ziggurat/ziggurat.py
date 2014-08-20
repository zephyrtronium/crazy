#!/usr/bin/env python3

# Copyright (c) 2014 Branden J Brown

# This software is provided 'as-is', without any express or implied
# warranty. In no event will the authors be held liable for any damages
# arising from the use of this software.

# Permission is granted to anyone to use this software for any purpose,
# including commercial applications, and to alter it and redistribute it
# freely, subject to the following restrictions:

# 1. The origin of this software must not be misrepresented; you must not
#    claim that you wrote the original software. If you use this software
#    in a product, an acknowledgment in the product documentation would be
#    appreciated but is not required.
# 2. Altered source versions must be plainly marked as such, and must not be
#    misrepresented as being the original software.
# 3. This notice may not be removed or altered from any source distribution.

"""
Calculate the parameters of the n-segment ziggurat for any monotonically
decreasing probability distribution.

The program uses mpmath for quadrature of the pdf's tail as well as to provide
a wide variety of mathematical functions for the use of the pdf. Calling the
program is as simple as providing the pdf, e.g.

	python3 ziggurat.py "exp(-0.5*x*x)" >normal.out

The parameters will be calculated with periodic updates printed to stderr and
the results printed to stdout in Go syntax.

It is assumed that the pdf is supported on the interval [0, inf) and that the
only variable is x.

See http://www.jstatsoft.org/v05/i08/paper.

"""

import itertools
import sys

import mpmath

def solve(f, x0=None, nseg=128, verbose=False):
	"""Find r, v, and the x coordinates for f."""
	r = x0
	if r is None:
		if verbose:
			print('Calculating initial guess... ', end='', file=sys.stderr)
		# The area we seek is nseg equal-area rectangles surrounding f(x), not
		# f(x) itself, but we can get a good approximation from it.
		v = mpmath.quad(f, [0, mpmath.inf]) / nseg
		r = mpmath.findroot(lambda x: x*f(x) + mpmath.quad(f, [x, mpmath.inf]) - v, x0=1, x1=100, maxsteps=100)
		if verbose:
			print(r, file=sys.stderr)
	# We know that f(0) is the maximum because f must decrease monotonically.
	maximum = f(0)
	txi = []
	tv = mpmath.mpf()
	def mini(r):
		nonlocal txi, tv
		xi = [r]
		y = f(r)
		v = r*f(r) + mpmath.quad(f, [r, mpmath.inf])
		if verbose:
			print('Trying r={0} (v={1})'.format(r, v), file=sys.stderr)
		for i in itertools.count():
			xm1 = xi[i]
			h = v / xm1
			y += h
			if y >= maximum or mpmath.almosteq(y, maximum, abs_eps=mpmath.mp.eps * 2**10):
				break
			# We solve for x via secant method instead of using f's inverse.
			x = mpmath.findroot(lambda x: f(x) - y, xm1)
			xi.append(x)
		xi.append(mpmath.mpf())
		if len(xi) == nseg:
			if mpmath.almosteq(y, maximum, abs_eps=mpmath.mp.eps * 2**10):
				txi, tv = xi[::-1], v
				return 0
			# If y > maximum, then v is too large, which means r is too far
			# left, so we want to return a negative value. The opposite holds
			# true when y < maximum.
			return maximum - y
		return len(xi) - nseg + h*mpmath.sign(len(xi) - nseg)
	r = mpmath.findroot(mini, r)
	assert len(txi) == nseg
	if verbose:
		print('Done calculating r, v, x[i].', file=sys.stderr)
	return r, tv, txi

def tables(f, r, v, xi, verbose=False):
	"""Calculate k[i], w[i], and f[i]."""
	ki = [None] * len(xi)
	wi = ki[:]
	fi = ki[:]
	for i, x in enumerate(xi):
		if verbose and i & 7 == 0:
			print('\r{0}/{1}'.format(i, len(xi)), end='', file=sys.stderr)
		if i == 0:
			ki[0] = mpmath.floor(2**32 * r*f(r)/v)
			wi[0] = v / f(r) / 2**32
		else:
			ki[i] = mpmath.floor(2**32 * xi[i-1]/x)
			wi[i] = x / 2**32
		fi[i] = f(x)
	if verbose:
		print('\r{0}/{0}'.format(len(xi)), file=sys.stderr)
	assert all(v is not None for v in ki)
	assert all(v is not None for v in wi)
	assert all(v is not None for v in fi)
	return ki, wi, fi

def format(r, v, xi, ki, wi, fi):
	'''Turn parameters into Go "pseudocode."'''
	return '''

const R = {r}
const V = {v}

var X = [{n}]float32{{{xi}}}

var K = [{n}]uint32{{{ki}}}

var W = [{n}]float32{{{wi}}}

var F = [{n}]float32{{{fi}}}

'''.format(
		r=r,
		v=v,
		n=len(xi),
		xi=', '.join(str(i) for i in xi),
		ki=', '.join(hex(int(i)) for i in ki),
		wi=', '.join(str(i) for i in wi),
		fi=', '.join(str(i) for i in fi),
	)

def main(*args):
	e = compile(args[0], '<f>', 'eval')
	globs = mpmath.__dict__
	globs['__builtins__'] = None
	mpmath.mp.prec = 80
	f = lambda x: eval(e, globs, {'x': x})
	r, v, xi = solve(f, verbose=True)
	ki, wi, fi = tables(f, r, v, xi, verbose=True)
	mpmath.mp.prec = 53
	print(format(r, v, xi, ki, wi, fi))

if __name__ == '__main__':
	main(*sys.argv[1:])