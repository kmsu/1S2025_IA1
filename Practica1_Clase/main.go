package main

import (
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

func inc() int {
	id++
	return id - 1
}

func heuristic(start, end string, h int) int {
	switch h {
	case 1:
		// Casillas fuera de lugar
		tilesOut := 0
		for i := 0; i < len(start); i++ {
			if start[i] != '0' && start[i] != end[i] {
				tilesOut++
			}
		}
		return tilesOut
	case 2:
		// Manhattan
		man := 0
		for i := 0; i < len(start); i++ {
			if start[i] != '0' {
				target := strings.IndexRune(end, rune(start[i]))
				man += int(math.Abs(float64(i/3-target/3)) + math.Abs(float64(i%3-target%3)))
			}
		}
		return man
	case 3:
		// Esperanza
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

// Calcula las posiciones válidas a donde se puede mover el espacio vacío
func validMoves(pos int) []int {
	moves := []int{}
	row, col := pos/3, pos%3

	if row > 0 {
		moves = append(moves, pos-3) // up
	}
	if row < 2 {
		moves = append(moves, pos+3) // down
	}
	if col > 0 {
		moves = append(moves, pos-1) // left
	}
	if col < 2 {
		moves = append(moves, pos+1) // right
	}

	return moves
}

func successors(n []interface{}, end string, h int) [][]interface{} {
	succ := [][]interface{}{}
	s := n[0].(string)
	zeroPos := strings.Index(s, "0")

	for _, move := range validMoves(zeroPos) {
		r := []rune(s)
		r[zeroPos], r[move] = r[move], r[zeroPos]
		child := string(r)
		succ = append(succ, []interface{}{child, heuristic(child, end, h), inc()})
	}
	return succ
}

func bestfirst(start, end string, h int) string {
	dot := "graph G {\n"
	list := [][]interface{}{
		{start, heuristic(start, end, h), inc()},
	}
	dot += fmt.Sprintf("N%d [label=\"%s\"];\n", list[0][2], list[0][0])

	visited := map[string]bool{start: true}
	cont := 0
	for len(list) > 0 {
		current := list[0]
		list = list[1:]

		if current[0] == end {
			dot += "}"
			return dot
		}

		temp := successors(current, end, h)
		for _, val := range temp {
			if visited[val[0].(string)] {
				continue
			}
			visited[val[0].(string)] = true
			dot += fmt.Sprintf("N%d [label=\"%s\"];\n", val[2], val[0])
			dot += fmt.Sprintf("N%d -- N%d [label=\"%d\"];\n", current[2], val[2], val[1])
			list = append(list, val)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i][1].(int) < list[j][1].(int)
		})

		cont++
		if cont > 1000 {
			fmt.Println("La búsqueda se ha detenido (posible ciclo infinito o sin solución).")
			dot += "}"
			return dot
		}
	}

	dot += "}"
	return dot
}

func main() {
	args := os.Args
	var input string
	if len(args) < 4 {
		fmt.Println("Uso: go run main.go <inicio> <fin> <heuristica>")
		fmt.Println("Ejemplo: go run main.go 724506831 123456780 2")
		fmt.Println("1 = Casillas fuera de lugar")
		fmt.Println("2 = Manhattan")
		fmt.Println("3 = Esperanza")
		return
	} else {
		input = strings.Join(args[1:], " ")
	}

	parts := strings.Split(input, " ")
	start := parts[0]
	end := parts[1]
	h, _ := strconv.Atoi(parts[2])

	result := bestfirst(start, end, h)
	fmt.Println(result)

	// Exportar a Graphviz
	dir := filepath.Dir("./dot.dot")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el directorio:", err)
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
	switch h {
	case 1:
		salida = "TilesOut.png"
	case 2:
		salida = "Manhattan.png"
	default:
		salida = "Esperanza.png"
	}

	cmd := exec.Command("dot", "-Tpng", "dot.dot", "-o", salida)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error al generar el PNG:", err)
		fmt.Println("Detalle:", string(out))
	}
}
