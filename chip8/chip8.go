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
  c.pc += 2

  // DECODE maybe in aparte func met kleine letterrr

  switch c.opcode & 0xF000 {
  case 0x000:


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
