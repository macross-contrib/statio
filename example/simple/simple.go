package main

import (
	"github.com/macross-contrib/statio"
	"macross"
)

func main() {
	m := macross.New()

	// if Allow DirectoryIndex
	//m.Use(statio.Serve("/", statio.LocalFile("/tmp", true)))
	// set prefix
	//m.Use(statio.Serve("/static", statio.LocalFile("/tmp", true)))

	m.Use(statio.Serve("/", statio.LocalFile("/tmp", false)))
	m.Get("/ping", func(self *macross.Context) error {
		return self.String("pong")
	})
	// Listen and Server in 0.0.0.0:8080
	m.Listen(":8080")
}
