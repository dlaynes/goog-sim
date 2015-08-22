package main

import (
	"./simulator"
	"./tools"

	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	//"log"
	"strconv"
)

func main() {
	var profiler = &tools.Profiler{}
	profiler.Init(40)

	profiler.StartTask("full_simulation")

	profiler.StartTask("load_resources")

	var attackerGroup = &simulator.FleetGroup{}
	var defenderGroup = &simulator.FleetGroup{}

	var resources = map[string]*simulator.Resource{}
	ParseJson("./data/resources.json", &resources)

	var rapidfire = map[string]map[string]float64{}
	ParseJson("./data/rapidfire.json", &rapidfire)

	var attackers = []*simulator.Player{}
	ParseJson("./data/attackers.json", &attackers)

	var defenders = []*simulator.Player{}
	ParseJson("./data/defenders.json", &defenders)

	profiler.EndTask("load_resources")

	//fmt.Printf("%#v\n", resources)
	//fmt.Printf("%#v\n", rapidfire)

	//fmt.Printf("%#v\n", attackers)
	//fmt.Printf("%#v\n", defenders)

	profiler.StartTask("init_groups")

	for _, res := range resources {
		res.Init(rapidfire)
	}

	/*
		profiler.StartTask("calc_fleet_size")
		aSize := 0
		dSize := 0
		for _, pl := range attackers {
			for _, amount := range pl.OriginalFleet {
				aSize += amount
			}
		}
		fmt.Println("Size attackers" + strconv.Itoa(aSize))

		for _, pl2 := range defenders {
			for _, amount2 := range pl2.OriginalFleet {
				dSize += amount2
			}
		}
		fmt.Println("Size defenders" + strconv.Itoa(dSize))
		profiler.EndTask("calc_fleet_size")
	*/

	profiler.StartTask("init_attackers")
	attackerGroup.Init()
	for _, player := range attackers {
		player.Expand(attackerGroup, resources)
	}
	profiler.EndTask("init_attackers")

	profiler.StartTask("init_defenders")
	defenderGroup.Init()
	for _, player2 := range defenders {
		player2.Expand(defenderGroup, resources)
	}
	profiler.EndTask("init_defenders")

	attackerGroup.CalcStatistics(0)
	defenderGroup.CalcStatistics(0)

	profiler.EndTask("init_groups")

	/* Battle */
	simulator.SeedRand()

	profiler.StartTask("full_battle")
	loops := 6
	idx := ""

	//battleStatus := 0
	exitBattle := false
	//successAt := true
	//successDf := true

	for i := 0; i < loops; i++ {
		idx = strconv.Itoa(i + 1)

		if len(attackerGroup.Ships) < 1 {
			fmt.Println("Attacker group has no remaining ships in battle")
			exitBattle = true
		}
		if len(defenderGroup.Ships) < 1 {
			fmt.Println("Defender group has no remaining ships in battle")
			exitBattle = true
		}

		if exitBattle {
			break
		}
		profiler.StartTask("full_round_" + idx)

		profiler.StartTask("battle_round_att_" + idx)
		_ = attackerGroup.Attack(defenderGroup)
		profiler.EndTask("battle_round_att_" + idx)

		profiler.StartTask("battle_round_df_" + idx)
		_ = defenderGroup.Attack(attackerGroup)
		profiler.EndTask("battle_round_df_" + idx)

		profiler.StartTask("battle_clean_att_" + idx)
		attackerGroup.Clean()
		profiler.EndTask("battle_clean_att_" + idx)

		profiler.StartTask("battle_clean_df_" + idx)
		defenderGroup.Clean()
		profiler.EndTask("battle_clean_df_" + idx)

		//fmt.Println("Round " + idx + " has ended")

		profiler.EndTask("full_round_" + idx)

		attackerGroup.CalcStatistics(i + 1)
		defenderGroup.CalcStatistics(i + 1)

	}
	profiler.EndTask("full_battle")

	//fmt.Printf("%#v\n", attackerGroup.Ships)
	//fmt.Printf("%#v\n", defenderGroup.Ships)

	profiler.EndTask("full_simulation")

	/* Results */

	tasks := profiler.GetTasks()
	for n := 0; n < len(tasks); n++ {
		theTask := tasks[n]
		fmt.Printf("Task "+theTask.Label+" took %v \n", theTask.EndTime.Sub(theTask.StartTime))
	}

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
