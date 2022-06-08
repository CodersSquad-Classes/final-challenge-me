# Little Pacman

> My little version of the classic game, Pac-Man

## Video

https://drive.google.com/drive/folders/1bJTk1DHh3vgNULsOAod8d7-YcyVKLZe-?usp=sharing

## Introduction

The presented project is a tribute to the classic pac-man game. This version is a pretty simple re-interpretation of the game, which is run completly in the console. The map of the game can be loaded, and there can be N total number of ghosts. The enemy ghost have a pretty simple AI, which is move randomly to any direction available. This document pretends to show the software architecture of the software project.

## Requirements

[x] The game's maze layout can be static.
[x] The pacman gamer must be controlled by the user.
[x] Enemies are autonomous entities that will move a random way.
[x] Enemies and pacman should respect the layout limits and walls.
[x] Enemies number can be configured on game's start.
[x] Each enemy's behaviour will be implemented as a separated thread.
[x] Enemies and pacman threads must use the same map or game layout data structure resource.
[x] Display obtained pacman's scores.
[x] Pacman loses when an enemy touches it.
[x] Pacman wins the game when it has taken all coins in the map.

## Proposed Architecture

For the creation of the project I decided to develop the game directly in console, with a non-traditional GUI. To do this we can defined our inputs:


| Input                      | Output                                 |
|----------------------------|----------------------------------------|
|     Keyboard up-arrow      |     pac-man moves up if possible       |
|     Keyboard down-arrow    |     pac-man moves down if possible     |
|     Keyboard left-arrow    |     pac-man moves left if possible     |
|     Keyboard right-arrow   |     pac-man moves right if possible    |
|     Keyboard escape key    |     Finishes the game                  |


The software works with go-routines allowing me to control multiple parts of the system almost simultaneously. We can define the routines as:

* GUI Manager: controls the GUI drawing 
* Input manager: controls the input from the user
* Enemy ghost: One thread per enemy, it randomly moves the ghost through the maze


## Map Parsing

The user can create its map manually, and the program will read for a file named `./maps/map.txt`

This file is constructed by:
- `#` Wall
- `g` Ghost
- `C` Pacman
- `*` Coin

E.g.

```
WWWWWWWW
W***C**W
Wg*****W
WWWWWWWW
```
This last map, generates a 8x4 size map, with Pacman and one ghost, and it's filled with coins.

## Structures description

_Item_
Attributes: int x, int y.

This structure defines the basic element for the game, with coordinates

_Map_
Attributes: 
    height     int
	width      int
	pacman     Item
	ghosts     []Item
	walls      [][]bool
	coins      [][]bool
	totalCoins int
	direction  int

This structure defines all data of the map, location of all items. The walls and coins is a boolean 2d array which is a representation of the map, if the coin/wall is in a specific coordinate then is true. If the coin is taken, then it turns it to false.

## End of game

- If the user closes the game (ctrl+c or escape)
- If a ghost touches the user
- If the user takes all coins

## Used libraries
[tcell:](github.com/gdamore/tcell/v2) Provides an API that allows the programmer to write text-based user interfaces in the terminal. The API provides functions to move the cursor, create windows, produce colors, play with mouse, etc. 
