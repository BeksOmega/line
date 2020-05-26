package line

import "image"

// BresenhamPlotter is a function which takes in: an x coordinate, a y coordinate,
// and an error value for the current coordinate. The error value is associated
// with the distance between the ideal line and the center of the pixel at (x, y).
type BresenhamPlotter func(x, y, e int)

// Bresenham plots a line between the two points using Bresenham's line algorithm.
func Bresenham(x0, y0, x1, y1 int, plot BresenhamPlotter) {
	p0 := image.Pt(x0, y0)
	p1 := image.Pt(x1, y1)
	d := p1.Sub(p0)

	switch {
	case !isLine(d):
		return

	case isVert(d):
		drawVert(p0, p1, d, plot)
	case isHoriz(d):
		drawHoriz(p0, p1, d, plot)
	case isDiagonal(d):
		drawDiagonal(p0, p1, d, plot)

	default:
		bresenhamAll(p0, p1, d, plot)
	}
}

func isLine(d image.Point) bool {
	return d.X != 0 || d.Y != 0
}

func isVert(d image.Point) bool {
	return d.X == 0
}

func isHoriz(d image.Point) bool {
	return d.Y == 0
}

func isDiagonal(d image.Point) bool {
	return abs(d.X) == abs(d.Y)
}


func drawVert(p0, p1, d image.Point, plot BresenhamPlotter) {
	inc := sign(d.Y)
	x := p0.X
	for y := p0.Y; y < p1.Y; y += inc {
		plot(x, y, 0)
	}
}

func drawHoriz(p0, p1, d image.Point, plot BresenhamPlotter) {
	inc := sign(d.X)
	y := p0.Y
	for x := p0.X; x < p1.X; x += inc {
		plot(x, y, 0)
	}
}

func drawDiagonal(p0, p1, d image.Point, plot BresenhamPlotter) {
	s := mapPt(d, sign)
	for p := p0; p.X != p1.X; p = p.Add(s) {
		plot(p.X, p.Y, 0)
	}
}

// bresenhamAll takes the input and modifies it so that it can be understood
// by bresenham simple, which only supports lines with a slope between (0, 1).
func bresenhamAll(p0, p1, d image.Point, plot BresenhamPlotter) {
	if abs(d.X) < abs(d.Y) {  // Switch x and y.
		plot = func (x, y, e int) {
			plot(y, x, e)
		}
		p0 = image.Pt(p0.Y, p0.X)
		d = image.Pt(d.Y, d.X)
	}
	if d.X < 0 {  // End is less than start.
		temp := p0
		p0 = p1
		p1 = temp
	}
	bresenhamSimple(p0, p1, d, plot)
}

// bresenhamSimple plots lines with a slope between (0, 1). Each pixel is
// plotted by passing the x coord, y coord, and error value to the plotter
// function.
func bresenhamSimple(p0, p1, d image.Point, plot BresenhamPlotter) {
	yi := sign(d.Y)
	d.Y = abs(d.Y)
	y := p0.Y
	e := 2*d.Y - d.X

	for x := p0.X; x != p1.X; x++ {
		plot(x, y, e)
		if e > 0 {
			y += yi
			e -= 2*d.X
		}
		e += 2*d.Y
	}
}


// http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go
func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

func sign(n int) int {
	return n / abs(n)
}

func mapPt(p image.Point, f func(int) int) image.Point {
	return image.Pt(
		f(p.X),
		f(p.Y),
	)
}
