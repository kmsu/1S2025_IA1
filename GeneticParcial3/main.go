package main

import (
	"fmt"
	"math"
	"strconv"
)

// Función objetivo: minimizar f(x) = 25x - x^2
func fitness(x int) int {
	return 25*x - x*x
}

// Convertir entero a binario de 5 bits
func toBinary5(x int) string {
	bin := strconv.FormatInt(int64(x), 2)
	for len(bin) < 5 {
		bin = "0" + bin
	}
	return bin
}

// Cruce zigzag entre dos padres
func zigzagCross(p1, p2 int) (int, int) {
	b1 := toBinary5(p1)
	b2 := toBinary5(p2)
	child1 := ""
	child2 := ""
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			child1 += string(b1[i])
			child2 += string(b2[i])
		} else {
			child1 += string(b2[i])
			child2 += string(b1[i])
		}
	}
	// Convertir binarios a enteros
	c1, _ := strconv.ParseInt(child1, 2, 64)
	c2, _ := strconv.ParseInt(child2, 2, 64)
	return int(c1), int(c2)
}

// Promedio
func average(arr []int) float64 {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return float64(sum) / float64(len(arr))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

var allFitness []int //para calculo de fitnes mas alto de las 4 generaciones
func main() {
	// Población inicial
	population := []int{22, 24, 4, 27}
	var fitnessHistory []float64
	genCount := 1
	var generation2 []int

	bestFitnessMaximizar := math.MinInt
	bestFitnessMinimizar := math.MaxInt64

	bestIndividualMax := -1
	bestIndividualMin := -1

	for {
		// Calcular fitness
		fit := []int{
			fitness(population[0]),
			fitness(population[1]),
			fitness(population[2]),
			fitness(population[3]),
		}

		//maximizando
		for _, x := range population {
			f := fitness(x)
			if f > bestFitnessMaximizar {
				bestFitnessMaximizar = f
				bestIndividualMax = x
			}
		}

		//minimizando
		for _, x := range population {
			f := fitness(x)
			if f < bestFitnessMinimizar {
				bestFitnessMinimizar = f
				bestIndividualMin = x
			}
		}

		fmt.Printf("\n🔁 Generación %d:\n", genCount+1)
		for i, x := range population {
			fmt.Printf("Individuo %d: x = %d → f(x) = %d\n", i+1, x, fitness(x))
		}

		// Guardar todos los valores de fitness de todas las generaciones
		allFitness = append(allFitness, fit...)

		// Guardar promedio
		avg := average(fit)
		fitnessHistory = append(fitnessHistory, avg)

		if genCount == 2 {
			generation2 = append(generation2, population...)
		}

		// Verificar convergencia si hay al menos 2 generaciones
		if len(fitnessHistory) > 1 {
			prev := fitnessHistory[len(fitnessHistory)-2]
			curr := fitnessHistory[len(fitnessHistory)-1]
			conv := min(prev, curr) / max(prev, curr)
			//NOTA no pueden haber convergencias negativas, se trabaja entre 0 y  1
			if conv >= 0.7 {
				break
			}
		}

		// Torneos (1,2) y (3,4)
		var p1, p2 int
		if fitness(population[0]) < fitness(population[1]) {
			p1 = population[0]
		} else {
			p1 = population[1]
		}
		if fitness(population[2]) < fitness(population[3]) {
			p2 = population[2]
		} else {
			p2 = population[3]
		}

		// Zigzag
		h1, h2 := zigzagCross(p1, p2)

		// Reemplazo: mejor padre, peor hijo, mejor hijo, peor padre
		parents := []int{p1, p2}
		children := []int{h1, h2}

		// Ordenar padres
		if fitness(parents[0]) > fitness(parents[1]) {
			parents[0], parents[1] = parents[1], parents[0]
		}
		// Ordenar hijos
		if fitness(children[0]) > fitness(children[1]) {
			children[0], children[1] = children[1], children[0]
		}

		// Nueva población: mejor padre, peor hijo, mejor hijo, peor padre
		population = []int{parents[0], children[1], children[0], parents[1]}
		genCount++
	}

	// Calcular porcentaje de promedios negativos
	negativos := 0
	for _, promedio := range fitnessHistory {
		if promedio < 0 {
			negativos++
		}
	}

	total := len(fitnessHistory)
	porcentaje := float64(negativos) / float64(total)

	fmt.Printf("\nTotal de promedios: %d\n", total)
	fmt.Printf("Promedios negativos: %d\n", negativos)
	fmt.Printf("Porcentaje de promedios negativos: %.2f\n", porcentaje)

	// Calcular valores de convergencia positivos
	// Convergencias se calculan desde la segunda generación en adelante
	positiveConvergenceCount := 0
	for i := 1; i < len(fitnessHistory); i++ {
		prev := fitnessHistory[i-1]
		curr := fitnessHistory[i]
		convergencia := min(prev, curr) / max(prev, curr)
		if convergencia > 0 {
			positiveConvergenceCount++
		}
	}

	fmt.Printf("\nValores de convergencia positivos: %d\n", positiveConvergenceCount)

	// Calcular el fitness más alto entre todas las generaciones
	maxFitness := allFitness[0]
	for _, f := range allFitness {
		if f > maxFitness {
			maxFitness = f
		}
	}
	fmt.Printf("\n✅ El valor de fitness más alto entre las 4 generaciones fue: %d\n", maxFitness)

	fmt.Printf("\n🏆 Individuo más apto al finalizar (maximizando): x = %d → f(x) = %d\n", bestIndividualMax, bestFitnessMaximizar)
	fmt.Printf("\n🏆 Individuo más apto al finalizar (minimizando): x = %d → f(x) = %d\n", bestIndividualMin, bestFitnessMinimizar)

}
