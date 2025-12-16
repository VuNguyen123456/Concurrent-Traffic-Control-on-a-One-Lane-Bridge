// Vu Nguyen: G01390056
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TrafficPatterns:
//		Car: 100 unit
//		Truck: 200 unit
// 		To cross the bridge each took 2 second

// Bridge Restrictions:
// 		No collision (1 direction on bridge and),
// 		No collapse total on bridge < 750

// Traffic Control Policies:
//		Once traffic flow in 1 direction all other vehicle arrive at that direction have priorities
// 		If there's avehocle waiting at 1 end of bridge then only 6 consecutive behicle may cross then it can go
// 		Bridge capa must be fully used when there's sufficient traffic
//		No traffic deadlock
//		Vehicle arrive at bridege should queue up in single line (keep your own Queue and have P associate with it's own condition in Q => When consumer send signal then only head of Q will wake up)

// Each vehicles is a Threads; 1 Vehicale = 1 goroutin
// Cannot use broadcast

// Global val: All access to this is critical section => modifying them required locking
var BridgeWeight int = 0
var currentFlow string
var consecutive int = 0
var VehiclesOnBridge []Vehicle // Keep track of vehicle on bridge with their ID
var NorthQueue []Vehicle
var SouthQueue []Vehicle

type Vehicle struct {
	vehicle_id   int
	vehicle_type string
	direction    string
	Cond         *sync.Cond
}

func firstInQ(vehicle Vehicle, Q []Vehicle) bool {
	if len(Q) > 0 && Q[0].vehicle_id == vehicle.vehicle_id {
		return true
	}
	return false
}

func carWeight(vehicle Vehicle) int {
	if vehicle.vehicle_type == "Car" {
		return 100
	}
	return 200
}

// func get2QueueFromVehicle(vehicle Vehicle) ([]Vehicle, []Vehicle, string, string) {
// 	var opposingQ []Vehicle
// 	var thisQ []Vehicle
// 	direction := ""
// 	opposeDirection := ""
// 	if vehicle.direction == "North" {
// 		opposingQ = SouthQueue
// 		thisQ = NorthQueue
// 		direction = "North"
// 		opposeDirection = "South"
// 	} else if vehicle.direction == "South" {
// 		opposingQ = NorthQueue
// 		thisQ = SouthQueue
// 		direction = "South"
// 		opposeDirection = "North"
// 	} else {
// 		direction = ""
// 		opposeDirection = ""
// 	}
// 	return thisQ, opposingQ, direction, opposeDirection
// }

// Each of this function will take a thread as parameter right?
func Arrive(vehicle Vehicle, mu *sync.Mutex) {
	// Check all traffic and bridge restriction
	// Put into queue
	// Wait if not allowed to cross
	mu.Lock()
	defer mu.Unlock()

	var opposingQ *[]Vehicle
	var thisQ *[]Vehicle
	if vehicle.direction == "North" {
		NorthQueue = append(NorthQueue, vehicle)
		opposingQ = &SouthQueue
		thisQ = &NorthQueue
		// fmt.Printf("DEBUG: Added to NorthQueue. North size: %d, South size: %d\n", len(NorthQueue), len(SouthQueue))
	} else {
		SouthQueue = append(SouthQueue, vehicle)
		opposingQ = &NorthQueue
		thisQ = &SouthQueue
		// fmt.Printf("DEBUG: Added to SouthQueue. North size: %d, South size: %d\n", len(NorthQueue), len(SouthQueue))
	}

	fmt.Printf("Vehicle #%d (%sbound, Type: %s) arrived.\n",
		vehicle.vehicle_id, vehicle.direction, vehicle.vehicle_type)

	weight := carWeight(vehicle)
	// Wait until this vehicle can cross
	// Wait if bridge is overweight.
	// Wait if flow is opposite.
	// Wait if not first in your queue.
	// Wait if your direction has already had 6 consecutive and opposing queue is non-empty AND bridge is empty.
	for BridgeWeight+weight >= 750 ||
		(currentFlow != "" && currentFlow != vehicle.direction) ||
		(len(*opposingQ) != 0 && consecutive >= 6) ||
		!firstInQ(vehicle, *thisQ) {
		vehicle.Cond.Wait()
	}
	// for BridgeWeight+weight >= 750 || (currentFlow != "" && currentFlow != vehicle.direction) || (len(*opposingQ) != 0 && consecutive >= 6 && len(VehiclesOnBridge) == 0) || !firstInQ(vehicle, *thisQ) {
	// 	vehicle.Cond.Wait()
	// }
	// Ready to cross
}

func printVehiclesOnBridge() string {
	if len(VehiclesOnBridge) == 0 {
		return "None"
	}
	fmt.Printf("[")
	for i := 0; i < len(VehiclesOnBridge); i++ {
		fmt.Printf(" Vehicle #%d (Type: %s) is on the bridge (%s). | ",
			VehiclesOnBridge[i].vehicle_id, VehiclesOnBridge[i].vehicle_type, VehiclesOnBridge[i].direction)
	}
	fmt.Printf("]\n")
	return ""
}

