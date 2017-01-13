package main

import (
  "fmt"
  "time"

  "github.com/darkincred/chip8-emu/chip8"
  "github.com/darkincred/chip8-emu/disp"

  "github.com/veandco/go-sdl2/sdl"

)

// CREATE THE CHIP8 struct



func main() {
  fmt.Printf("Testing...\n")


  myChip, _ := chip8.NewCPU()
  myDisp, _ := disp.SetupGraphics()



  myChip.Init() // Initializing

  // Load a Chip8 screen and using Blocking OP for enter?

  myChip.Load("c8games/BRIX") // Load a game




  for myDisp.RunFlag {
    myChip.Cycle()

    if myChip.DrawFlag { // Frameskips  settings here
      updateGraphics(myChip,myDisp)
    }

    pollEvent(myDisp,myChip)


    time.Sleep(0 * time.Millisecond) // later n options
  }

}


func updateGraphics(c *chip8.CPU, d *disp.Disp) {

  for yLoc := 0; yLoc < 32; yLoc++ {
    for xLoc := 0; xLoc < 64; xLoc++ {

      r := c.Gfx[yLoc][xLoc] ^ 1
      if r == 0 {

        d.Renderer.SetDrawColor(255,255,255,1)
        for x := 0; x < 10; x++ {
          for y := 0; y < 10; y++ {

            d.Renderer.DrawPoint((xLoc *10) + y,(yLoc * 10) +x)
          }
        }
      } else {
          d.Renderer.SetDrawColor(0,0,0,1)
        for x := 0; x < 10; x++ {
          for y := 0; y < 10; y++ {

            d.Renderer.DrawPoint((xLoc *10) + y,(yLoc * 10) +x)
          }
        }

      }
      }

    }

  d.Renderer.Present()
  c.DrawFlag = false;
}

func pollEvent(d *disp.Disp, c *chip8.CPU) {
  var event sdl.Event
  for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
    switch t := event.(type) {
      case *sdl.KeyDownEvent:
  fmt.Printf("Key Down: %x\n", keyMap[t.Keysym.Sym])
        c.Key[keyMap[t.Keysym.Sym]] = 1
      case *sdl.KeyUpEvent:
        fmt.Printf("Key Up: %x\n", keyMap[t.Keysym.Sym])
        c.Key[keyMap[t.Keysym.Sym]] = 0
      case *sdl.QuitEvent:
        d.RunFlag = false
      default:
        event = nil

    }
  }
}

var keyMap = map[sdl.Keycode]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'a': 0x04, 'z': 0x05, 'e': 0x06, 'r': 0x0D,
	'q': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'w': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}
