package chip8

import (
  "fmt"
  "io/ioutil"
  "math/rand"

  "github.com/darkincred/chip8-emu/input"

)




// Datatype for the chip8 structure

type CPU struct {

  opcode uint16; // Current opcode
  memory [4096] byte; // Memory of 4096 bytes
  V [16] byte;  // 16 Registers
  I uint16; // Index
  pc uint16; // Program counter

  Gfx [32][64] byte;
  DrawFlag bool;

  delay_timer byte;
  sound_timer byte;

  stack [16] uint16; // Stack
  sp uint16; // Stack pointer

  KeyPad *input.Keyboard

}




// Constructor

func NewCPU() (*CPU, error) {

  return &CPU{}, nil
}



// INIT
func (c *CPU) Init() {

  fmt.Printf("Initializing...\n")

  c.pc = 0x200; // Set Program Counter to start of loaded program

  c.KeyPad = &input.Keyboard{}

  for i := 0; i < 16; i++ {
    c.V[i] = 0x00// Load fontset in memory
  }

  for i := 0; i < 80; i++ {
    c.memory[i] = fontset[i]// Load fontset in memory
  }
}

// LOAD

func (c *CPU) Load(filename string)  {

  fmt.Printf("Loading: %s...\n", filename)
  program, _:= ioutil.ReadFile(filename)

  for i := 0; i < len(program); i++ {
    c.memory[c.pc + uint16(i)] = program[i]
  }
}

// EMULATION CYCLE

func (c *CPU) Cycle() {

  //fmt.Printf("Emulating one Cycle...\n")

  c.opcode = c.fetch() // Fetch an opcode

  fmt.Printf("Current instruction: 0x%X\n", c.opcode)

  c.decode() // Decode and execute opcode
}


// HElPERS

func (c *CPU) fetch() (uint16) {
    opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc + 1])
    c.pc += 2
    return opcode
}

