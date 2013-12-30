package main

import (
	"flag"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MainSuite struct{}

var _ = Suite(&MainSuite{})

func (s *MainSuite) TestGenApi(c *C) {
	flag.Set("pkg", "github.com/hypermusk/hypermusk/tests/api")
	flag.Set("lang", "objc")
	flag.Set("outdir", "tests")
	flag.Set("impl", "github.com/hypermusk/hypermusk/tests/apiimpl")
	genApi()
}