func print2Queues() {
	fmt.Printf("\nWaiting vehicles (Northbound): [")
	for i := 0; i < len(NorthQueue); i++ {
		fmt.Printf(" Vehicle #%d (Type: %s) is waiting in the North queue. | ", NorthQueue[i].vehicle_id, NorthQueue[i].vehicle_type)
	}
	fmt.Printf("]\n")
	fmt.Printf("\nWaiting vehicles (Southbound): [")
	for i := 0; i < len(SouthQueue); i++ {
		fmt.Printf(" Vehicle #%d (Type: %s) is waiting in the South queue. | ", SouthQueue[i].vehicle_id, SouthQueue[i].vehicle_type)
	}
	fmt.Printf("]\n")
}

func Cross(vehicle Vehicle, mu *sync.Mutex) {
	// Delay vehicle 2 sec to imitate crossing
	mu.Lock()
	// if vehicle.Cond == nil {
	// 	vehicle.Cond = sync.NewCond(mu)
	// }
	// Update traffic mechanism
	if len(NorthQueue) > 0 && vehicle.direction == "North" {
		NorthQueue = NorthQueue[1:]
	} else if len(SouthQueue) > 0 && vehicle.direction == "South" {
		SouthQueue = SouthQueue[1:]
	}
	BridgeWeight += carWeight(vehicle)
	VehiclesOnBridge = append(VehiclesOnBridge, vehicle)
	consecutive += 1
	if currentFlow == "" {
		currentFlow = vehicle.direction
	}
	fmt.Printf("Vehicle #%d is now crossing the bridge.\n", vehicle.vehicle_id)
	fmt.Printf("\nVehicles on the bridge: Going %s; Total Weight: %d\n", currentFlow, BridgeWeight)
	fmt.Printf("\t%s", printVehiclesOnBridge())
	print2Queues()
	// fmt.Printf("\nWaiting vehicles (Southbound): %s\n", printQueue(SouthQueue))
	// fmt.Printf("\nWaiting vehicles (Northbound): %s\n", printQueue(NorthQueue))
	mu.Unlock()
	// Consecutive check will be updated on Leave
	time.Sleep(2 * time.Second)
}

// WRONG? !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func wakeUpVehicles(Q *[]Vehicle) {
	for i := 0; i < len(*Q); i++ {
		(*Q)[i].Cond.Signal()
	}
}

