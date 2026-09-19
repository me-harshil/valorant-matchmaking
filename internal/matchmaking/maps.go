package matchmaking

import "math/rand"

var mapPool = []string{
	"Abyss", "Ascent", "Bind", "Breeze", "Corrode", "Fracture",
	"Haven", "Icebox", "Lotus", "Pearl", "Split", "Summit", "Sunset",
}

func RandomMap() string {
	return mapPool[rand.Intn(len(mapPool))]
}
