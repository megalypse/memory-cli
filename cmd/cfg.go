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

func SetWidth(w int) {
	widthMu.Lock()
	defer widthMu.Unlock()

	width = w
}

func GetWidthRatio(ratio float64) int {
	widthMu.Lock()
	defer widthMu.Unlock()

	return int(float64(width) * ratio)
}

func GetHeightRatio(ratio float64) int {
	heightMu.Lock()
	defer heightMu.Unlock()

	return int(float64(height) * ratio)
}
