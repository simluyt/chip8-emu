package chip8

import (
  "fmt"
  "io/ioutil"
  "math/rand"

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

  Key[16] byte


}




// Constructor

func NewCPU() (*CPU, error) {

  return &CPU{}, nil
}



// INIT
func (c *CPU) Init() {

  fmt.Printf("Initializing...\n")

  c.pc = 0x200; // Set Program Counter to start of loaded program


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

  c.opcode = c.fetch() // Fetch an opcode

  //fmt.Printf("Current instruction: 0x%X\n", c.opcode)

  c.decode() // Decode and execute opcode
}


// HElPERS

func (c *CPU) fetch() (uint16) {
    opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc + 1])
    c.pc += 2
    return opcode
}

func (c *CPU) decode() { // Decode add error support

  var x byte = byte((c.opcode & 0x0F00)>>8)
  var y byte = byte((c.opcode & 0x00F0)>>4)

  var nn byte = byte(c.opcode & 0x00FF)
  var nnn uint16 = c.opcode & 0x0FFF

  //Perhaps 2 vars that point to the adresses of Vx and Vy and one to cf

  switch (c.opcode & 0xF000) >> 12 {

    case 0:

      switch c.opcode & 0x000F {

        case 0:  //Clear the screen
          c.clearDisp()
        case 0xE:  //Return from subroutine
          c.sp--
          c.pc = c.stack[c.sp]
      }

    case 1: // 0x1NNN -> Program counter jumps to NNN

      c.pc = c.opcode & 0x0FFF

    case 2: // 0x2NNN -> Calls subroutine at NNN

      c.stack[c.sp] = c.pc
      c.sp++
      c.pc = nnn

    case 3: // 0x3XNN -> Skips next instruction if VX == NN

      if c.V[x] == nn {

        c.pc += 2
      }

    case 4: // 0x4XNN -> Skips next instruction if VX != NN

      if c.V[x] != nn {

        c.pc += 2
      }

    case 5: // 0x5XY0 -> Skips next instruction if VX == VY (Should this produce err if not ending in 0?)

      if c.V[x] != c.V[y] {

        c.pc += 2
      }

    case 6: // 0x6XNN -> Stores NN in register VX

      c.V[x] = nn

    case 7: // 0x6XNN -> Adds NN to register VX

      c.V[x] += nn

    case 8:

        switch c.opcode & 0xF {

          case 0: // 0x8XY0 --> Assign value VY to VX

            c.V[x] = c.V[y]

          case 1: // 0x8XY1 --> Assign value VY to VX OR VY

            c.V[x] = (c.V[x] | c.V[y])

          case 2: // 0x8XY2 --> Assign value VY to VX AND VY

            c.V[x] = (c.V[x] & c.V[y])

          case 3: // 0x8XY3 --> Assign value VY to VX XOR VY

            c.V[x] = (c.V[x] ^ c.V[y])


          case 4: // 0x8XY4 --> Adds value VY to VX. colFlag == carry

            result := uint16(c.V[x]) + uint16(c.V[y])

            if result > 0xFF {

                c.V[0xF] = 1
            } else {

                c.V[0xF] = 0;
            }

            c.V[x] = byte(result & 0xFF)

          case 5: // 0x8XY5 --> Subtracts value VY from VX. colFlag != borrow

            var cf byte

            if c.V[x] > c.V[y] {

              cf = 0x1
            }

            c.V[0xF] = cf
            c.V[x] -= c.V[y]

          case 6: // 0x8XY6 --> Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.

            cf := (c.opcode & 0x1)

            if cf >= 1 {
              c.V[0xF] = 1
            }

            c.V[x] = c.V[x]>>1

          case 7: // 0x8XY7 --> Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.

            var cf byte

            if c.V[y] > c.V[x] {

              cf = 0x1
            }

            c.V[0xF] = cf
            c.V[x] = c.V[y] - c.V[x]

          case 0xE: // 0x8XYE --> Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.

            var cf byte

            if (c.V[x] & 0x80) == 0x80 {

              cf = 1
            }

            c.V[0xF] = cf
            //c.V[x] = c.V[x]<<1
            c.V[x] = c.V[x] * 2

          }

    case 9: // 0x9XY0 -> Skips next instruction if VX == VY

      if c.V[x] == c.V[y] {

        c.pc += 2
      }

    case 0xA: // 0xANNN --> Sets I to NNN

      c.I = nnn

    case 0xB: // 0xBNNN --> Program counter jumps to NNN + V0

      c.pc = nnn + uint16(c.V[0])

    case 0xC: // 0xCXNN --> Assign a random(0..255) AND NN to VX

      var rnd byte = byte(rand.Intn(255))

      c.V[x] = rnd & nn

    case 0xD:

      var h uint16 = c.opcode & 0x000F

      c.drawSprite(c.V[x],c.V[y],h)

    case 0xE:

      switch nn {

        case 0x9E:

          if c.keyDown(c.V[x]) {

            c.pc += 2
          }

        case 0xA1:

          if !c.keyDown(c.V[x]) {

            c.pc += 2
          }

      }


    case 0xF:

      switch nn {

          case 0x07: // 0xFX07 --> Sets VX to the value of the delay timer.

              c.V[x] = c.delay_timer

          case 0x0A: // 0xFX0A --> A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event

              var press bool

              for i := 0; i < 16; i++ {
                if c.Key[i] == 1 {

                  c.V[x] = byte(i)
                  press = true
                }
              }

              if !press {
                c.pc -= 2
              }

          case 0x15: // 0xFX15 --> Sets the delay timer to VX.

              c.delay_timer = c.V[x]

          case 0x18: // 0xFX18 --> Sets the sound timer to VX.

              c.sound_timer = c.V[x]

          case 0x1E: // 0xFX1E --> Adds VX to I. (TODO: Don't forget the overflow rule for Spacefight)

              c.I += uint16(c.V[x])

          case 0x29: // 0xFX29 --> Sets I to the location of the sprite for the character(4x5) in VX.

              c.I = uint16(c.V[(c.opcode & 0x0F00) >> 8]) * 5

          case 0x33: // 0xFX33 --> BCD

              c.memory[c.I]     = c.V[x] / 100;
              c.memory[c.I + 1] = (c.V[x] / 10) % 10;
              c.memory[c.I + 2] = (c.V[x] % 100) % 10;

          case 0x55: //0FX55 --> Stores V0 to VX (including VX) in memory starting at address I

              for i := 0; i <= int(x); i++ {

                c.memory[uint16(c.I) + uint16(i)] = c.V[i]
              }

          case 0x65: //0FX65 --> Fills V0 to VX (including VX) with values from memory starting at address I
              for i := 0; i <= int(x); i++ {

                c.V[i] = c.memory[uint16(c.I) + uint16(i)]
              }

          }

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
  }

  c.DrawFlag = true
}

func (c *CPU) drawSprite(vX byte, vY byte, height uint16) {
  //TODO COLLISIONS + REWRITE
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


func (c *CPU) keyDown(index byte) (bool) {

  if c.Key[index] == 1  {

    return true
  } else {

    return false
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
  0xF0, 0x80, 0xF0, 0x80, 0x80 } // F
  // 0x36, 0x36, 0x60, 0x36, 0x36, // H
  // 0x28, 0x28, 0x28, 0x28, 0x28, // I
  // 0x60, 0x36, 0x60, 0x20,0x20 }// P
