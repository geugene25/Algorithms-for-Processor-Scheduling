/*“I Greg Eugene affirm that this program is entirely my own work and
that I have neither developed my code together with any another person, nor copied any code from any
other person, nor permitted my code to be copied or otherwise used by any other person, nor have I copied,
modified, or otherwise used programs created by others. I acknowledge that any violation of the above terms
will be treated as academic dishonesty.”
*/

//scheduling algs FCFS SJFP RR

//go imports
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Process struct {
	name       string //proc name
	PID        int    //proc name for sorting
	arrival    int    //arrival time
	wait       int    // wait time
	burst      int    //burst time
	selectB    bool   //iff selected true else false
	finished   bool   //iff done running true else false
	selected   int
	turnaround int //turnaround
	complete   int //completed time
}

var (
	processcount int       //num of procs
	runfor       int       //total time prog runs
	use          string    //which string to use
	quantum      int       //for RR
	procs        []Process //slice of structs
)

//main fucntion where readin and stuff occurs
func main() {
	//os.Args[1] is cmd line file
	//int slice unspecified first time wait will be 0

	fileRead := os.Args[1] //read the file in system
	fileOut := os.Args[2]
	file, _ := os.Open(fileRead) //open file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) //scan in file as words

	for scanner.Scan() {
		checkStr := scanner.Text()
		//all info save use will need to have text converted to string
		if checkStr == "processcount" {
			scanner.Scan()
			processcount, _ = strconv.Atoi(scanner.Text())
		} else if checkStr == "runfor" {
			scanner.Scan()
			runfor, _ = strconv.Atoi(scanner.Text())
		} else if checkStr == "use" {
			scanner.Scan()
			use = scanner.Text()
			if use == "rr" {
				scanner.Scan()
				quantum, _ = strconv.Atoi(scanner.Text())
				fmt.Printf("%3d", quantum)
			} else {
				continue
			} //if not rr quant is 0
		}
	}
	//initalize the Process struct values
	procs = initStruct(fileRead, processcount)
	switch use {
	case "fcfs":
		{
			FCFS(procs, processcount, runfor, fileOut)
			break
		}
	case "sjf":
		{
			SJFP(procs, processcount, runfor, fileOut)
			break
		}
	case "rr":
		{
			//RR(prcos, processcount,runfor)
			break
		}

	}
}

//AHHHHHHH IT WORKS
func FCFS(procs []Process, processcount int, runfor int, output string) {
	fileOut, _ := os.Create(output)
	arrivalQ := make([]Process, processcount) //make a queue like list for arrivals
	amount := 0                               //amount items in queue
	pos := 0                                  //position in list
	time := 0
	procs = selSort("Arrival", procs)

	fmt.Fprintf(fileOut, "%3d processes\n", processcount)
	fmt.Fprintf(fileOut, "Using First-Come First-Served\n")

	for time < runfor { //run time is finished
		for i := 0; i < processcount; i++ {
			if procs[i].arrival == time { //iff process just arrived
				fmt.Fprintf(fileOut, "Time %3d : %s arrived\n", time, procs[i].name)
				arrivalQ[i] = procs[i] //add current process to queue
				amount++               // added to queue knows it has things to run
			}
		}

		if amount == 0 { //if queue empty then its idle
			fmt.Fprintf(fileOut, "Time %3d : Idle\n", time)
		}
		if amount > 0 { //when queue not empty
			//iff procs selected and its selected time + burst time is current time
			if arrivalQ[pos].selectB && (arrivalQ[pos].selected+arrivalQ[pos].burst == time) {
				arrivalQ[pos].finished = true //its done
				arrivalQ[pos].complete = time //at this time
				arrivalQ[pos].selectB = false //dequeue
				amount--                      //one item out of queue
				fmt.Fprintf(fileOut, "Time %3d : %s finished\n", time, procs[pos].name)

				if pos < (processcount - 1) {
					pos++ //moves on to next item in the queue
				}
			}

			//iff not selected and not done then select next one
			if !arrivalQ[pos].selectB && !arrivalQ[pos].finished && amount > 0 {
				arrivalQ[pos].selectB = true
				arrivalQ[pos].selected = time
				fmt.Fprintf(fileOut, "Time %3d : %s selected (burst %3d)\n", time, procs[pos].name, procs[pos].burst)
			} else if amount == 0 {
				fmt.Fprintf(fileOut, "Time %3d : Idle\n", time)
			}

		}
		time++ //keep time rolling
	}

	//time calculations
	procs = arrivalQ
	for i := 0; i < processcount; i++ {
		procs[i].turnaround = procs[i].complete - procs[i].arrival
		procs[i].wait = procs[i].turnaround - procs[i].burst
	}

	fmt.Fprintf(fileOut, "Finished at time  %d\n\n", runfor)
	procs = selSort("PID", procs)

	for i := 0; i < processcount; i++ {
		fmt.Fprintf(fileOut, "%s wait %3d turnaround %3d\n", procs[i].name, procs[i].wait, procs[i].turnaround)
	}

}

