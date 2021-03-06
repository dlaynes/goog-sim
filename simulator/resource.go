// Copyright 2015 Donato Cassel Laynes Gonzales
//
// This file is part of GoOgame - Battle Simulator

package simulator

import (
	//"fmt"
	"strconv"
)

const (
	PosTypePlanet = iota + 1
	PosTypeMoon
	PosTypeDebris
)

const (
	ResTypeShip = iota + 1
	ResTypeDefense
	ResTypeBuilding
	ResTypeMaterial
	ResTypeTech
)

type Position struct {
	Galaxy  int
	System  int
	Planet  int
	PosType int
}

type Resource struct {
	Id                 int
	Name               string
	Metal              int
	Crystal            int
	Deuterium          int
	Size               float64
	Pos                *Position
	Energy             int
	Factor             float64
	Speed              int
	Capacity           int
	Attack             float64
	Defense            float64
	Hull               float64
	Motors             map[string]int
	CurrentMotor       string
	Speeds             map[string]int
	CurrentSpeed       string
	Consumptions       map[string]float64
	CurrentConsumption float64

	Rapidfires map[int]float64
	ResType    int
}

func (this *Resource) InitPlanet(g int, s int, p int, t int) {
	this.Pos = &Position{Galaxy: g, System: s, Planet: p, PosType: t}
}

func (this *Resource) Init(rapidfire map[string]map[string]float64) {
	//id := int(this.Id)

	this.initRapidfire(rapidfire)
}

func (this *Resource) initRapidfire(rapidfire map[string]map[string]float64) {
	rf := rapidfire[strconv.Itoa(this.Id)]
	ln := len(rf)
	this.Rapidfires = make(map[int]float64)

	//

	if ln > 0 {
		for id, r := range rf {
			i, _ := strconv.Atoi(id)
			this.Rapidfires[i] = (1.0 - (1.0 / r))

			//fmt.Println("Adding a new rapidfire rule: " + strconv.FormatFloat(this.Rapidfires[i], 'g', 1, 64))
		}
	}

	//fmt.Printf("%#v\n", this.Rapidfires)
}
