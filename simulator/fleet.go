// Copyright 2015 Donato Cassel Laynes Gonzales
//
// This file is part of GoOgame - Battle Simulator

package simulator

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Player struct {
	Id            int
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
	Rapidfires map[int]float64
}

type ShipUnit struct {
	T *ShipType
	A float64
	S float64
	H float64
	//U bool
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

func (this *Player) Expand(group *FleetGroup, resources map[string]*Resource) {
	this.processShipTypes(resources)

	ships := &group.Ships

	//fmt.Println("Expandiendo")

	for _, t := range this.ShipTypes {
		*ships = append(*ships, t.Expand()...)
	}
}

func (this *Player) ExpandTo(group *FleetGroup, resources map[string]*Resource) {
	this.processShipTypes(resources)

	for _, t := range this.ShipTypes {
		t.ExpandTo(group)
	}
}

func (this *Player) processShipTypes(resources map[string]*Resource) {

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

func (this *ShipType) Init(resource *Resource, amount int, mtech int, dtech int, htech int) {
	this.Amount = amount

	d := 1 + (float64(mtech) * 0.1)
	s := 1 + (float64(dtech) * 0.1)
	h := (1 + (float64(htech) * 0.1)) * 0.1

	this.Res = resource
	this.BaseAttack = d * resource.Attack
	this.BaseShield = s * resource.Defense
	this.BaseHull = h * resource.Hull

	//fmt.Println("Looking for Rapidire values on " + resource.Id)
	//fmt.Printf("%#v\n", resource)

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

func (this *ShipType) ExpandTo(group *FleetGroup) {

	ships := group.Ships
	l := len(ships)

	fmt.Println("Len " + strconv.Itoa(l))
	fmt.Println("Capacity " + strconv.Itoa(cap(ships)))
	fmt.Println("Amount " + strconv.Itoa(this.Amount))

	for i := l; i < this.Amount+l; i++ {
		//fmt.Println("Appending to " + strconv.Itoa(i))
		ships = append(ships, &ShipUnit{T: this, A: this.BaseAttack, S: this.BaseShield, H: this.BaseHull})
	}
	group.Ships = ships
}

/* TO DO */
func (this *ShipType) CalcCapacity() {

}

func (this *ShipType) LogBattle() {

}

/* FleetGroup functions */
func (this *FleetGroup) Init() {
	this.Ships = make([]*ShipUnit, 0)
}

func (this *FleetGroup) InitWith(length int) {
	this.Ships = make([]*ShipUnit, 0, length)
}

func (this *FleetGroup) Attack(otherGroup *FleetGroup) bool {

	m := len(otherGroup.Ships)
	var Dm, Dc, De, xp float64
	var c int
	var fPtr, uPtr *ShipUnit

	running := true

	this.TurnAttacks = len(this.Ships)
	c = this.TurnAttacks

	//r := &rand.Rand

	//function pointers
	ri := rand.Intn
	rf := rand.Float64

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

			uPtr = otherGroup.Ships[ri(m)]

			if uPtr.H != 0.0 {
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
							if xp > 0.3 && rf() < xp {
								//boom!
								uPtr.H = 0.0
								uPtr.T.Explosions += 1
							}

						} else {
							//Kaboom!
							uPtr.H = 0.0
							uPtr.T.Explosions += 1
						}

					} else {
						uPtr.S -= Dm           // The shield survived the shot. We decrease the shield points of the target
						this.TurnDefense += Dm // We update the shield protection statistics
					}
				} else {
					this.TurnDefense += Dm
					running = false
				}
			}

			//fmt.Printf("%#v\n", fPtr.T.Rapidfires)

			//Rapidfire calculations
			if val, ok := fPtr.T.Rapidfires[uPtr.T.Res.Id]; ok {
				//Do we get another turn?

				//fmt.Println("Rapidfire value " + strconv.FormatFloat(val, 'g', 1, 64))

				if rf() < val {
					this.TurnAttacks++
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

/*
func (this *FleetGroup) AttackParallel(otherGroup *FleetGroup) bool {
	m := len(otherGroup.Ships)
	//var Dm, Dc, De, xp float64
	//var c int
	//var fPtr, uPtr *ShipUnit

	//running := true

	this.TurnAttacks = len(this.Ships)
	c := this.TurnAttacks

	sem := make(chan bool, c)

	for i := 0; i < c; i++ {
		go func(fPtr *ShipUnit) {

			Dm := fPtr.T.BaseAttack
			Dc := Dm * 100.0
			running := true

			for running {
				this.TurnDamage += Dm //We shoot! and we update the statistics accordingly

				uPtr := otherGroup.Ships[rand.Intn(m)]

				if uPtr.H != 0.0 {
					//Check if the shot is strong enough against Large Shield Domes
					if Dc > uPtr.T.BaseShield {
						if Dm > uPtr.S {
							//Shield wasn't strong enough to survive the shot
							De := Dm - uPtr.S          //New damage after substracting the shield points
							this.TurnDefense += uPtr.S //shield protection statistics
							uPtr.S = 0.0

							//Check if the ships "health" is greater than the damage
							if De < uPtr.H {
								uPtr.H -= De

								xp := (uPtr.T.BaseHull - uPtr.H) / uPtr.H
								if xp > 0.3 && rand.Float64() < xp {
									//boom!
									uPtr.H = 0.0
									uPtr.T.Explosions += 1
								}

							} else {
								//Kaboom!
								uPtr.H = 0.0
								uPtr.T.Explosions += 1
							}

						} else {
							uPtr.S -= Dm           // The shield survived the shot. We decrease the shield points of the target
							this.TurnDefense += Dm // We update the shield protection statistics
						}
					} else {
						this.TurnDefense += Dm
						running = false
					}
				}

				//fmt.Printf("%#v\n", fPtr.T.Rapidfires)

				//Rapidfire calculations
				if val, ok := fPtr.T.Rapidfires[uPtr.T.Res.Id]; ok {
					//Do we get another turn?

					//fmt.Println("Rapidfire value " + strconv.FormatFloat(val, 'g', 1, 64))

					if rand.Float64() < val {
						this.TurnAttacks++
					} else {
						running = false
					}

				} else {
					running = false
				}
			}

			//end goroutine
			sem <- true
		}(this.Ships[i])
	}

	//Huh
	for i := 0; i < c; i++ {
		<-sem
	}

	return true
}
*/

func (this *FleetGroup) Clean() {

	fmt.Println("Length before cleanup " + strconv.Itoa(len(this.Ships)))

	//Adding the approximate capacity. function is now 4x/5x faster
	//However when applying this patch elsewhere, this makes the init_groups section 10x slower, so...
	newShips := make([]*ShipUnit, 0, len(this.Ships))

	c := len(this.Ships)
	for i := 0; i < c; i++ {
		ship := this.Ships[i]
		if ship.H != 0.0 {
			ship.S = ship.T.BaseShield
			newShips = append(newShips, ship)
		} else {
			ship = nil
		}
	}
	this.Ships = newShips

	fmt.Println("Length after cleanup  " + strconv.Itoa(len(this.Ships)))
}

func (this *FleetGroup) CalcStatistics(round int) {
	//TO DO: generate statistics report

	this.TurnDamage = 0.0
	this.TurnDefense = 0.0
}
