//bestfirst

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

func successors(n []interface{}, end string, h int) [][]interface{} {
	suc := [][]interface{}{}
	s := n[0].(string)

	for i := 0; i < len(s)-1; i++ {
		tmp := string(s[i])
		child := s[:i] + string(s[i+1]) + tmp + s[i+2:]
		suc = append(suc, []interface{}{child, heuristic(child, end, h), inc()})
	}
	return suc
}

func bestfirst(start, end string, h int) string {
	dot := "graph G {\n"
	list := [][]interface{}{
		{start, heuristic(start, end, h), inc()},
	}
	dot += fmt.Sprintf("N%d [label=\"%s\"];\n", list[0][2], list[0][0])

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
			dot += fmt.Sprintf("N%d [label=\"%s\"];\n", val[2], val[0])
			dot += fmt.Sprintf("N%d -- N%d [label=\"%d\"];\n", current[2], val[2], val[1])
		}

		list = append(list, temp...)

		sort.Slice(list, func(i, j int) bool {
			return list[i][1].(int) < list[j][1].(int)
		})

		cont++
		if cont > 100 {
			fmt.Println("The search is looped!")
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
		fmt.Println("Ejemplo: go run main.go halo hola 1")
		fmt.Println("1 = casillas fuera de lugar")
		fmt.Println("2 = manhatan")
		fmt.Println("3 = esperanza")
		input = "halo hola 1"
	} else {
		input = strings.Join(args[1:], " ")
	}

	parts := strings.Split(input, " ")
	start := parts[0]
	end := parts[1]
	h, _ := strconv.Atoi(parts[2])

	result := bestfirst(start, end, h)
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
