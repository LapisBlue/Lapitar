package main

import (
	"fmt"
	"github.com/LapisBlue/Tar/head"
	"github.com/LapisBlue/Tar/server"
	"github.com/LapisBlue/Tar/skin"
	"image/png"
	"log"
	"os"
)

const (
	width  = 256
	height = 256
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Usage: tar <player>")
		os.Exit(1)
	}

	player := os.Args[1]

	if player == "start" {
		log.Fatalln(server.Start(os.Args[2]))
		os.Exit(0)
	}

	log.Println("Downloading skin:", player)
	sk, err := skin.Download(player)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Creating head for", player)

	img, err := head.Render(sk, 45, 256, 256, 4, true, true, true)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(player + ".png")
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatalln(err)
	}

	file.Close()
}
