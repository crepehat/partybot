package partybot

import "testing"

func TestReadGridFile(t *testing.T) {
	nameGrid, err := ReadGridFile("./grid.csv")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Got grid file: %+v", nameGrid)
}

func TestNewGrid(t *testing.T) {
	nameGrid, err := ReadGridFile("./grid.csv")
	if err != nil {
		t.Error(err)
	}
	_, err = NewGrid(nameGrid)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Loaded grid")
}

func TestStartGrid(t *testing.T) {
	nameGrid, err := ReadGridFile("./grid.csv")
	if err != nil {
		t.Error(err)
	}
	grid, err := NewGrid(nameGrid)
	if err != nil {
		t.Error(err)
	}

	grid.Start()
	t.Logf("Initialised grid %+v", grid)
}
