package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const (
	DelayMs = 750
)

var grid Grid

func init() {
	n := len(os.Args)
	if n != 3 && n != 4 {
		PrintHelp()
		os.Exit(1)
	}
	w, errw := strconv.Atoi(os.Args[1])
	h, errh := strconv.Atoi(os.Args[2])
	if errw != nil || errh != nil {
		os.Exit(1)
	}
	switch n {
	case 3:
		rand.Seed(time.Now().UTC().UnixNano())
		grid = NewRandGrid(w, h)
	case 4:
		grid = NewGridFromString(os.Args[3], w, h)
	}
}

func PrintHelp() {
	fmt.Printf(`usage: %v <width> <height>
       %[1]v <width> <height> <data>

`, os.Args[0])
}

func main() {
	t := time.NewTicker(DelayMs * time.Millisecond)
	fmt.Println(grid)
	grid = NextGridFrom(grid)
	update := func() {
		ClearConsole()
		fmt.Print(grid)
		grid = NextGridFrom(grid)
	}
	update()
	for range t.C {
		update()
	}
}

func ClearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

const (
	GridMargin = 1
	Threshold  = 0.6
)

type Grid struct {
	Data          []int
	Width, Height int
	Stride        int
}

func NewGrid(w, h int) Grid {
	if w <= 0 || h <= 0 {
		os.Exit(1)
	}
	g := Grid{Width: w, Height: h, Stride: w + GridMargin*2}
	g.Data = make([]int, g.Stride*(g.Height+GridMargin*2))
	return g
}

func NewGridFromData(data []int, w, h int) Grid {
	if len(data) != w*h {
		os.Exit(1)
	}
	g := NewGrid(w, h)
	for y := GridMargin; y < g.Height+GridMargin; y++ {
		for x := GridMargin; x < g.Width+GridMargin; x++ {
			g.Data[g.Stride*y+x] = data[g.Width*(y-GridMargin)+(x-GridMargin)]
		}
	}
	return g
}

func NewGridFromString(str string, w, h int) Grid {
	if len(str) != w*h {
		os.Exit(1)
	}
	g := NewGrid(w, h)
	for y := GridMargin; y < g.Height+GridMargin; y++ {
		for x := GridMargin; x < g.Width+GridMargin; x++ {
			d := 0
			if str[g.Width*(y-GridMargin)+(x-GridMargin)] == '1' {
				d = 1
			}
			g.Data[g.Stride*y+x] = d
		}
	}
	return g
}

func NewRandGrid(w, h int) Grid {
	g := NewGrid(w, h)
	g.Each(func(idx int) {
		if rand.Float64() > Threshold {
			g.Data[idx] = 1
		}
	})
	return g
}

func (g *Grid) Each(do func(idx int)) {
	for y := GridMargin; y < g.Height+GridMargin; y++ {
		for x := GridMargin; x < g.Width+GridMargin; x++ {
			do(g.Stride*y + x)
		}
	}
}

func NextGridFrom(h Grid) Grid {
	g := NewGrid(h.Width, h.Height)
	g.Each(func(idx int) {
		cnt := h.Data[idx-g.Stride-1] + h.Data[idx-g.Stride] + h.Data[idx-g.Stride+1] + h.Data[idx-1] + h.Data[idx+1] + h.Data[idx+g.Stride-1] + h.Data[idx+g.Stride] + h.Data[idx+g.Stride+1]
		switch {
		case h.Data[idx] == 0 && cnt == 3:
			g.Data[idx] = 1
		case h.Data[idx] == 1 && (cnt <= 1 || cnt >= 4):
			g.Data[idx] = 0
		default:
			g.Data[idx] = h.Data[idx]
		}
	})
	return g
}

func (g Grid) String() string {
	s := ""
	for y := GridMargin; y < g.Height+GridMargin; y++ {
		for x := GridMargin; x < g.Width+GridMargin; x++ {
			c := "□"
			if g.Data[g.Stride*y+x] > 0 {
				c = "■"
			}
			s += fmt.Sprint(c)
		}
		s += "\n"
	}
	return s[:len(s)-1]
}
