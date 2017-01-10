package main

import (
  "fmt"
  "time"

  "github.com/darkincred/chip8-emu/chip8"
  "github.com/darkincred/chip8-emu/disp"


)

// CREATE THE CHIP8 struct



func main() {
  fmt.Printf("Testing...\n")


  myChip, _ := chip8.NewCPU()
 myDisp, _ := disp.SetupGraphics()




  myChip.Init() // Initializing
  myChip.Load("PONG") // Load a game
  //myChip.Test()
  // for i := 0; i < 10; i++ {
  //   myChip.Cycle()
  // }
  //
  //
  for myDisp.IsRunning() {
    myChip.Cycle()

    if myChip.DrawFlag {
      updateGraphics(myChip,myDisp)
    }

    //  myChip.Key, myChip.KeyState = myDisp.GetKey()

    if myChip.KeyPad.QuitFlag {
      myDisp.Running = false
    }

    time.Sleep(0 * time.Millisecond) // later n options
  }
  //myChip.Test()
  // for yLoc := 0; yLoc < 32 ; yLoc++ {
  //   for xLoc := 0; xLoc < 64 ; xLoc++ {
  //     fmt.Printf("%d",myChip.Gfx[yLoc][xLoc])
  //   }
  //   fmt.Printf("\n")
  // }

}

func updateGraphics(c *chip8.CPU, d *disp.Disp) {

  for yLoc := 0; yLoc < 32; yLoc++ {
    for xLoc := 0; xLoc < 64; xLoc++ {
      if c.Gfx[yLoc][xLoc] == 1 {
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
