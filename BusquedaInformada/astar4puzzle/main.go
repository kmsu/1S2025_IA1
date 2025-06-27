package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var id = 1

func heuristic(start, end string, h int) int {
	if h == 1 {
		//Casillas fuera de lugar
		tilesOut := 0
		for i := 0; i < len(start); i++ {
			if start[i] != end[i] {
				tilesOut++
			}
		}
		return tilesOut
	} else if h == 2 {
		//manhatan
		man := 0
		for i := 0; i < len(start); i++ {
			c := string(start[i])
			pos := strings.Index(end, c)
			if pos >= 0 {
				man += int(math.Abs(float64(i - pos)))
			}
		}
		return man
	} else if h == 3 {
		//esperanza (diferencia absoluta entre el estado final e inicial)
		esperanza := 0
		for i := 0; i < len(start); i++ {
			valStart := int(start[i])
			valEnd := int(end[i])
			esperanza += int(math.Abs(float64(valStart - valEnd)))
		}
		return esperanza
	}
	return 0
}

func inc() int {
	id++
	return id
}

type Node struct {
	state string
	cost  int
	id    int
	level int
}

func successors(n Node, end string, h int) []Node {
	var suc []Node
	for i := 0; i < len(n.state)-1; i++ {
		level := n.level + 1
		chars := []rune(n.state)
		chars[i], chars[i+1] = chars[i+1], chars[i]
		child := string(chars)
		suc = append(suc, Node{
			state: child,
			cost:  heuristic(child, end, h) + level,
			id:    inc(),
			level: level,
		})
	}
	return suc
}

func aStar(start, end string, h int) string {
	dot := "graph{"
	list := []Node{{
		state: start,
		cost:  heuristic(start, end, h),
		id:    id,
		level: 0,
	}}
	dot += fmt.Sprintf("\n%d [label=\"%s\"];\n", list[0].id, list[0].state)

	cont := 0 //contador de iteraciones
	for len(list) > 0 {
		current := list[0]
		list = list[1:]

		if current.state == end {
			dot += "}"
			return dot
		}

		temp := successors(current, end, h)
		for _, val := range temp {
			dot += fmt.Sprintf("%d [label=\"%s\"]; ", val.id, val.state)
			dot += fmt.Sprintf("%d--%d [label=\"%d+%d\"];\n", current.id, val.id, val.cost-val.level, val.level)
		}

		list = append(list, temp...)
		sort.Slice(list, func(i, j int) bool { return list[i].cost < list[j].cost })

		cont++
		if cont > 150 {
			fmt.Println("The search is looped!")
			dot += "}"
			return dot
		}
	}
	dot += "}"
	return dot
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter start text and end text separated by a space: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = "halo hola 1"
	}
	parts := strings.Split(input, " ")
	if len(parts) < 3 {
		fmt.Println("Invalid entry.")
		return
	}
	start := parts[0]
	end := parts[1]
	h, _ := strconv.Atoi(parts[2])

	result := aStar(start, end, h)
	fmt.Println(result)

	//crear arbol (graficamente)

	dir := filepath.Dir("./dot.dot")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el directorio: ", err)
		return
	}

	file, err := os.Create("./dot.dot")
	if err != nil {
		fmt.Println("Error al crear el archivo .dot:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(result)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}

	file.Close()

	salida := ""
	if h == 1 {
		salida = "TilesOut.png"
	} else if h == 2 {
		salida = "Manhatan.png"
	} else {
		salida = "Esperanza.png"
	}

	cmd := exec.Command("dot", "-Tpng", "dot.dot", "-o", salida)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error al generar el reporte PNG:", err)
		fmt.Println("Detalle del error:", string(out))
		return
	}
}
