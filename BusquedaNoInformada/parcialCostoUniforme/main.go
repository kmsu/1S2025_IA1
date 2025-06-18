package main

import (
	"fmt"
	"sort"
)

/*
Grafo: (A,B,2)(A,C,3)(A,D,4)(B,C,1)(C,B,1)(B,D,3)(C,D,2)
*/
func successors(node []int) [][]int {
	name := node[0]
	cost := node[1]

	switch name {
	case 1:
		return [][]int{
			{2, cost + 2},
			{3, cost + 3},
			{4, cost + 4},
		}
	case 2:
		return [][]int{
			{3, cost + 1},
			{4, cost + 3},
		}
	case 3:
		return [][]int{
			{2, cost + 1},
			{4, cost + 2},
		}

	}
	return [][]int{}
}

func uniformCost(begin, end int) {
	list := [][]int{{begin, 0}}

	for len(list) > 0 {
		current := list[0]
		list = list[1:]
		fmt.Println("Current Node: ", current)

		if current[0] == end {
			fmt.Println("SOLUTION")
			return
		}

		tmp := successors(current)
		fmt.Println("Successors:", tmp)

		if len(tmp) > 0 {
			list = append(list, tmp...)
			sort.Slice(list, func(i, j int) bool {
				return list[i][1] < list[j][1]
			})
			fmt.Println("New List: ", list)
		}
	}

	fmt.Println("NO-SOLUTION")
}

func main() {
	uniformCost(1, 4)
}
