package main

import (
	"errors"
	"fmt"

	"github.com/oowhyy/astroapp/pkg/vector"
)

var (
	maxDistMult = 2.0
	ErrFar      = errors.New("flew away")
	ErrCrash    = errors.New("crashed")
	PixelsPerAU = 200.0

	maxSimSteps = 1000000
)

var okSpeedLim bool

func main() {

	startX := []float64{
		0.387, // merc
		0.723, // ven
		1, // ea
		1.524, // ma
		5.203,  // jup
		9.537,  // satu
		19.191, // ura
		30.069, // nept
	}
	for _, val := range startX {
		gConst := 0.004
		steps := 1000000.0
		speedLim := 100 * PixelsPerAU
		sunMass := 330000.0
		bestSteps := 0
		var bestStart float64
		okSpeedLim = false
		for vely := float64(0.0); vely < speedLim; vely += PixelsPerAU / steps {
			if okSpeedLim {
				break
			}
			planet := &Body{
				Name: "earth",
				Pos:  vector.FromFloats(val*PixelsPerAU, 0),
				Vel:  vector.FromFloats(0, vely),
				Mass: 1, // cancels out later anyway
			}
			curSteps := simulate(gConst, sunMass, planet)
			if curSteps > bestSteps {
				bestSteps = curSteps
				bestStart = vely
			}
		}
		fmt.Printf("BestSteps: %d BestStart: %f\n", bestSteps, bestStart/PixelsPerAU)
	}
}

type Body struct {
	Name string
	Pos  vector.Vector
	Vel  vector.Vector
	Acc  vector.Vector

	Mass float64
}

func (b *Body) Update() {
	b.Vel.Scale(0.99996)
	b.Pos.Add(b.Vel)
	b.Vel.Add(b.Acc)
	b.Acc.Reset()
}

func (b *Body) ApplyForce(force vector.Vector) {
	acc := vector.Scaled(force, 1.0/b.Mass)
	b.Acc.Add(acc)
}

func (b *Body) DistTo(b2 *Body) float64 {
	return vector.Diff(b.Pos, b2.Pos).Len()
}

// assume forceSource is frozen at (0,0)
func simulate(gConst float64, centerMass float64, body *Body) int {
	startDist := body.Pos.Len()
	// startVel := body.Vel.Y
	centerBody := &Body{
		Name: "sun",
		Mass: centerMass,
	}
	recordDist := 0.0
	maxDist := body.Pos.Len() * maxDistMult
	recordMin := body.Pos.Len() * 10
	for i := 0; i < maxSimSteps; i++ {
		diff := vector.Diff(centerBody.Pos, body.Pos)
		distTo := diff.Len()
		recordDist = max(recordDist, distTo)
		recordMin = min(recordMin, distTo)
		if distTo > maxDist {
			if !okSpeedLim && i < 10 && recordMin >= startDist {
				okSpeedLim = true
				// fmt.Println("speed limit", startVel/PixelsPerAU)
			}
			return i
		}
		if distTo < 10 {
			return 0
		}
		diff.Scale(1 / distTo)
		forceMag := gConst * body.Mass * centerBody.Mass / (distTo * distTo)
		diff.Scale(forceMag)

		body.ApplyForce(diff)
		body.Update()

	}
	return maxSimSteps
}
