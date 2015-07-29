package main

import (
	"./simulator"
	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	//"log"
)

func main() {
	var attackerGroup = &simulator.FleetGroup{}
	var defenderGroup = &simulator.FleetGroup{}

	var resources = map[string]simulator.Resource{}
	err := json.Unmarshal(ReadFile("./data/resources.json"), &resources)
	check(err)

	var rapidfire = map[string]map[string]float64{}
	err2 := json.Unmarshal(ReadFile("./data/rapidfire.json"), &rapidfire)
	check(err2)

	var attackers = map[string]simulator.Player{}
	err3 := json.Unmarshal(ReadFile("./data/attackers.json"), &attackers)
	check(err3)

	var defenders = map[string]simulator.Player{}
	err4 := json.Unmarshal(ReadFile("./data/defenders.json"), &defenders)
	check(err4)

	//log.Println(resources)
	//log.Println(rapidfire)

	//fmt.Printf("%#v\n", resources)
	//fmt.Printf("%#v\n", rapidfire)

	//fmt.Printf("%#v\n", attackers)
	//fmt.Printf("%#v\n", defenders)

	for _, res := range resources {
		res.Init(rapidfire)
	}

	for _, player := range attackers {
		player.Expand(attackerGroup, &resources)
	}
	for _, player2 := range defenders {
		player2.Expand(defenderGroup, &resources)
	}

	fmt.Printf("Ok!.\n")
}

func ReadFile(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	check(err)

	//log.Println(json.Unmarshal(dat))

	return dat
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
