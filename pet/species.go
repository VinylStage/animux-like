package pet

type SpeciesType string

const (
	Cat     SpeciesType = "cat"
	Dog     SpeciesType = "dog"
	Penguin SpeciesType = "penguin"
)

type Species struct {
	Name        string
	HungerDecay int // amount of decay per tick
	HappyDecay  int
	CleanDecay  int
}

var SpeciesData = map[SpeciesType]Species{
	Cat:     {"Cat", 2, 1, 1},
	Dog:     {"Dog", 3, 2, 2},
	Penguin: {"Penguin", 1, 1, 1},
}

func GetSpeciesList() []SpeciesType {
	return []SpeciesType{Cat, Dog, Penguin}
}
