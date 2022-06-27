package processing

import "gopkg.in/gographics/imagick.v3/imagick"

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	mw.SetSize(400, 400)

	fc := imagick.NewPixelWand()
	fc.SetColor("none")
	mw.NewImage(400, 400, fc)
	mw.SetImageArtifact("gradient:angle", "35")

	red := imagick.NewPixelWand()
	if ok := red.SetColor("#ff0000"); !ok {
		panic(mw.GetLastError())
	}

	green := imagick.NewPixelWand()
	if ok := green.SetColor("red"); !ok {
		panic(mw.GetLastError())
	}

	blue := imagick.NewPixelWand()
	if ok := blue.SetColor("blue"); !ok {
		panic(mw.GetLastError())
	}

	stops := []imagick.StopInfo{
		imagick.NewStopInfo(red.GetMagickColor(), 0),
		imagick.NewStopInfo(green.GetMagickColor(), 0.2),
		imagick.NewStopInfo(green.GetMagickColor(), 0.5),
		imagick.NewStopInfo(blue.GetMagickColor(), 1),
	}

	err := mw.GradientImage(imagick.GRADIENT_TYPE_LINEAR, imagick.SPREAD_METHOD_PAD, stops)
	if err != nil {
		panic(err)
	}

	//mw.WriteImage("out.png")

	lw := imagick.NewMagickWand()
	pw := imagick.NewPixelWand()
	dw := imagick.NewDrawingWand()

	pw.SetColor("none")
	lw.NewImage(400, 400, pw)

	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.RoundRectangle(15, 15, 385, 385, 100, 100)
	lw.DrawImage(dw)

	lw.CompositeImage(mw, imagick.COMPOSITE_OP_SRC_OUT, false, 0, 0)
	lw.WriteImage("mask_result.png")
}
