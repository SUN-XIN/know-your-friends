package geo

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
)

// Lat+Lng : coordinates on earth (]-90/+90[ | [-180/+180[)
// note : origin is SW ! (-180,-90)
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var (
	World = LLBox{
		N: 90, S: -90,
		E: 180, W: -180,
	}
)

type LLBox struct {
	N float64 `json:"n"`
	E float64 `json:"e"`
	S float64 `json:"s"`
	W float64 `json:"w"`
}

func (b LLBox) Validate() error {
	if b.N < b.S {
		return fmt.Errorf("north must be greater than south")
	}
	if b.E < b.W {
		return fmt.Errorf("east must be greater than west")
	}
	return nil
}

func (b LLBox) Empty() bool {
	return (b.W == 0 && b.E == 0) || (b.N == 0 && b.S == 0)
}

func (b *LLBox) String() string {
	return fmt.Sprintf("N:%f,E:%f,S:%f,W:%f", b.N, b.E, b.S, b.W)
}

// build a correct box from two points A & B, normalize the points
func NewLLBox(n, e, s, w float64) *LLBox {
	return &LLBox{N: n, E: e, S: s, W: w}
}

// box contains X,Y ?
func (box *LLBox) ContainsLL(lat, long float64) bool {
	if long < box.W || long > box.E {
		return false
	}

	if lat > box.N || lat < box.S {
		return false
	}

	return true
}

// box contains X,Y ?
func (box *LLBox) ContainsLatLng(latLng LatLng) bool {
	if latLng.Lng < box.W || latLng.Lng > box.E {
		return false
	}

	if latLng.Lat > box.N || latLng.Lat < box.S {
		return false
	}

	return true
}

// box1 contains box2 ?
func (box LLBox) ContainsZone(b LLBox) bool {
	return box.ContainsLL(b.N, b.E) &&
		box.ContainsLL(b.S, b.W)
}

// does this region contains a point
func (box *LLBox) Contains(p *LatLng) bool {
	return box.ContainsLL(p.Lat, p.Lng)
}

// Union computes the bounding box of 2 boxes
func (box *LLBox) Union(box2 *LLBox) *LLBox {
	return &LLBox{N: math.Max(box.N, box2.N), W: math.Min(box.W, box2.W), S: math.Min(box.S, box2.S), E: math.Max(box.E, box2.E)}
}

func (box *LLBox) RandomLatLngInside() (p LatLng) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p.Lat = (r.Float64() * (box.N - box.S)) + box.S
	p.Lng = (r.Float64() * (box.E - box.W)) + box.W
	return
}

// max distance in a zone: Diagonal
func (box *LLBox) MaxDistanceInKM() float64 {
	return DistanceFromPointsInKM(box.S, box.W, box.N, box.E)
}

func (box *LLBox) CenterPoint() (lat float64, lng float64) {
	return (box.N + box.S) / 2.0, (box.W + box.E) / 2.0
}

func DistanceFromPointsInKM(plat1, plng1, plat2, plng2 float64) float64 {
	dLat := (plat2 - plat1) * (math.Pi / 180.0)
	dLon := (plng2 - plng1) * (math.Pi / 180.0)

	lat1 := plat1 * (math.Pi / 180.0)
	lat2 := plat2 * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}
