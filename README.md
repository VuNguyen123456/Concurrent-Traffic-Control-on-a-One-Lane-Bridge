# Concurrent Traffic Control on a One-Lane Bridge (Go)

This project implements a **concurrent traffic control system** in Go that simulates vehicles crossing an imaginary **one-lane bridge**. Each vehicle is modeled as a **goroutine**, and synchronization is enforced using **mutex locks and condition variables** to guarantee safety, fairness, and efficiency.

The goal is to gain hands-on experience with **concurrent programming**, **thread synchronization**, and **condition variables** using Go.

---

## Problem Overview

The bridge supports traffic traveling in **two directions** (northbound and southbound) but has only **one lane**, meaning traffic may flow in **only one direction at a time**.

Vehicles consist of:
- **Cars** (weight = 100)
- **Trucks** (weight = 200)

Each vehicle takes **2 seconds** to cross the bridge.

---

## Bridge Restrictions

The system must enforce the following safety constraints:

- **R1 – Direction Restriction:**  
  Vehicles may only cross in one direction at a time (to prevent collisions).

- **R2 – Weight Restriction:**  
  The total weight of vehicles on the bridge must not exceed **750 units** (to prevent collapse).

Violating either restriction is considered a system failure.

---

## Traffic Control Policies

The bridge controller enforces the following policies:

- **P1 – Direction Priority:**  
  Once traffic starts flowing in one direction, newly arriving vehicles in the same direction are prioritized.

- **P2 – Fairness:**  
  If vehicles are waiting on one side, **no more than 6 consecutive vehicles** may cross from the opposite side before allowing the waiting side to proceed.

- **P3 – Maximum Utilization:**  
  Subject to safety and fairness, the bridge capacity must be fully utilized when traffic is available.  
  Deadlocks must be avoided.

- **P4 – FIFO Queuing:**  
  Vehicles queue in arrival order on each side and must enter the bridge in that same order.  
  If a truck at the head of the queue cannot enter due to weight limits, vehicles behind it must also wait.

---

## Concurrency Model

Each vehicle is represented by a **goroutine** executing the following lifecycle:

```go
OneVehicle(vehicleID, vehicleType, direction) {
    Arrive()
    Cross()

Function Responsibilities

Arrive()

Blocks until it is safe to enter the bridge

Enforces all restrictions and policies

Uses mutex locks and condition variables

Cross()

Simulates crossing the bridge by sleeping for 2 seconds

Leave()

Updates shared state

Signals waiting vehicles when appropriate

All shared state is protected inside critical sections, and busy waiting is strictly avoided.
    Leave()
}

Output and Logging

The program prints detailed runtime information, including:

Vehicle arrivals and departures

Vehicles currently crossing the bridge

Total weight on the bridge

Waiting queues on both sides (in order)

Example output format:

Vehicle #5 (Northbound, Type: Car) arrived.
Vehicle #5 is now crossing the bridge.
Vehicles on the bridge: [#3 (Northbound, Truck), #5 (Northbound, Car)]
Total Weight: 300
Waiting vehicles (Southbound): [#7 (Truck), #8 (Car)]
Vehicle #5 exited the bridge.

Arrival Schedules

The program is executed under five predefined arrival schedules, each with different traffic patterns and delays.
Each schedule introduces 30 vehicles total, with randomized direction and vehicle type based on given probabilities.

Examples include:

Car-only traffic

Truck-heavy traffic

Bursty arrivals with delays

Mixed traffic patterns

Vehicle direction is chosen with equal probability (50% northbound, 50% southbound).

Implementation Constraints

Language: Go

Synchronization:

sync.Mutex

sync.Cond

Only Signal() is allowed for waking goroutines (Broadcast() is prohibited)

No busy waiting

All shared variables accessed only within critical sections

FIFO order must be preserved despite Go’s nondeterministic scheduling

Learning Objectives

This project demonstrates:

Practical use of mutexes and condition variables

Correct handling of concurrency and fairness

Avoidance of deadlocks and starvation

Designing concurrent systems under realistic constraints

Modeling real-world resource contention with goroutines

References

Go Concurrency Documentation

Operating Systems: Three Easy Pieces (Chapters 25–30)

Go sync and time packages

Notes

Due to nondeterministic thread scheduling in Go:

Vehicles may not execute or exit in creation order

Correctness is defined by policy enforcement, not execution order

---
