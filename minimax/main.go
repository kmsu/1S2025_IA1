package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var idCounter int = 0
var dotLines []string

// Lista de heurísticas fijas para los nodos hoja si las necesitara
var fixedHeuristics = []int{55, 12, 75, 20, 22, 30, 45, 51}
var heuristicIndex = 0

func inc() int {
	idCounter++
	return idCounter
}

// Nodo representa un nodo en el árbol: [id, heurística]
type Nodo struct {
	id         int
	heuristica int
}

func addNode(id int, label string, extra string) {
	line := fmt.Sprintf("    %d [label=\"%s\"%s];", id, label, extra)
	dotLines = append(dotLines, line)
}

func addEdge(from, to int, label string) {
	line := fmt.Sprintf("    %d -- %d [label=\"%s\"];", from, to, label)
	dotLines = append(dotLines, line)
}

func succ(node Nodo, depth int, bfactor int, arrow string) []Nodo {
	var children []Nodo
	for i := 0; i < bfactor; i++ {
		childID := inc()
		heuristic := 0
		if depth == 1 {

			//Asignar de forma aleatoria los valores de los nodos hojas
			//heuristic = rand.Intn(41) - 10 // [-10, 30]

			// Asignar valores fijos a los nodos hoja en lugar de random
			if heuristicIndex < len(fixedHeuristics) {
				heuristic = fixedHeuristics[heuristicIndex]
				heuristicIndex++
			}
		}
		addNode(childID, strconv.Itoa(heuristic), "")
		addEdge(node.id, childID, arrow)
		children = append(children, Nodo{id: childID, heuristica: heuristic})
	}
	return children
}

// minimax(root, 3, true, 2) -> 3 es la profundidad y 2 la cantidad de nodos hijos para cada nodo
func minimax(node Nodo, depth int, maximizing bool, bfactor int) Nodo {
	if depth == 0 {
		return node
	}

	arrow := "max"
	if !maximizing {
		arrow = "min"
	}

	children := succ(node, depth, bfactor, arrow)

	var best Nodo
	if maximizing {
		best = Nodo{id: -1, heuristica: -999}
		for _, child := range children {
			result := minimax(child, depth-1, false, bfactor)
			if result.heuristica > best.heuristica {
				best = result
			}
		}
	} else {
		best = Nodo{id: -1, heuristica: 999}
		for _, child := range children {
			result := minimax(child, depth-1, true, bfactor)
			if result.heuristica < best.heuristica {
				best = result
			}
		}
	}

	node.heuristica = best.heuristica
	addNode(best.id, strconv.Itoa(best.heuristica), ", shape=doublecircle")
	addNode(node.id, strconv.Itoa(node.heuristica), "")
	return node
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Iniciar DOT
	dotLines = append(dotLines, "graph G {")
	dotLines = append(dotLines, "    rankdir=TB;")
	dotLines = append(dotLines, "    node [shape=circle, fontname=\"Arial\", color=blue4, fontcolor=blue4];")
	dotLines = append(dotLines, "    edge [color=blue4];")

	rootID := inc()
	addNode(rootID, "0", ", color=red, fontcolor=red, shape=doublecircle")

	root := Nodo{id: rootID, heuristica: 0}

	//raiz, profundidad (niveles debajo de la raiz), maximizar/minimizar, numero de nodos sucesores
	//result := minimax(root, 4, true, 3) -> 4 niveles debajo de la raiz y cada uno con 3 nodos hijos para cada nodo.
	result := minimax(root, 3, true, 2)

	dotLines = append(dotLines, "}")
	fmt.Println("Resultado final:", result)

	// Imprimir el .dot generado
	fmt.Println("\n--- DOT Output ---")
	for _, line := range dotLines {
		fmt.Println(line)
	}

	//crear arbol (graficamente)
	dotContent := strings.Join(dotLines, "\n")

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

	_, err = file.WriteString(dotContent)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}

	file.Close()

	cmd := exec.Command("dot", "-Tpng", "dot.dot", "-o", "minimax.png")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error al generar el reporte PNG:", err)
		fmt.Println("Detalle del error:", string(out))
		return
	}

}