func (c *CPU) decode() { // Decode add error support

  switch c.opcode & 0xF000 {

    case 0x0000:
      switch c.opcode & 0x000F {
        // Ommiting 0x0NNN for now...
        case 0x0000:  //Clear the screen
          c.clearDisp()
        case 0x000E:  //Return from subroutine
          c.sp--
          c.pc = c.stack[c.sp]
        default:
          fmt.Printf("Invalid instruction: 0x%X", c.opcode)
        }

    case 0x1000: // 0x1NNN -> Program counter jumps to NNN
      c.pc = c.opcode & 0x0FFF

      fmt.Printf("JUMP to 0x%X\n", c.pc)

    case 0x2000: // 0x2NNN -> Calls subroutine at NNN
      c.stack[c.sp] = c.pc
      c.sp++
      c.pc = c.opcode & 0x0FFF

    case 0x3000: // 0x3XNN -> Skips next instruction if VX == NN
      x := (c.opcode & 0x0F00)>>8
      nn := c.opcode & 0x00FF

      fmt.Printf("SKIP NEXT INSTRUCTION\n")

      if c.V[x] == byte(nn) {

        c.pc += 2 // Opcodes are 2 memory spaces (bytes)
      }

    case 0x4000: // 0x4XNN -> Skips next instruction if VX != NN
      x := (c.opcode & 0x0F00)>>8
      nn := c.opcode & 0x00FF

      if c.V[x] != byte(nn) {
        fmt.Printf("SKIP NEXT INSTRUCTION\n")
        c.pc += 2// Opcodes are 2 memory spaces (bytes)
      }

    case 0x5000: // 0x5XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)
      x := (c.opcode & 0x0F00)>>8
      y := (c.opcode & 0x00F0)>>4

      if c.V[x] != c.V[y] {
        fmt.Printf("SKIP NEXT INSTRUCTION\n")
        c.pc += 2// Opcodes are 2 memory spaces (bytes)
      }

    case 0x6000: // 0x6XNN -> Stores NN in register VX
      x := (c.opcode & 0x0F00)>>8
      nn := c.opcode & 0x00FF
      fmt.Printf("STORE 0x%X in REG 0x%X\n", nn,x)

      c.V[x] = byte(nn)

    case 0x7000: // 0x6XNN -> Adds NN to register VX
      x := (c.opcode & 0x0F00)>>8
      nn := c.opcode & 0x00FF
      fmt.Printf("ADDS 0x%X to REG 0x%X\n", nn,x)
      c.V[x] += byte(nn)

    case 0x8000:
      x := (c.opcode & 0x0F00)>>8
      y := (c.opcode & 0x00F0)>>4

        switch c.opcode & 0xF00F {

          case 0x8000: // 0x8XY0 --> Assign value VY to VX
            c.V[x] = c.V[y]

          case 0x8001: // 0x8XY1 --> Assign value VY to VX OR VY
            c.V[x] = (c.V[x] | c.V[y])

          case 0x8002: // 0x8XY2 --> Assign value VY to VX AND VY
            c.V[x] = (c.V[x] & c.V[y])

          case 0x8003: // 0x8XY3 --> Assign value VY to VX XOR VY
            c.V[x] = (c.V[x] ^ c.V[y])


          case 0x8004: // 0x8XY4 --> Adds value VY to VX if there is a carry set VF to 1 else to 0
            if (c.V[y] > (0xFF - c.V[x])){

                c.V[0xF] = 1 // Carry
            } else {

                c.V[0xF] = 0;
            }
            c.V[x] += c.V[y]

          case 0x8005: // 0x8XY5 --> Subtracts value VY from VX if there is a borrow set VF to 0 else to 1
            if ((c.V[x] & 0x0F)) < ((c.V[y] & 0x0F)){

              c.V[0xF] = 0 //borrow

            } else {

                c.V[0xF] = 1;
            }
            c.V[x] -= c.V[y]

          case 0x8006: // 0x8XY6 --> Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.
            c.V[0xF] = byte(c.opcode & 0x000F)
            c.V[x] = c.V[x]>>4
            fmt.Printf("BITSHIFT RIGHT REG 0x%X : 0x%X\n",x, c.V[x])

          case 0x8007: // 0x8XY7 --> Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
            if ((c.V[y] & 0x000F)) < ((c.V[x] & 0x000F)) {

              c.V[0xF] = 0 //borrow

              } else {

                c.V[0xF] = 1;
              }
            c.V[x] = c.V[y] - c.V[x]
          case 0x800E: // 0x8XYE --> Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.
            c.V[0xF] = byte(c.V[x]>>8)
            c.V[x] = c.V[x]<<4
            fmt.Printf("BITSHIFT RIGHT REG 0x%X : 0x%X\n",x, c.V[x])
          default:
            fmt.Printf("Invalid instruction: 0x%X\n", c.opcode)
          }// Multiple Opcodes

    case 0x9000: // 0x9XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)
      x := (c.opcode & 0x0F00)>>8
      y := (c.opcode & 0x00F0)>>4

      if c.V[x] == c.V[y] {
        c.pc += 2// Opcodes are 2 memory spaces (bytes)
      }

    case 0xA000: // 0xANNN --> Sets I to NNN
      c.I = c.opcode & 0x0FFF
      fmt.Printf("SET ADDR I to 0x%X\n",c.opcode & 0x0FFF)

    case 0xB000: // 0xBNNN --> Program counter jumps to NNN + V0
      c.pc = (c.opcode & 0x0FFF) + uint16(c.V[0])

    case 0xC000: // 0xCXNN --> Assign a random(0..255) AND NN to VX
      x := (c.opcode & 0x0F00)>>8
      nn := (c.opcode & 0x00FF)
      rnd := rand.Intn(255)
      c.V[x] = byte(rnd) & byte(nn)

      fmt.Printf("ASSIGN rand: %X to REG 0x%X\n",rnd ,x)

    case 0xD000:
      x := byte((c.opcode & 0x0F00)>>8)
      y := byte((c.opcode & 0x00F0)>>4)

      fmt.Printf("Locations: 0x%X and 0x%X\n", c.V[y],c.V[x])
      fmt.Printf("Locations: %d and %d\n", c.V[y],c.V[x])

      h := (c.opcode & 0x000F)

      c.drawSprite(c.V[x],c.V[y],h)

    case 0xE000:
      switch c.opcode & 0xE0FF {
        case 0xE09E:
          if c.keyDown(c.V[(c.opcode & 0x0F00) >> 8]) {
            c.pc += 2
          }

        case 0xE0A1:
          if !c.keyDown(c.V[(c.opcode & 0x0F00) >> 8]) {
            c.pc += 2
          }

        default:
          fmt.Printf("Invalid instruction: 0x%X", c.opcode)

      }// Multiple Opcodes


    case 0xF000:
      x := (c.opcode & 0x0F00)>>8

      switch c.opcode & 0x00FF {

          case 0x0007: // 0xFX07 --> Sets VX to the value of the delay timer.
              c.V[x] = c.delay_timer

          case 0x000A: // 0xFX0A --> A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event

            c.getKey()
            c.V[x]  = c.KeyPad.Latest

          case 0x0015: // 0xFX15 --> Sets the delay timer to VX.
              c.delay_timer = c.V[x]

          case 0x0018: // 0xFX18 --> Sets the sound timer to VX.
              c.sound_timer = c.V[x]

          case 0x001E: // 0xFX1E --> Adds VX to I. (TODO: Don't forget the overflow rule for Spacefight)
              c.I += uint16(c.V[x])

          case 0x0029: // 0FX29 --> Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
              c.I = uint16(c.V[(c.opcode & 0x0F00) >> 8])
          case 0x0033:
              c.memory[c.I]     = c.V[(c.opcode & 0x0F00) >> 8] / 100;
              c.memory[c.I + 1] = (c.V[(c.opcode & 0x0F00) >> 8] / 10) % 10;
              c.memory[c.I + 2] = (c.V[(c.opcode & 0x0F00) >> 8] % 100) % 10;

          case 0x0055: //0FX55 --> Stores V0 to VX (including VX) in memory starting at address I
              for i := 0; i <= int(x); i++ {
                c.memory[uint16(c.I) + uint16(i)] = c.V[i]
              }

          case 0x0065: //0FX65 --> Fills V0 to VX (including VX) with values from memory starting at address I
              for i := 0; i <= int(x); i++ {
                c.V[i] = c.memory[uint16(c.I) + uint16(i)]
              }

          default:
              fmt.Printf("Invalid instruction: 0x%X", c.opcode)

          }// Multiple Opcodes

    default:
      fmt.Printf("Invalid instruction: 0x%X", c.opcode)

    }
    c.timerUpdate()
}

