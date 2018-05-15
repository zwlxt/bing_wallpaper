package main

import (
	"image/color"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	WallpaperDir     string
	FontFile         string
	FontSize         float64
	Width            int
	OffsetX, OffsetY int
	TextColor        color.RGBA
	BackgroundColor  color.RGBA
}

func (c Config) Save(filename string) {
	y, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, y, 0644)
}

func (c *Config) Load(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		panic(err)
	}
}

func Default() Config {
	return Config{
		FontFile:        "simsun.ttc",
		FontSize:        20,
		Width:           500,
		OffsetX:         1500,
		OffsetY:         50,
		TextColor:       color.RGBA{255, 255, 255, 255},
		BackgroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 100},
		WallpaperDir:    "./",
	}
}
