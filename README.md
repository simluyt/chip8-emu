# chip8-emu


# Short Introduction

This is my first emulator project. I chose a CHIP8 emulator because according to sources online it was the easiest project to get started

I'm using Go because I wanted to get a little bit more familiar with this language than the classic "hello world" application.


# Chip8

I use go-sdl2 for rendering the graphics. These are sdl2 bindings made by veandco


# Status

This is a WIP 

Almost all the Opcodes are implemented correctly, but they need refactoring.

## TODO

+ Scores draw imperfectly, guess this has to do with imperfect opcode implemntation
+ Collisions with the borders of the window don't seem to work in PONG or BRIX
+ idk wtf MAZE even supposed to do

## Works more or less

+ Input --> Have yet to test if I press unmapped key
+ Display --> Draws weird scores tho

