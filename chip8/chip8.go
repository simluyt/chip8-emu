package chip8

import (
  "fmt"
  "io/ioutil"
)


// Datatype for the chip8 structure

type Chip8 struct {

  // OPCODE
  opcode uint16;


  // MEMORY

  memory [4096] byte;


  // REGISTER

  V [16] byte;

// INDEX AND PROGRAM COUNTER

  I uint16;
  pc uint16;

// GRAPHICS OUTPUT

  drawflag bool;

  gfx [64][32] byte;

// TIMERS

  delay_timer byte;
  sound_timer byte;

// STACK

  stack [16] uint16;
  sp uint16;

// KEYPAD

  key [16] byte;

// PKG TEST


}



// INIT METHOD
func (c *Chip8) Init() {

  fmt.Printf("Initializing...\n")

  // Golang automaticaly should set all vars to 0 when declared.

  c.pc = 0x200;

  for i := 0; i < 80; i++ {
    c.memory[i] = fontset[i]
  }



}

func (c *Chip8) Load(filename string)  {

  program, _:= ioutil.ReadFile(filename)

  for i := 0; i < len(program); i++ {
    c.memory[c.pc + uint16(i)] = program[i]
  }
}

// EMULATION CYCLE

func (c *Chip8) Cycle() {

  fmt.Printf("Emulating one Cycle...\n")

  // FETCH
  c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc + 1])


  c.pc += 2 // Points the Program Counter to the next first part of an opcode

  // DECODE maybe in aparte func met kleine letterrr

  switch c.opcode & 0xF000 {
    case 0x0000:
      switch c.opcode & 0x000F {
        // Ommiting 0x0NNN for now...
        case 0x0000:  //Clear the screen
        case 0x000E:  //Return from subroutine
          c.sp--
          c.pc = c.stack[c.sp]
        default:
          fmt.Printf("Invalid instruction: 0x%X", c.opcode)
        }
    case 0x1000: // 0x1NNN -> Program counter jumps to NNN
      c.pc = c.opcode & 0x0FFF

    case 0x2000: // 0x2NNN -> Calls subroutine at NNN
      c.stack[c.sp] = c.pc
      c.sp++
      c.pc = c.opcode & 0x0FFF

case 0x3000: // 0x3XNN -> Skips next instruction if VX == NN
  x := c.opcode & 0x0F00>>2
  nn := c.opcode & 0x00FF


  if c.V[x] == byte(nn) {
    c.pc += 2 // Opcodes are 2 memory spaces (bytes)
  }

case 0x4000: // 0x4XNN -> Skips next instruction if VX != NN
  x := c.opcode & 0x0F00>>2
  nn := c.opcode & 0x00FF

  if c.V[x] != byte(nn) {
    c.pc += 2// Opcodes are 2 memory spaces (bytes)
  }

case 0x5000: // 0x5XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)
  x := c.opcode & 0x0F00>>2
  y := c.opcode & 0x00F0>>1

  if c.V[x] != c.V[y] {
    c.pc += 2// Opcodes are 2 memory spaces (bytes)
  }

case 0x6000:  // 0x6XNN -> Stores NN in register VX
  x := c.opcode & 0x0F00>>2
  nn := c.opcode & 0x00FF

  c.V[x] = byte(nn)

case 0x7000:  // 0x6XNN -> Adds NN to register VX
  x := c.opcode & 0x0F00>>2
  nn := c.opcode & 0x00FF

  c.V[x] += byte(nn)

case 0x8000:
  switch c.opcode & 0xF00F {
  case 0x8000:
  case 0x8001:
  case 0x8002:
  case 0x8003:
  case 0x8004:
  case 0x8005:
  case 0x8006:
  case 0x8007:
  case 0x800E:
  default:
    fmt.Printf("Invalid instruction: 0x%X", c.opcode)
  }
case 0x9000:
case 0xA000:
case 0xB000:
case 0xC000:
case 0xD000:
case 0xE000:
case 0xF000:
default:
  fmt.Printf("Invalid instruction: 0x%X", c.opcode)

}
}


// TEST FUNC

func (c *Chip8) Test() {
  //fmt.Printf("Testing...\n")




  // for i := 0x200; i < 0x230; i++ {
  //   fmt.Printf("%x\n", c.memory[i])
  // }



}


// CHIPS8_FONTSET

var fontset = [80]byte{
  0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
  0x20, 0x60, 0x20, 0x20, 0x70, // 1
  0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
  0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
  0x90, 0x90, 0xF0, 0x10, 0x10, // 4
  0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
  0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
  0xF0, 0x10, 0x20, 0x40, 0x40, // 7
  0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
  0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
  0xF0, 0x90, 0xF0, 0x90, 0x90, // A
  0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
  0xF0, 0x80, 0x80, 0x80, 0xF0, // C
  0xE0, 0x90, 0x90, 0x90, 0xE0, // D
  0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
  0xF0, 0x80, 0xF0, 0x80, 0x80 }  // F