func (c *CPU) timerUpdate() {

  if c.delay_timer > 0 {
    c.delay_timer--
  }

  if c.sound_timer > 0 {
    if c.sound_timer == 1 {
      fmt.Printf("BEEP!\n")
    }
    c.sound_timer--
  }
}

// Gfx HElPERS

func (c *CPU) clearDisp() {
  for xLoc := 0; xLoc < 32 ; xLoc++ {
    for yLoc := 0; yLoc < 64 ; yLoc++ {
      c.Gfx[xLoc][yLoc] = 0;
    }
  } // Clear the display
  c.DrawFlag = true
}

func (c *CPU) drawSprite(vX byte, vY byte, height uint16) { // OPnieuw

  var pixel byte

  c.V[0xF] = 0;

for yLine := 0; uint16(yLine) < height; yLine++ {
    pixel = c.memory[c.I + uint16(yLine)]
    for xLine := 0; xLine < 8; xLine++ {
      if ((pixel & ( 0x80 >> byte(xLine))) != 0) {
        if c.Gfx[(vY + byte(yLine))][(vX + byte(xLine))] == 1 {
          c.V[0xF] = 1
        }
      c.Gfx[(vY + byte(yLine))][(vX + byte(xLine))] ^= 1
      }
    }
  }

  c.DrawFlag = true
}

func (c *CPU) getKey() {
  c.KeyPad.Poll()
}

func (c *CPU) keyDown(index byte) (bool) {
  if c.KeyPad.Keys[index] == 1  {
    return true
  } else {
    return false
  }
}

// TEST FUNC

func (c *CPU) Test() {

  c.Init()
  c.Load("PONG")

  for i := 0; i <= 550; i++ {
    c.Cycle()
  }

//   c.memory[0x200] = 0x6A
//   c.memory[0x201] = 0x02
//   c.memory[0x202] = 0x6B
//   c.memory[0x203] = 0x0C
//
//   c.Cycle()
//   c.Cycle()
//
//   fmt.Printf("REG 1: 0x%X and REG 2: 0x%X\n", c.V[0x1], c.V[0x2])
//
//   c.memory[0x204] = 0x6C
//   c.memory[0x205] = 0x02
//
//   c.memory[0x206] = 0x6D
//   c.memory[0x207] = 0x3F
//
//   c.Cycle()
//   c.Cycle()
//
//   fmt.Printf("REG 1: 0x%X and REG 2: 0x%X\ndrawFlag: 0x%X\n", c.V[0x1], c.V[0x2],c.V[0xF])
//
//   fmt.Printf("Screen before draw\n")
//
//   for xLoc := 0; xLoc < 32 ; xLoc++ {
//     for yLoc := 0; yLoc < 64 ; yLoc++ {
//       fmt.Printf("%d",c.Gfx[xLoc][yLoc])
//     }
//     fmt.Printf("\n")
//   }
//
// fmt.Printf("Screen after draw\n")
//
// c.memory[0x208] = 0xA2
// c.memory[0x209] = 0xEA
// c.Cycle()
// c.memory[0x20A] = 0xDA
// c.memory[0x20B] = 0xB6
//
// c.Cycle()
//
// c.memory[0x20C] = 0xDC
// c.memory[0x20D] = 0xD6
//
// c.Cycle()

for yLoc := 0; yLoc < 32 ; yLoc++ {
  for xLoc := 0; xLoc < 64 ; xLoc++ {
    fmt.Printf("%d",c.Gfx[yLoc][xLoc])
  }
  fmt.Printf("\n")
}
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
