package statio_test

import (
	"github.com/insionng/macross"
	"github.com/macross-contrib/statio"
	"testing"
)

func TestStatic(t *testing.T) {
	m := macross.New()
	m.Use(statio.Serve("/", statio.LocalFile("public", false)))
	go m.Run(":8888")
}
