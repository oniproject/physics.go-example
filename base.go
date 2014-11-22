package main

import (
	"github.com/oniproject/physics.go"
	"github.com/oniproject/physics.go/behaviors"
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	//"github.com/oniproject/physics.go/renderers"
	. "./debug-renderer"
	"log"
	"time"
)

func main() {
	w, h := 800, 600

	renderer, _ := NewRendererSDL("Simple", w, h)

	world := physics.NewWorldImprovedEuler()
	world.SetRenderer(renderer)
	step := func(interface{}) {
		world.Render()
	}
	world.On("step", &step)

	circle := bodies.NewCircle(80)
	circle.SetPosition(float64(w)*0.4, float64(h)*0.3)
	circle.SetVelocity(0.3, 0)
	circle.State().Angular.Vel = 0.001

	//bounds := //geom.NewAABB_byWH(float64(w), float64(h))
	bounds := geom.NewAABB_byMM(0, 0, float64(w), float64(h))
	log.Println(bounds)
	world.Add(circle)
	world.Add(behaviors.NewBodyImpulseResponse())
	world.Add(behaviors.NewConstantAcceleration(0, 0.0004))
	world.Add(behaviors.NewEdgeCollisionDetection(bounds, 0.99, 0.99))

	c := time.Tick(time.Second / 30)
	for now := range c {
		world.Step(now)
	}
}
