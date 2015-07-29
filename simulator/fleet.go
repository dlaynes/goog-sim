package simulator

import (
	"fmt"
)

type Player struct {
	Id            string
	Name          string
	MainPlanet    *Position
	ShipTypes     []*ShipType
	MilitaryTech  int
	DefenseTech   int
	HullTech      int
	OriginalFleet map[string]int
}

/* unused */
type Fleet struct {
}

type ShipType struct {
	Res        *Resource
	BaseAttack float64
	BaseShield float64
	BaseHull   float64
	Amount     int
	Explosions int
	Statistics []int //Remaining ships after every battle
	Rapidfires map[string]float64
}

type ShipUnit struct {
	Type     *ShipType
	Attack   float64
	Shield   float64
	Hull     float64
	Exploded bool
}

type FleetGroup struct {
	Ships           []*ShipUnit
	ThisTurnDamage  float64
	ThisTurnDefense float64
	ThisTurnAttacks int
}

/* Player functions */

func NewPlayer() *Player {
	pl := &Player{}

	return pl
}

func (this *Player) Expand(group *FleetGroup, resources map[string]Resource) {
	this.processShipTypes(resources)

	ships := &group.Ships

	//fmt.Println("Expandiendo")

	for _, t := range this.ShipTypes {
		*ships = append(*ships, t.Expand()...)

	}
}

func (this *Player) processShipTypes(resources map[string]Resource) {

	var i = 0

	fmt.Println("Cantidad Items ", len(this.OriginalFleet))

	for id, amount := range this.OriginalFleet {

		st := &ShipType{}
		st.Init(resources[id], amount, this.MilitaryTech, this.DefenseTech, this.HullTech)
		this.ShipTypes = append(this.ShipTypes, st)
		i++
	}
}

/* ShipType functions */

func (this *ShipType) Init(resource Resource, amount int, mtech int, dtech int, htech int) {
	this.Amount = amount

	d := 1 + (float64(mtech) * 0.1)
	s := 1 + (float64(dtech) * 0.1)
	h := (1 + (float64(htech) * 0.1)) * 0.1

	this.BaseAttack = d * resource.Attack
	this.BaseShield = s * resource.Defense
	this.BaseHull = h * resource.Hull

	this.Rapidfires = resource.Rapidfires
}

func (this *ShipType) Expand() []*ShipUnit {
	sh := make([]*ShipUnit, this.Amount)

	for i := 0; i < this.Amount; i++ {
		sh[i] = &ShipUnit{Type: this, Attack: this.BaseAttack, Shield: this.BaseShield, Hull: this.BaseHull}
	}

	fmt.Printf("%#v\n", sh)

	return sh
}

/* TO DO */
func (this *ShipType) CalcCapacity() {

}

func (this *ShipType) LogBattle() {

}

/* FleetGroup functions */

func (this *FleetGroup) Init() {
	this.Ships = []*ShipUnit{}
}

func (this *FleetGroup) attack(other *FleetGroup) bool {

	return true
}

func (this *FleetGroup) clean() {

}
