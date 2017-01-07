package disp

import (

  "fmt"
  "os"
  "github.com/veandco/go-sdl2/sdl"
)

var (
  WindowTitle = "myChip8"
  WindowWidth = 320
  WindowHeight = 640
)

type Disp struct {
  Window *sdl.Window
  Renderer *sdl.Renderer
  Running bool
}

func SetupGraphics() (*Disp, error){

  window, err := sdl.CreateWindow(WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowHeight, WindowWidth, sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
    return nil, err
	}


	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return nil, err
	}


    renderer.Clear()
    renderer.Present()



return &Disp{
		Window:   window,
		Renderer: renderer,
	}, nil

}
func (disp *Disp) HandleEvents(event sdl.Event) {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			disp.Running = false
		}
	}
}
