/*

  TODO Catch errors


*/
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
  RunFlag bool
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



return &Disp{
		Window:   window,
		Renderer: renderer,
    RunFlag: true,
	}, nil

}
