package main

import (
    "Driver-go/elevio"
    "fmt"
)

func main() {

    numFloors := 9

    elevio.Init("localhost:12345", numFloors)

    var d elevio.MotorDirection = elevio.MD_Stop
    elevio.SetMotorDirection(d)

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors := make(chan int)
    drv_obstr := make(chan bool)
    drv_stop := make(chan bool)
    var target int = 0
    var currentFloor int = 0
    var moving bool = false

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

    for {
        select {
        case a := <-drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            if !moving {
                target = a.Floor
                moving = true
                fmt.Printf("Target is set to %d\n", target)
                fmt.Printf("Floor is %d\n", currentFloor)
                if a.Floor == currentFloor {
                    elevio.SetMotorDirection(elevio.MD_Stop)
                    fmt.Printf("Stopped")
                    moving = false
                } else if a.Floor > currentFloor {
                    elevio.SetMotorDirection(elevio.MD_Up)
                    fmt.Printf("Going up")
                } else {
                    elevio.SetMotorDirection(elevio.MD_Down)
                    fmt.Printf("Going down")
                }
            }

        case a := <-drv_floors:
            fmt.Printf("%+v\n", a)
            currentFloor = a

            if a == target {
                elevio.SetMotorDirection(elevio.MD_Stop)
                fmt.Printf("Reached target")
                moving = false
            } else if a < target {
                elevio.SetMotorDirection(elevio.MD_Up)
                fmt.Printf("still going up")
            } else {
                elevio.SetMotorDirection(elevio.MD_Down)
                fmt.Printf("still going down")
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
        }
    }
}
