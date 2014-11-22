package renderers

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
	"github.com/oniproject/physics.go/renderers"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type RendererSDL struct {
	Window *sdl.Window
}

func NewRendererSDL(title string, w, h int) (r *RendererSDL, err error) {
	r = &RendererSDL{}
	r.Window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	return
}

const _dt = 1000 / 30

/*func (r *RendererSDL) Run() {
	for {
		r.Window.UpdateSurface()
		sdl.Delay(_dt)
	}
}*/

func (r *RendererSDL) SetWorld(world renderers.World) {}

func (r *RendererSDL) Render(arr []bodies.Body, meta interface{}) {
	surface := r.Window.GetSurface()
	w, h := r.Window.GetSize()

	// clear
	surface.FillRect(&sdl.Rect{0, 0, int32(w), int32(h)}, 0xffffffff)

	ww, hh := float64(w), float64(h)
	drawLine(surface, geom.Vector{ww, 0}, geom.Vector{0, hh}, 0xff00cc00)
	drawLine(surface, geom.Vector{ww / 2, 0}, geom.Vector{ww / 2, hh}, 0xff00cc00)
	drawLine(surface, geom.Vector{ww / 2, hh}, geom.Vector{ww / 2, 0}, 0xff00cc00)

	// aabb
	for _, body := range arr {
		aabb := body.AABB(0)
		rect := &sdl.Rect{
			int32(aabb.X - aabb.HW),
			int32(aabb.Y - aabb.HH),
			int32(aabb.HW * 2),
			int32(aabb.HH * 2),
		}
		surface.FillRect(rect, 0xff0000cc)
	}

	for _, body := range arr {
		switch g := body.Geometry().(type) {
		case *geometries.Circle:
			state := body.State()
			drawCircle(surface, int(state.Pos.X), int(state.Pos.Y), int(g.Radius), 0xffcc0000)
			fillCircle(surface, int(state.Pos.X), int(state.Pos.Y), int(g.Radius), 0xffcc0000)
			apos := state.Angular.Pos
			a := geom.Vector{math.Cos(apos), math.Sin(apos)}
			pos := state.Pos.Plus(a.Times(g.Radius))
			drawLine(surface, state.Pos, pos, 0xff00cc00)
		}
	}

	r.Window.UpdateSurface()
}

/*
	CreateView(geometry geometries.Geometry , styles) interface{}
	DrawBody(body bodies.Body, view interface{})
	DrawMeta(meta interface{})
	Render(bodies []bodies.Body, meta interface{})
	SetWorld(world util.EventTarget)
*/

func setPixel(surface *sdl.Surface, x, y int, pixel uint32) {
	if x < 0 || y < 0 || x >= int(surface.W) || y >= int(surface.H) {
		return
	}
	pos := y*surface.Pitch + x*4
	if pos < 0 {
		return
	}
	setPixelPos(surface, pos, pixel)
}
func setPixelPos(surface *sdl.Surface, pos int, pixel uint32) {
	surface.Pixels()[pos+3] = byte((pixel >> 24) & 0xff) // a
	surface.Pixels()[pos+2] = byte((pixel >> 16) & 0xff) // r
	surface.Pixels()[pos+1] = byte((pixel >> 8) & 0xff)  // g
	surface.Pixels()[pos+0] = byte((pixel >> 0) & 0xff)  // b
}

func drawCircle(surface *sdl.Surface, n_cx, n_cy, radius int, pixel uint32) {
	err := -float64(radius)
	x := float64(radius) - 0.5
	y := float64(0.5)
	cx := float64(n_cx) - 0.5
	cy := float64(n_cy) - 0.5

	for x >= y {
		setPixel(surface, int(cx+x), int(cy+y), pixel)
		setPixel(surface, int(cx+y), int(cy+x), pixel)

		if x != 0 {
			setPixel(surface, int(cx-x), int(cy+y), pixel)
			setPixel(surface, int(cx+y), int(cy-x), pixel)
		}

		if y != 0 {
			setPixel(surface, int(cx+x), int(cy-y), pixel)
			setPixel(surface, int(cx-y), int(cy+x), pixel)
		}

		if x != 0 && y != 0 {
			setPixel(surface, int(cx-x), int(cy-y), pixel)
			setPixel(surface, int(cx-y), int(cy-x), pixel)
		}

		err += y
		y++
		err += y

		if err >= 0 {
			x--
			err -= x
			err -= x
		}
	}
}

func fillCircle(surface *sdl.Surface, cx, cy, radius int, pixel uint32) {
	const BPP = 4

	r := float64(radius)

	for dy := float64(1); dy <= r; dy += 1.0 {
		dx := math.Floor(math.Sqrt(2.0*r*dy - dy*dy))

		ay := int(float64(cy) + r - dy)
		by := int(float64(cy) - r + dy)

		for x := int(float64(cx) - dx); x <= cx+int(dx); x++ {
			setPixel(surface, x, ay, pixel)
			setPixel(surface, x, by, pixel)
		}

	}
}

func drawLine(surface *sdl.Surface, a, b geom.Vector, pixel uint32) {

	steep := math.Abs(b.Y-a.Y) > math.Abs(b.X-a.X)

	if steep {
		a.X, a.Y = a.Y, a.X
		b.X, b.Y = b.Y, b.X
	}

	if a.X > b.X {
		a.X, b.X = b.X, a.X
		a.Y, b.Y = b.Y, a.Y
	}

	if a.X == b.X {
		// vertical
		for y := int(a.Y); y < int(b.Y); y++ {
			setPixel(surface, int(a.X), y, pixel)
		}
		return
	}

	dx := b.X - a.X
	dy := math.Abs(b.Y - a.Y)

	err := dx * 0.5
	ystep := -1
	if a.Y < b.Y {
		ystep = 1
	}

	y := int(a.Y)
	maxX := int(b.X)

	for x := int(a.X); x < maxX; x++ {
		if steep {
			setPixel(surface, y, x, pixel)
		} else {
			setPixel(surface, x, y, pixel)
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
}
