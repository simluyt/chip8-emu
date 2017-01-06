package main

import (
  "fmt"

  "github.com/darkincred/chip8-emu/chip8"
)

// CREATE THE CHIP8 struct



func main() {
  fmt.Printf("Testing...\n")

  myChip := &chip8.Chip8{}

  myChip.Init()


}
