package main

import (
  "fmt"

  "github.com/veandco/go-sdl2/sdl"
  "github.com/darkincred/chip8-emu/chip8"
  "github.com/darkincred/chip8-emu/disp"
)

// CREATE THE CHIP8 struct



func main() {
  fmt.Printf("Testing...\n")

  myChip := &chip8.Chip8{}


  //myChip.Test()
  MyDisp, _ := disp.SetupGraphics()

  var event sdl.Event
  MyDisp.Running = true

  if MyDisp.Running {
    fmt.Printf("Running..\n")
    }


  myChip.Load("PONG")


  for MyDisp.Running {
    //myChip.Cycle()
    MyDisp.HandleEvents(event)

  }


}
