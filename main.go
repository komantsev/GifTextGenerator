package main

import (
	"./giftextgenerator"
	"fmt"
	"image/color"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/generate", func(writer http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(writer, "ParseForm() err: %v", err)
				return
			}
			//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			text := r.FormValue("text")
			textSize := r.FormValue("text_size")
			textColor := r.FormValue("text_color")
			bgColor := r.FormValue("bg_color")

			writer.Header().Set("Content-Disposition", "attachment; filename=test.gif")
			writer.Header().Set("Content-Type", "image/gif")

			size, _ := strconv.ParseFloat(textSize, 64)
			giftextgenerator.Generate(writer, text, "CooperOriginal.ttf", size, ParseHexColor(textColor), ParseHexColor(bgColor))
		}
	})

	http.ListenAndServe(":8011", nil)
}

func ParseHexColor(s string) (c color.RGBA) {
	c.A = 0xff
	fmt.Sscanf(s,"#%02x%02x%02x", &c.R, &c.G, &c.B)
	return
}