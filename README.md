# Concurrent Traffic Control on a One-Lane Bridge (Go)

This project implements a **concurrent traffic control system** in Go that simulates vehicles crossing an imaginary **one-lane bridge**. Each vehicle is modeled as a **goroutine**, and synchronization is enforced using **mutex locks and condition variables** to guarantee safety, fairness, and efficiency.

The purpose of this project is to gain hands-on experience with **concurrent programming**, **thread synchronization**, and **condition variables** in Go, while modeling a realistic shared-resource problem.

---

## Problem Overview

The bridge supports traffic traveling in **two directions** (northbound and southbound) but has only **one lane**, meaning traffic may flow in **only one direction at a time**.

Vehicles consist of:
- **Cars** (weight = 100 units)
- **Trucks** (weight = 200 units)

Each vehicle takes **2 seconds** to cross the bridge.

---

## Bridge Restrictions

The system enforces the following safety constraints:

- **R1 – Direction Restriction**  
  Vehicles may cross the bridge in only one direction at a time to prevent collisions.

- **R2 – Weight Restriction**  
  The total weight of vehicles on the bridge must not exceed **750 units** at any time to prevent collapse.

Violating either restriction is considered a system failure.

---

## Traffic Control Policies

The bridge controller enforces the following policies:

- **P1 – Direction Priority**  
  Once traffic begins flowing in one direction, newly arriving vehicles traveling in the same direction are prioritized.

- **P2 – Fairness**  
  If vehicles are waiting on one side of the bridge, **no more than 6 consecutive vehicles** may cross from the opposite direction before allowing the waiting side to proceed.

- **P3 – Maximum Utilization**  
  Subject to safety and fairness constraints, the bridge should be fully utilized whenever traffic is available.  
  Deadlocks must be avoided.

- **P4 – FIFO Queuing**  
  Vehicles queue in arrival order on each side of the bridge and must enter the bridge in that same order.  
  If a truck at the head of a queue cannot enter due to weight limits, vehicles behind it must also wait.

---

## Concurrency Model
Each vehicle is represented by a **goroutine** executing the following lifecycle:

OneVehicle(vehicleID, vehicleType, direction):
Arrive()
Cross()
Leave()


### Function Responsibilities

**Arrive()**
- Blocks until it is safe for the vehicle to enter the bridge
- Enforces all restrictions and traffic control policies
- Uses mutex locks and condition variables

**Cross()**
- Simulates crossing the bridge by sleeping for 2 seconds

**Leave()**
- Updates shared state
- Signals waiting vehicles when appropriate

All shared state is accessed only inside **critical sections**, and **busy waiting is strictly avoided**.

---

## Output and Logging

The program prints detailed runtime information, including:
- Vehicle arrivals and departures
- Vehicles currently crossing the bridge
- Total weight on the bridge
- Waiting queues on both sides, displayed in FIFO order

Example output:
Vehicle #5 (Northbound, Type: Car) arrived.
Vehicle #5 is now crossing the bridge.
Vehicles on the bridge: [#3 (Northbound, Truck), #5 (Northbound, Car)]
Total Weight: 300
Waiting vehicles (Southbound): [#7 (Truck), #8 (Car)]
Vehicle #5 exited the bridge.

---

## Arrival Schedules

The program is executed using **five predefined arrival schedules**, each with different traffic patterns and delays.

- Each schedule introduces **30 vehicles total**
- Vehicle direction is chosen with **equal probability** (50% northbound, 50% southbound)
- Vehicle type (car or truck) is determined using schedule-specific probabilities

Traffic patterns include:
- Car-only traffic
- Truck-heavy traffic
- Bursty arrivals with delays
- Mixed traffic workloads

---

## Implementation Constraints

- **Language:** Go
- **Synchronization Mechanisms:**
  - `sync.Mutex`
  - `sync.Cond`
- Only `Signal()` may be used to wake blocked goroutines (`Broadcast()` is prohibited)
- Busy waiting is not allowed
- All shared variables are accessed only within critical sections
- FIFO ordering must be preserved despite Go’s nondeterministic scheduling

---

## Learning Objectives

This project demonstrates:
- Practical use of **mutexes and condition variables**
- Safe and fair concurrent resource management
- Avoidance of **deadlocks and starvation**
- Designing concurrent systems under strict constraints
- Modeling real-world contention using goroutines

---

## References

- Go Concurrency Documentation  
- *Operating Systems: Three Easy Pieces*, Chapters 25–30  
- Go `sync` and `time` packages  

---

## Notes

Due to nondeterministic goroutine scheduling in Go:
- Vehicles may not execute or exit in creation order
- Correctness is defined by policy enforcement, not execution order

