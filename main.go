package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

type Queue struct {
	orders []int
}

func (q *Queue) Enqueue(floor int) {
	q.orders = append(q.orders, floor)
}

func (q *Queue) Dequeue() int {
	if len(q.orders) == 0 {
		return -1
	}
	floor := q.orders[0]
	q.orders = q.orders[1:]
	return floor
}

func (q *Queue) Peek() int {
	if len(q.orders) == 0 {
		return -1
	}
	return q.orders[0]
}

func (q *Queue) sneak(floor int) {
	q.orders = append([]int{floor}, q.orders...)
}
	"Driver-go/elevio"
	"fmt"
	"time"
)

type Queue struct {
	orders []int
}

func (q *Queue) Enqueue(floor int) {
	q.orders = append(q.orders, floor)
}

func (q *Queue) Dequeue() int {
	if len(q.orders) == 0 {
		return -1
	}
	floor := q.orders[0]
	q.orders = q.orders[1:]
	return floor
}

func (q *Queue) Peek() int {
	if len(q.orders) == 0 {
		return -1
	}
	return q.orders[0]
}

func (q *Queue) sneak(floor int) {
	q.orders = append([]int{floor}, q.orders...)
}

func main() {
func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Stop
	var currentFloor int = -1

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			//queue.Enqueue(a.Floor)
			if a.Button == elevio.BT_HallUp || a.Button == elevio.BT_HallDown {
				queue.Enqueue(a.Floor)
			} else if a.Button == elevio.BT_Cab {
				queue.sneak(a.Floor)
			}
			targetFloor := queue.Peek()
			if targetFloor > currentFloor {
				d = elevio.MD_Up
			} else if targetFloor < currentFloor {
				d = elevio.MD_Down
			} else {
				d = elevio.MD_Stop
				elevio.SetButtonLamp(elevio.BT_HallUp, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_HallDown, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_Cab, currentFloor, false)
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			currentFloor = a
			if currentFloor == queue.Peek() {
				d = elevio.MD_Stop
				elevio.SetMotorDirection(d)

				elevio.SetButtonLamp(elevio.BT_HallUp, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_HallDown, currentFloor, false)
				elevio.SetButtonLamp(elevio.BT_Cab, currentFloor, false)
				// wait for 3 seconds
				elevio.SetDoorOpenLamp(true)
				elevio.SetFloorIndicator(currentFloor)
				time.Sleep(3 * time.Second)
				elevio.SetDoorOpenLamp(false)

				queue.Dequeue()
				if len(queue.orders) > 0 {
					targetFloor := queue.Peek()
					if targetFloor > currentFloor {
						d = elevio.MD_Up
					} else if targetFloor < currentFloor {
						d = elevio.MD_Down
					} else {
						d = elevio.MD_Stop
					}
					elevio.SetMotorDirection(d)
				}
			}

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			queue = Queue{} // Clear the queue
		}
	}
}
