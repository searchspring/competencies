package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestProcessInherited(t *testing.T) {
	result, err := processInherits("#original \n##Skills\n<inherit doc=\"../test1.md\"/>\n   a  ", false)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
	require.Equal(t, "#original \n##Skills\n<skills>\nbreakdancing\n</skills>\n<skills>\nfigure skating\n</skills>\n<skills>\nkung fu\n</skills>\n   a", result)
}
