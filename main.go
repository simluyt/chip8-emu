package main

import (
  "fmt"
  "os"

  "github.com/darkincred/chip8-emu/chip8"
  
)

// CREATE THE CHIP8 struct



func main() {
  fmt.Printf("Testing...\n")

  myChip := &chip8.Chip8{}


  myChip.Test()
  //myChip.Load("PONG")

  //
  // for {
  //   myChip.Cycle()
  //
  //   if myChip.V[0xF] {
  //
  //   }
  //
  // }


}
