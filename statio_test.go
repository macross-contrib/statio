package statio_test

import (
	"github.com/macross-contrib/statio"
	"macross"
	"testing"
)

func TestStatic(t *testing.T) {
	m := macross.New()
	m.Use(statio.Serve("/", statio.LocalFile("public", false)))
	go m.Run(":8888")
}
