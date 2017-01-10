package input

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {

  Keys [16] byte

  //PollFlag bool
  QuitFlag bool
  Latest byte

}

func (key *Keyboard) Poll() {
    var event sdl.Event
  	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
  		switch t := event.(type) {
        case *sdl.KeyDownEvent:
          key.Latest = byte(t.Keysym.Unicode)
          key.Keys[keyMapAzerty[key.Latest]] = 1
        case *sdl.KeyUpEvent:
          key.Latest = byte(t.Keysym.Unicode)
          key.Keys[keyMapAzerty[key.Latest]] = 0
        case *sdl.QuitEvent:
          key.QuitFlag = true
        default:
          event = nil
          continue
      }
    }
}



var keyMapAzerty = map[byte]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'a': 0x04, 'z': 0x05, 'e': 0x06, 'r': 0x0D,
	'q': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'w': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}

var keyMap = map[rune]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
	'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}
