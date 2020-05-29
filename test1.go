package main

import "./giftextgenerator"

import "os"

func main() {
	// save to out.gif
	ff, _ := os.OpenFile("approve.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer ff.Close()

	giftextgenerator.Generate(ff, "Test Test Test", "CooperOriginal.ttf", 12, "FF0000", "FFFFFF")
}
