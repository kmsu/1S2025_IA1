package main

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// --- Estructuras ---
type Edge struct {
	to       string
	distance int
	heur     int
}

type Node struct {
	id       int
	name     string
	parentID int
	g        int
	f        int
	path     []string
}

type Item struct {
	nodeID   int
	name     string
	priority int
	g        int
	path     []string
	parentID int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return pq[i].index < pq[j].index
	}
	return pq[i].priority < pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i]; pq[i].index, pq[j].index = i, j }
func (pq *PriorityQueue) Push(x any)   { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// --- Variables globales ---
var idCounter = 0

func nextID() int {
	idCounter++
	return idCounter
}

// --- Leer grafo y heurísticas desde CSV ---
func buildGraphFromCSV(filename string, destination string) (map[string][]Edge, map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	graph := make(map[string][]Edge)
	heuristics := make(map[string]int)

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		origin := row[1]
		dest := row[2]
		dist, _ := strconv.Atoi(row[3])
		heur, _ := strconv.Atoi(row[4])

		graph[origin] = append(graph[origin], Edge{to: dest, distance: dist, heur: heur})

		if dest == destination {
			if val, ok := heuristics[origin]; !ok || heur < val {
				heuristics[origin] = heur
			}
		}
	}
	heuristics[destination] = 0
	return graph, heuristics, nil
}

// --- A* con generación del árbol de búsqueda ---
func aStarTree(graph map[string][]Edge, heuristics map[string]int, start, goal string) ([]Node, []int) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	startID := nextID()
	heap.Push(&pq, &Item{
		nodeID:   startID,
		name:     start,
		priority: heuristics[start],
		g:        0,
		path:     []string{start},
		parentID: 0,
	})

	var nodes []Node
	var goalPath []int

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)

		nodes = append(nodes, Node{
			id:       current.nodeID,
			name:     current.name,
			parentID: current.parentID,
			g:        current.g,
			f:        current.priority,
			path:     current.path,
		})

		if current.name == goal {
			// reconstruir path por ids
			for i := len(nodes) - 1; i >= 0; i-- {
				if nodes[i].name == goal && nodes[i].f == current.priority {
					n := nodes[i]
					for n.parentID != 0 {
						goalPath = append([]int{n.id}, goalPath...)
						// buscar padre
						for _, x := range nodes {
							if x.id == n.parentID {
								n = x
								break
							}
						}
					}
					goalPath = append([]int{n.id}, goalPath...)
					break
				}
			}
			break
		}

		for _, edge := range graph[current.name] {
			newG := current.g + edge.distance
			newF := newG + heuristics[edge.to]
			newPath := append([]string{}, current.path...)
			newPath = append(newPath, edge.to)

			childID := nextID()
			heap.Push(&pq, &Item{
				nodeID:   childID,
				name:     edge.to,
				priority: newF,
				g:        newG,
				path:     newPath,
				parentID: current.nodeID,
			})
		}
	}
	return nodes, goalPath
}

// --- Generar .dot del árbol de búsqueda ---
func generateSearchTreeDOT(nodes []Node, pathIDs []int) {
	dot := "digraph G {\n"
	nodeMap := make(map[int]Node)
	highlight := make(map[int]bool)

	for _, id := range pathIDs {
		highlight[id] = true
	}

	for _, n := range nodes {
		nodeMap[n.id] = n
		label := fmt.Sprintf("%s\\nf=%d", n.name, n.f)
		style := ""
		if highlight[n.id] {
			style = ", color=red, fontcolor=red, style=filled, fillcolor=lightyellow"
		}
		dot += fmt.Sprintf("  %d [label=\"%s\"%s];\n", n.id, label, style)
		if n.parentID != 0 {
			if highlight[n.id] && highlight[n.parentID] {
				dot += fmt.Sprintf("  %d -> %d [color=red, penwidth=2.0];\n", n.parentID, n.id)
			} else {
				dot += fmt.Sprintf("  %d -> %d;\n", n.parentID, n.id)
			}
		}
	}
	dot += "}\n"

	err := os.WriteFile("arbol_busqueda.dot", []byte(dot), 0644)
	if err != nil {
		fmt.Println("Error al guardar el archivo .dot:", err)
	} else {
		fmt.Println("Archivo arbol_busqueda.dot generado correctamente.")
	}
}

// --- main ---
func main() {
	graph, heuristics, err := buildGraphFromCSV("tabla.csv", "E")
	if err != nil {
		fmt.Println("Error leyendo el CSV:", err)
		return
	}

	nodes, pathIDs := aStarTree(graph, heuristics, "A", "E")
	if len(pathIDs) == 0 {
		fmt.Println("No se encontró ruta.")
		return
	}

	// Mostrar ruta textual
	var path []string
	for _, id := range pathIDs {
		for _, n := range nodes {
			if n.id == id {
				path = append(path, n.name)
				break
			}
		}
	}
	fmt.Println("Ruta encontrada:", path)

	lastID := pathIDs[len(pathIDs)-1]
	var finalF int
	for _, n := range nodes {
		if n.id == lastID {
			finalF = n.f
			break
		}
	}
	fmt.Println("f(x) mínimo:", finalF)

	var finalG int
	for _, n := range nodes {
		if n.id == lastID {
			finalG = n.g
			break
		}
	}
	fmt.Println("Costo total g(x):", finalG)

	generateSearchTreeDOT(nodes, pathIDs)
}
