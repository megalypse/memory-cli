package cmd

import "sync"

var heightMu = sync.Mutex{}
var widthMu = sync.Mutex{}

var height int
var width int

func SetHeight(h int) {
	heightMu.Lock()
	defer heightMu.Unlock()

	height = h
}

func GetHeight() int {
	heightMu.Lock()
	defer heightMu.Unlock()

	return height
}

func SetWidth(w int) {
	widthMu.Lock()
	defer widthMu.Unlock()

	width = w
}

func GetWidth() int {
	widthMu.Lock()
	defer widthMu.Unlock()

	return width
}

func GetWidthRatio(ratio float64) int {
	return int(float64(GetWidth()) * ratio)
}

func GetHeightRatio(ratio float64) int {
	return int(float64(GetHeight()) * ratio)
}
