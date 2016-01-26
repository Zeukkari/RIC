package ops

import "github.com/phzfi/RIC/server/images"

type LiquidRescale struct {
	Width, Height uint
}

func (r LiquidRescale) Apply(img images.Image) error {
	// The third argument to LiquidRescaleImage is the maximum transversal step, or how many pixels a seam can move sideways.
	// The fourth is the rigidity, which makes seams less steep. Recommended if transversal step is greater than one.
	return img.LiquidRescaleImage(r.Width, r.Height, 1, 0)
}
