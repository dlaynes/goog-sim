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
	ParseJson("./data/resources.json", &resources)

	var rapidfire = map[string]map[string]float64{}
	ParseJson("./data/rapidfire.json", &rapidfire)

	var attackers = map[string]simulator.Player{}
	ParseJson("./data/attackers.json", &attackers)

	var defenders = map[string]simulator.Player{}
	ParseJson("./data/defenders.json", &defenders)

	//fmt.Printf("%#v\n", resources)
	//fmt.Printf("%#v\n", rapidfire)

	//fmt.Printf("%#v\n", attackers)
	//fmt.Printf("%#v\n", defenders)

	attackerGroup.Init()
	defenderGroup.Init()

	for _, res := range resources {
		//id.(string)
		//res.(simulator.Resource)
		res.Init(rapidfire)
	}

	for _, player := range attackers {
		player.Expand(attackerGroup, resources)
	}
	for _, player2 := range defenders {
		player2.Expand(defenderGroup, resources)
	}

	//fmt.Printf("%#v\n", attackerGroup.Ships)
	//fmt.Printf("%#v\n", defenderGroup.Ships)

	fmt.Printf("Ok!.\n")
}

func ParseJson(file string, targetPtr interface{}) {
	err := json.Unmarshal(ReadFile(file), &targetPtr)
	check(err)
}

func ReadFile(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	check(err)

	//fmt.Println(dat)

	return dat
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