func Leave(vehicle Vehicle, mu *sync.Mutex) {
	// Update traffic mechanism and other vehicle cross the bridge
	// Remove vehicle from bridge?
	// Wake up
	mu.Lock()
	defer mu.Unlock()

	// Remove that vehicle on bridge
	for i := 0; i < len(VehiclesOnBridge); i++ {
		if VehiclesOnBridge[i].vehicle_id == vehicle.vehicle_id {
			VehiclesOnBridge = append(VehiclesOnBridge[:i], VehiclesOnBridge[i+1:]...)
			break
		}
	}

	BridgeWeight -= carWeight(vehicle)
	// var opposingQ []Vehicle
	// var thisQ []Vehicle
	// newDir := ""
	// if vehicle.direction == "North" {
	// 	opposingQ = SouthQueue
	// 	thisQ = NorthQueue
	// 	newDir = "South"
	// } else if vehicle.direction == "South" {
	// 	opposingQ = NorthQueue
	// 	thisQ = SouthQueue
	// 	newDir = "North"
	// } else {
	// 	newDir = ""
	// }

	// This part might be wrong because consecutive seems not right ASK TA
	// if len(VehiclesOnBridge) == 0 {
	// 	if consecutive >= 6 && len(opposingQ) > 0 {
	// 		currentFlow = newDir
	// 		consecutive = 0
	// 	} else if len(thisQ) > 0 {
	// 		currentFlow = vehicle.direction
	// 	} else {
	// 		currentFlow = ""
	// 		consecutive = 0
	// 	}
	// }

	fmt.Printf("\nVehicle #%d exited the bridge. Total Weight: %d\n", vehicle.vehicle_id, BridgeWeight)

	if len(VehiclesOnBridge) == 0 {
		// Bridge is empty → decide who goes next
		northWaiting := len(NorthQueue) > 0
		southWaiting := len(SouthQueue) > 0

		if currentFlow == "North" {
			if consecutive >= 6 && southWaiting {
				// Switch to South after 6
				currentFlow = "South"
				consecutive = 0
				wakeUpVehicles(&SouthQueue)
			} else if northWaiting {
				// Continue North , keep consecutive so when south came up it can switch after 6
				wakeUpVehicles(&NorthQueue)
			} else if southWaiting {
				// Switch To south if North empty (this is for when north is flowing but not enough cars)
				currentFlow = "South"
				consecutive = 0
				wakeUpVehicles(&SouthQueue)
			} else {
				// nobody waiting, might not need to reset consecutive?????? HEYYYYYYYYYYYYYYYYYYYYYYYYYYYY LOOOKKKKKKKKKK HEEEEEEEEERRRRRRRRRRREEEEEEEEEEEEE AGAINNNNNNNNNN
				currentFlow = ""
				consecutive = 0
			}
			// Opposite direction but like above
		} else if currentFlow == "South" {
			if consecutive >= 6 && northWaiting {
				currentFlow = "North"
				consecutive = 0
				wakeUpVehicles(&NorthQueue)
			} else if southWaiting {
				wakeUpVehicles(&SouthQueue)
			} else if northWaiting {
				currentFlow = "North"
				consecutive = 0
				wakeUpVehicles(&NorthQueue)
			} else {
				currentFlow = ""
				consecutive = 0
			}
		} else {
			// No current traffic flow yet → pick whoever is waiting
			if northWaiting {
				currentFlow = "North"
				wakeUpVehicles(&NorthQueue)
			} else if southWaiting {
				currentFlow = "South"
				wakeUpVehicles(&SouthQueue)
			}
		}
	} else {
		// Bridge not empty, keep going same direction if possible (There are still car left in that direction and consecutive still allow)
		if currentFlow == "North" && len(NorthQueue) > 0 && consecutive < 6 {
			wakeUpVehicles(&NorthQueue)
		} else if currentFlow == "South" && len(SouthQueue) > 0 && consecutive < 6 {
			wakeUpVehicles(&SouthQueue)
		}
	}

	// // Wake up (It's not hapening?)
	// // Any vehicle waiting
	// if len(VehiclesOnBridge) == 0 {
	// 	if currentFlow == "North" && len(NorthQueue) > 0 {
	// 		wakeUpVehicles(&NorthQueue)
	// 	} else if currentFlow == "South" && len(SouthQueue) > 0 {
	// 		wakeUpVehicles(&SouthQueue)
	// 	}
	// } else if consecutive == 0 {
	// 	if vehicle.direction == "North" && len(SouthQueue) > 0 {
	// 		wakeUpVehicles(&SouthQueue)
	// 	} else if vehicle.direction == "South" && len(NorthQueue) > 0 {
	// 		wakeUpVehicles(&NorthQueue)
	// 	}
	// 	// Wake up same direction if there's any
	// } else {
	// 	if vehicle.direction == "North" && len(NorthQueue) > 0 {
	// 		wakeUpVehicles(&NorthQueue)
	// 	} else if vehicle.direction == "South" && len(SouthQueue) > 0 {
	// 		wakeUpVehicles(&SouthQueue)
	// 	}
	// }
}

func OneVehicle(vehicle Vehicle, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as complete when done
	Arrive(vehicle, mu)
	Cross(vehicle, mu)
	Leave(vehicle, mu)
}

// OneVehicle(parameter-list)
//  {
//  Arrive(. . .)
//  Cross(. . .)
//  Leave(. . .)
//  }

type VehicleGroup struct {
	numVehicles int
	delayTime   int // in seconds
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	vehicleId := 0
	//Assuming car will get generate and provided here (go ....)

	fmt.Print("How many groups of vehicles do you want to simulate? ")
	var numGroups int
	fmt.Scanln(&numGroups)

	groups := make([]VehicleGroup, numGroups)
	for i := 0; i < numGroups; i++ {
		fmt.Printf("Enter number of vehicles in group %d: ", i+1)
		fmt.Scanln(&groups[i].numVehicles)
		// Delay is ask after each group except the last one
		if i < numGroups-1 {
			fmt.Printf("Enter delay (seconds) before group %d: ", i+2)
			fmt.Scanln(&groups[i].delayTime)
		}
	}

	rand.Seed(time.Now().UnixNano())
	fmt.Print("Enter Car probability (0.0–1.0): ")
	var probCar float64
	fmt.Scanln(&probCar)

	fmt.Print("Enter Northbound probability (0.0–1.0): ")
	var probNorth float64
	fmt.Scanln(&probNorth)
	fmt.Print("\n\n")
	for i := 0; i < numGroups; i++ {
		// Simulate a group of vehicles
		numVehiclesInGroup := groups[i].numVehicles

		// Genrate vehicles with random type and direction
		for j := 0; j < numVehiclesInGroup; j++ {
			// Randomly generate vehicle type and direction
			vehicleType := "Truck"
			if rand.Float64() < probCar {
				vehicleType = "Car"
			}

			direction := "South"
			if rand.Float64() < probNorth {
				direction = "North"
			}

			vehicle := Vehicle{
				vehicle_id:   vehicleId,
				vehicle_type: vehicleType,
				direction:    direction,
				Cond:         sync.NewCond(&mu),
			}
			vehicleId++
			wg.Add(1)
			go OneVehicle(vehicle, &mu, &wg)
		}
		// Delay before next group
		if i < numGroups-1 {
			time.Sleep(time.Duration(groups[i].delayTime) * time.Second)
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()
}
