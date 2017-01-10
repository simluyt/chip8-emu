package input

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {

  keys [16] byte
  event sdl.Event

}

func (*key Keyboard) Poll() {

  	for event = sdl.PollEvent(); event == *sdl.KeyDownEvent; event = sdl.PollEvent() {
  		switch t := event.(type) {
      case *sdl.KeyDownEvent:
        switch t.Keysym.Sym {
        case sdl.GetKeyFromName("&"):
          return 0x01, t.State;
        case sdl.GetKeyFromName("é"):
          return 0x02, t.State;
        case sdl.GetKeyFromName("\""):
          return 0x03, t.State;
        case sdl.GetKeyFromName("'"):
          return 0x0C, t.State;
        case sdl.GetKeyFromName("a"):
          return 0x04, t.State;
        case sdl.GetKeyFromName("z"):
          return 0x05, t.State;
        case sdl.GetKeyFromName("e"):
          return 0x06, t.State;
        case sdl.GetKeyFromName("r"):
          return 0x0D, t.State;
        case sdl.GetKeyFromName("q"):
          return 0x07, t.State;
        case sdl.GetKeyFromName("s"):
          return 0x08, t.State;
        case sdl.GetKeyFromName("d"):
          return 0x09, t.State;
        case sdl.GetKeyFromName("f"):
          return 0x0E, t.State;
        case sdl.GetKeyFromName("w"):
          return 0x0A, t.State;
        case sdl.GetKeyFromName("x"):
          return 0x00, t.State;
        case sdl.GetKeyFromName("c"):
          return 0x0B, t.State;
        case sdl.GetKeyFromName("v"):
          return 0x0F, t.State;
        default:
          continue
        }
      }

  		}
      return 0xFF, 0xFF
  	}

}

func (*key Keyboard) Stop() {

}
