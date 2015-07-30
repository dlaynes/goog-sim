package simulator

import (
	//"fmt"
	"math/rand"
	"time"
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
	T *ShipType
	A float64
	S float64
	H float64
	X bool
}

type FleetGroup struct {
	Ships       []*ShipUnit
	TurnDamage  float64
	TurnDefense float64
	TurnAttacks int
}

/* Globals */

func SeedRand() {
	rand.Seed(time.Now().Unix())
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

	//fmt.Println("Cantidad Items ", len(this.OriginalFleet))

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

	this.Res = &resource
	this.BaseAttack = d * resource.Attack
	this.BaseShield = s * resource.Defense
	this.BaseHull = h * resource.Hull

	this.Rapidfires = resource.Rapidfires
}

func (this *ShipType) Expand() []*ShipUnit {
	sh := make([]*ShipUnit, this.Amount)

	for i := 0; i < this.Amount; i++ {
		sh[i] = &ShipUnit{T: this, A: this.BaseAttack, S: this.BaseShield, H: this.BaseHull}
	}

	//fmt.Printf("%#v\n", sh)

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

func (this *FleetGroup) Attack(otherGroup *FleetGroup) bool {

	m := len(otherGroup.Ships)
	var Dm, Dc, De, xp float64
	var c int
	var fPtr, uPtr *ShipUnit
	var resId string

	running := true

	this.TurnAttacks = len(this.Ships)
	c = this.TurnAttacks

	//r := &rand.Rand

	//TO DO: add concurrency, and use binary operations ...
	for i := 0; i < c; i++ {

		//Init some variables for the current Ship
		fPtr = this.Ships[i]
		Dm = fPtr.T.BaseAttack
		Dc = Dm * 100.0
		running = true

		//Current Ship loop
		for running {
			this.TurnDamage += Dm //We shoot! and we update the statistics accordingly

			uPtr = otherGroup.Ships[rand.Intn(m)]

			if !uPtr.X {
				//Check if the shot is strong enough against Large Shield Domes
				if Dc > uPtr.T.BaseShield {
					if Dm > uPtr.S {
						//Shield wasn't strong enough to survive the shot
						De = Dm - uPtr.S           //New damage after substracting the shield points
						this.TurnDefense += uPtr.S //shield protection statistics
						uPtr.S = 0.0

						//Check if the ships "health" is greater than the damage
						if De < uPtr.H {
							uPtr.H -= De

							xp = (uPtr.T.BaseHull - uPtr.H) / uPtr.H
							if xp > 0.3 && rand.Float64() < xp {
								//boom!
								uPtr.X = true
								uPtr.T.Explosions += 1
							}

						} else {
							//Kaboom!
							uPtr.X = true
							uPtr.T.Explosions += 1
						}

					} else {
						uPtr.S -= Dm           // The shield survived the shot. We decrease the shield points of the target
						this.TurnDefense += Dm // We update the shield protection statistics accordingly
					}
				} else {
					this.TurnDefense += Dm
					running = false
				}
			}

			//Rapidfire calculations
			resId = uPtr.T.Res.Id
			val, ok := fPtr.T.Rapidfires[resId]

			if ok {
				if val > 0.0 {
					//Do we get another turn?
					if rand.Float64() < fPtr.T.Rapidfires[uPtr.T.Res.Id] {
						this.TurnAttacks++
					} else {
						running = false
					}
				} else {
					running = false
				}
			} else {
				running = false
			}
		}
	}

	return true
}

func (this *FleetGroup) Clean() {
	newShips := make([]*ShipUnit, 0)
	for _, ship := range this.Ships {
		if ship.X {
			ship.S = ship.T.BaseShield
			newShips = append(newShips, ship)
		} else {
			ship = nil
		}
	}
	this.Ships = newShips
}

func (this *FleetGroup) StartStatistics() {

}