//comment
func SJFP(procs []Process, processcount int, runfor int, output string) {
	fileOut, _ := os.Create(output)
	procs = selSort("Arrival", procs) // sort arrivals
	fmt.Fprintf(fileOut, "%3d process\n", processcount)
	fmt.Fprintf(fileOut, "Using preemptive Shortest Job First\n")

	arrivalQ := make(([]Process), 0, processcount)
	amount := 0
	time := 0

	for time < runfor {

		for i := 0; i < processcount; i++ {
			if procs[i].arrival == time {
				fmt.Fprintf(fileOut, "Time %3d : %s arrived\n", time, procs[i].name)
				arrivalQ = append(arrivalQ, procs[i]) //add 2 queue
				amount++                              //increase amount in queue
			}
		}

		if amount == 0 {
			fmt.Fprintf(fileOut, "Time %3d : Idle\n", time)
		}

		pos := 0  //current position
		prev := 0 //previous ID

		if amount > 0 {
			for arrivalQ[pos].finished {
				pos++
			}
			//if current pos < total and it is selected
			if pos < processcount && arrivalQ[pos].selectB {
				prev = arrivalQ[pos].PID //set previous ID to current val
			}

			for i := 0; i < len(arrivalQ); i++ {
				arrivalQ[i].selectB = false //set selected to fasle
			}

			arrivalQ = selSort("Burst", arrivalQ) //sort burst times

			for pos < processcount && arrivalQ[pos].finished {
				pos++ //incremet to next element
			}
			//set to true if seen before
			if pos < processcount && arrivalQ[pos].PID == prev {
				arrivalQ[pos].selectB = true
			}
			//completed
			if pos < processcount && arrivalQ[pos].selectB && arrivalQ[pos].burst == 0 && !arrivalQ[pos].finished {
				arrivalQ[pos].selectB = false
				arrivalQ[pos].finished = true
				arrivalQ[pos].complete = time
				amount--
				fmt.Fprintf(fileOut, "Time %3d : %s finished\n", time, arrivalQ[pos].name)
			}
			//point math
			if pos < processcount && arrivalQ[pos].finished && amount > 0 {
				pos++
			}
			//select a process if none chosen
			if pos < processcount && !arrivalQ[pos].selectB && !arrivalQ[pos].finished && amount > 0 {
				arrivalQ[pos].selectB = true
				arrivalQ[pos].selected = time
				fmt.Fprintf(fileOut, "Time %3d : %s selected (burst %3d)\n", time, arrivalQ[pos].name, arrivalQ[pos].burst)
			}
			if amount == 0 {
				fmt.Fprintf(fileOut, "Time %3d : Idle\n", time)
			}
		}
		time++

		if pos < processcount && arrivalQ[pos].burst > 0 {
			arrivalQ[pos].burst-- //dec busrt
		}
	}

	//turn around
	for i := 0; i < processcount; i++ {
		arrivalQ[i].turnaround = arrivalQ[i].complete - arrivalQ[i].arrival
	}

	arrivalQ = selSort("PID", arrivalQ)
	procs = selSort("PID", procs)

	for i := 0; i < processcount; i++ {
		arrivalQ[i].wait = arrivalQ[i].turnaround - procs[i].burst //take turnaround minus the bursts
	}
	fmt.Fprintf(fileOut, "Finished at time  %d\n\n", runfor)

	for i := 0; i < processcount; i++ {
		fmt.Fprintf(fileOut, "%s wait %3d turnaround %3d\n", arrivalQ[i].name, arrivalQ[i].wait, arrivalQ[i].turnaround)

	}
}

func RR(procs []Process, processcount int, runfor int) {
	for i := 0; i < processcount; i++ {
	}

}

//
func selSort(choice string, procs []Process) (p []Process) {
	p = make([]Process, len(procs))
	p = procs
	//3 choices arrival burst or PID
	switch choice {
	case "Arrival":
		{
			for i := 0; i < (len(p) - 1); i++ {
				selNum := i
				for j := i + 1; j < len(p); j++ {
					if p[j].arrival < p[selNum].arrival {
						selNum = j
					}
				}
				temp := p[selNum]
				p[selNum] = p[i]
				p[i] = temp
			}
		}

	case "Burst":
		{
			for i := 0; i < (len(p) - 1); i++ {
				selNum := i
				for j := i + 1; j < len(p); j++ {
					if p[j].burst < p[selNum].burst {
						selNum = j
					}
				}
				temp := p[selNum]
				p[selNum] = p[i]
				p[i] = temp
			}
		}
	case "PID":
		{
			for i := 0; i < (len(p) - 1); i++ {
				selNum := i
				for j := i + 1; j < len(procs); j++ {
					if p[j].PID < p[selNum].PID {
						selNum = j
					}
				}
				temp := p[selNum]
				p[selNum] = p[i]
				p[i] = temp
			}
		}
	}
	return p
}

//reading in struct info didnt work in main so make a function
func initStruct(file string, numProc int) (process []Process) {
	readIn, _ := os.Open(file)
	scanner := bufio.NewScanner(readIn)
	scanner.Split(bufio.ScanWords)
	process = make([]Process, numProc)

	var p Process
	count := 0
	for scanner.Scan() {
		checkStr := scanner.Text()

		if checkStr == "end" {
			break
		} else if checkStr == "name" {
			scanner.Scan()
			p.name = scanner.Text()
			p.PID++ //increment pointer
			count++
		} else if checkStr == "arrival" {
			scanner.Scan()
			p.arrival, _ = strconv.Atoi(scanner.Text())
			count++
		} else if checkStr == "burst" {
			scanner.Scan()
			p.burst, _ = strconv.Atoi(scanner.Text())
			count++
		}
		if count == 3 {
			process = append(process, p) //join the scanned slices to real slice
			count = 0                    //RESET COUNT IMPORTANT
		}

	}
	process = process[numProc:] //cut off for existing process only
	return process

}
