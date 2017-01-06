package chip8

import (
  "fmt"
  "io/ioutil"
  "math/rand"
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
  x := (c.opcode & 0x0F00)>>2
  nn := c.opcode & 0x00FF


  if c.V[x] == byte(nn) {
    c.pc += 2 // Opcodes are 2 memory spaces (bytes)
  }

case 0x4000: // 0x4XNN -> Skips next instruction if VX != NN
  x := (c.opcode & 0x0F00)>>2
  nn := c.opcode & 0x00FF

  if c.V[x] != byte(nn) {
    c.pc += 2// Opcodes are 2 memory spaces (bytes)
  }

case 0x5000: // 0x5XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)
  x := (c.opcode & 0x0F00)>>2
  y := (c.opcode & 0x00F0)>>1

  if c.V[x] != c.V[y] {
    c.pc += 2// Opcodes are 2 memory spaces (bytes)
  }

case 0x6000:  // 0x6XNN -> Stores NN in register VX
  x := (c.opcode & 0x0F00)>>2
  nn := c.opcode & 0x00FF

  c.V[x] = byte(nn)

case 0x7000:  // 0x6XNN -> Adds NN to register VX
  x := (c.opcode & 0x0F00)>>2
  nn := c.opcode & 0x00FF

  c.V[x] += byte(nn)

case 0x8000:
  x := (c.opcode & 0x0F00)>>2
  y := (c.opcode & 0x00F)>>1
  switch c.opcode & 0xF00F {

    case 0x8000: // 0x8XY0 --> Assign value VY to VX
      c.V[x] = c.V[y]
    case 0x8001: // 0x8XY1 --> Assign value VY to VX OR VY
      c.V[x] = c.V[x] | c.V[y]
    case 0x8002: // 0x8XY2 --> Assign value VY to VX AND VY
      c.V[x] = c.V[x] & c.V[y]
    case 0x8003: // 0x8XY3 --> Assign value VY to VX XOR VY
      c.V[x] = c.V[x] ^ c.V[y]
    case 0x8004: // 0x8XY4 --> Adds value VY to VX if there is a carry set VF to 1 else to 0
    case 0x8005: // 0x8XY5 --> Subtracts value VY to VX if there is a carry set VF to 1 else to 0
    case 0x8006:
    case 0x8007:
    case 0x800E:
  default:
    fmt.Printf("Invalid instruction: 0x%X", c.opcode)
  }
case 0x9000: // 0x9XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)
  x := (c.opcode & 0x0F00)>>2
  y := (c.opcode & 0x00F0)>>1

  if c.V[x] != c.V[y] {
    c.pc += 2// Opcodes are 2 memory spaces (bytes)
  }
case 0xA000: // 0xANNN --> Sets I to NNN
  c.I = c.opcode & 0x0FFF

case 0xB000: // 0xBNNN --> Program counter jumps to NNN + V0
  c.pc = (c.opcode & 0x0FFF) + uint16(c.V[0])

case 0xC000: // 0xCXNN --> Assign a random(0..255) AND NN to VX
  x := (c.opcode & 0x0F00)>>2
  c.V[x] = byte(rand.Intn(255)) & 0x00FF

case 0xD000: // Draw
case 0xE000:
  switch c.opcode & 0xE0F0 {
    case 0xE090:
    case 0xE0A0:
    default:
    fmt.Printf("Invalid instruction: 0x%X", c.opcode)

  }
case 0xF000:
  x := (c.opcode & 0x0F00)>>2

  switch c.opcode & 0x00F0 {
    case 0x0000:
      switch c.opcode & 0x000F {
        case 0x0007: // 0xFX07 --> Sets VX to the value of the delay timer.
          c.V[x] = c.delay_timer

        case 0x000A: // 0xFX0A --> A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)

        default:
          fmt.Printf("Invalid instruction: 0x%X", c.opcode)
        }

    case 0x0010:
      switch c.opcode & 0x000F {
        case 0x0005: // 0xFX15 --> Sets the delay timer to VX.
          c.delay_timer = c.V[x]
        case 0x0008: // 0xFX18 --> Sets the sound timer to VX.
          c.sound_timer = c.V[x]
        case 0x000E: // 0xFX1E --> Adds VX to I. (TODO: Don't forget the overflow rule for Spacefight)
          c.I = c.V[x]
        default:
          fmt.Printf("Invalid instruction: 0x%X", c.opcode)
      }

  case 0x0020: // 0FX29 --> Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.

  case 0x0030:
  case 0x0050:
  case 0x0060:

  }
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
