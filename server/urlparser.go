package main

import (
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strings"
	"errors"
)

func ExtToFormat(ext string) string {
	ext = strings.ToUpper(strings.TrimLeft(ext, "."))
	if ext == "JPG" {
		return "JPEG"
	}
	if ext == "TIF" {
		return "TIFF"
	}
	return ext
}

func ParseURI(uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker, conf config.Conf) (operations []ops.Operation, ext string, err, invalid error) {
	filename := string(uri.Path())
	w, h, cropx, cropy, mode, invalid := getParams(uri.QueryArgs())
	ow, oh, err := source.ImageSize(filename)
	if invalid != nil {
		return
	}
	if err != nil {
		return
	}

	operations = []ops.Operation{source.LoadImageOp(filename)}

	adjustWidth := func() {
		w = roundedIntegerDivision(h*ow, oh)
	}

	adjustHeight := func() {
		h = roundedIntegerDivision(w*oh, ow)
	}

	adjustSize := func() {
		if h == 0 && w != 0 {
			adjustHeight()
		} else if h != 0 && w == 0 {
			adjustWidth()
		} else if w == 0 && h == 0 {
			w, h = ow, oh
		}
	}

	denyUpscale := func() {
		h0 := h
		w0 := w
		if w > ow {
			h = roundedIntegerDivision(ow*h0, w0)
			w = ow
		}
		if h > oh || h > h0 {
			w = roundedIntegerDivision(oh*w0, h0)
			h = oh
		}
	}

	resize := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.Resize{w, h})
	}

	liquid := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.LiquidRescale{w, h})
	}

	crop := func() {
	  if w == 0 {
	    w = ow
	  }
	  if h == 0 {
	    h = oh
	  }
		operations = append(operations, ops.Crop{w, h, cropx, cropy})
	}

	cropmid := func() {
		if w == 0 || w > ow {
	    w = ow
	  }
	  if h == 0 || h > oh {
	    h = oh
	  }
	  midW := roundedIntegerDivision(ow, 2)
	  midH := roundedIntegerDivision(oh, 2)
	  cropx := midW - roundedIntegerDivision(w, 2)
	  cropy := midH - roundedIntegerDivision(h, 2)
		operations = append(operations, ops.Crop{w, h, cropx, cropy})
	}

	fit := func() {
		if w > ow {
			w = ow
		}
		if h > oh {
			h = oh
		}
		if w != 0 && h != 0 {
			if ow*h > w*oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{w, h})
		} else {
			resize()
		}
	}

	watermark := func() {
		heightOK := h > marker.Minheight && h < marker.Maxheight
		widthOK := w > marker.Minwidth && w < marker.Maxwidth
		if marker.AddMark && heightOK && widthOK {
			logging.Debug("Adding watermarkOp")
			operations = append(operations, ops.WatermarkOp(marker.WatermarkImage, marker.Horizontal, marker.Vertical))
		}
	}

	switch mode {
	case "resize":
		resize()
	case "fit":
		fit()
	case "liquid":
		liquid()
	case "crop":
		crop()
	case "cropmid":
		cropmid()
	default:
		resize()
	}
	watermark()

	ext = filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	operations = append(operations, ops.Convert{ExtToFormat(ext)})

	return
}

func roundedIntegerDivision(n, m int) int {
	if (n < 0) == (m < 0) {
		return (n + m/2) / m
	} else { // -5 / 6 should round to -1
		return (n - m/2) / m
	}
}


// returns validated parameters from request and error if invalid
func getParams(a *fasthttp.Args) (w, h, cropx, cropy int, mode string, e error) {
	w = a.GetUintOrZero("width")
	h = a.GetUintOrZero("height")
	cropx = a.GetUintOrZero("cropx")
	cropy = a.GetUintOrZero("cropy")
	mode = string(a.Peek("mode"))
	modes := map[string]bool {
		"": true,
		"fit": true,
		"crop": true,
		"cropmid": true,
		"liquid": true,
	}

	if strings.Contains(a.String(), "%") {
		e = errors.New("Invalid characters in request!")
		return
	}
	if w == 0 && a.Has("width") {
		e = errors.New("Invalid width!")
		return
	}
	if h == 0 && a.Has("height") {
		e = errors.New("Invalid height!")
		return
	}
	if cropx == 0 && a.Has("cropx") {
		e = errors.New("Invalid cropx!")
		return
	}
	if cropy == 0 && a.Has("cropy") {
		e = errors.New("Invalid cropy!")
		return
	}
	if !modes[mode] {
		e = errors.New("Invalid mode!")
		return
	}
	a.Del("width")
	a.Del("height")
	a.Del("mode")
	a.Del("cropx")
	a.Del("cropy")
	if a.Len() != 0 {
		e = errors.New("Invalid parameter " + a.String())
		return
	}
	return
}
