package simulator

import (
//"log"
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

func (this *Player) Expand(group *FleetGroup, resources *map[string]Resource) {
	this.processShipTypes(resources)

}

func (this *Player) processShipTypes(resources *map[string]Resource) {

	var i = 0

	for _, amount := range this.OriginalFleet {

		st := &ShipType{}
		st.Init(amount, this.MilitaryTech, this.DefenseTech, this.HullTech)
		this.ShipTypes = append(this.ShipTypes, st)
		i++
	}
}

/* ShipType functions */

func (this *ShipType) Init(amount int, mtech int, dtech int, htech int) {
	this.Amount = amount
}

/* FleetGroup functions */

func NewFleetGroup() *FleetGroup {
	group := &FleetGroup{}

	return group
}

func (this *FleetGroup) attack(other *FleetGroup) bool {

	return true
}

func (this *FleetGroup) clean() {

}
