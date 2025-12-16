# README.txt

## Source Files
- lab1.go – Contains the main program logic for simulating vehicles crossing the one-lane bridge.

## Compilation & Running
To compile the program, run:
    go run .\lab1.go

## User Input
When you run the program, it will ask you step by step:

1. Number of Groups  
   How many sets of vehicles you want to simulate.  
   Example: 3

2. Number of Vehicles in Each Group  
   For each group, how many vehicles should appear.  
   Example: Group 1 = 5, Group 2 = 4, Group 3 = 6

3. Delay Between Groups (seconds)  
   After each group (except the last one), how many seconds to wait before generating the next group.  
   Example: Delay before Group 2 = 2, Delay before Group 3 = 3

4. Car Probability (0.0 – 1.0)  
   Probability that a vehicle is a Car (otherwise it’s a Truck).  
   Example: 0.7  → 70% chance a vehicle is a Car

5. Northbound Probability (0.0 – 1.0)  
   Probability that a vehicle goes North (otherwise South).  
   Example: 0.5  → 50% chance North, 50% chance South

---

## Example Run
Input:
    How many groups of vehicles do you want to simulate? 2
    Enter number of vehicles in group 1: 3
    Enter delay (seconds) before group 2: 5
    Enter number of vehicles in group 2: 2
    Enter Car probability (0.0–1.0): 0.6
    Enter Northbound probability (0.0–1.0): 0.5

Explanation of input:
- 2 groups total
- Group 1 has 3 vehicles
- Wait 5 seconds before creating Group 2
- Group 2 has 2 vehicles
- 60% chance each vehicle is a Car, otherwise a Truck
- 50% chance each vehicle is Northbound, otherwise Southbound

---

## Example Run 2:
Input:
    How many groups of vehicles do you want to simulate? 3
    Enter number of vehicles in group 1: 3
    Enter delay (seconds) before group 2: 5
    Enter number of vehicles in group 2: 2
    Enter delay (seconds) before group 3: 3
    Enter number of vehicles in group 3: 4
    Enter Car probability (0.0–1.0): 0.2
    Enter Northbound probability (0.0–1.0): 0.5
