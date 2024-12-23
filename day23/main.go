package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"bautik.net/advent2024/helpers"
)

func addConnection(comp1, comp2 string, connectionMap map[string]map[string]bool) {
	if _, ok := connectionMap[comp1]; !ok {
		connectionMap[comp1] = map[string]bool{}
	}
	if _, ok := connectionMap[comp2]; !ok {
		connectionMap[comp2] = map[string]bool{}
	}
	connectionMap[comp1][comp2] = true
	connectionMap[comp2][comp1] = true
}

func bronKerbosch(adjacencyMatrix [][]int, r, p, x []int, maximalCliques *[][]int) {
	if len(p) == 0 && len(x) == 0 {
		// No more candidates, add the current clique to the results
		*maximalCliques = append(*maximalCliques, append([]int(nil), r...))
		return
	}

	for _, v := range p {
		neighbors := []int{}
		for i := 0; i < len(adjacencyMatrix); i++ {
			if adjacencyMatrix[v][i] == 1 {
				neighbors = append(neighbors, i)
			}
		}

		bronKerbosch(adjacencyMatrix, append(r, v), helpers.Intersection(p, neighbors), helpers.Intersection(x, neighbors), maximalCliques)

		p = helpers.Difference(p, []int{v})
		x = helpers.Union(x, []int{v})
	}
}

func findLargestCliqueIndices(adjacencyMatrix [][]int) []int {
	n := len(adjacencyMatrix)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	var cliques [][]int
	bronKerbosch(adjacencyMatrix, []int{}, p, []int{}, &cliques)

	largest := []int{}
	for _, clique := range cliques {
		if len(clique) > len(largest) {
			largest = clique
		}
	}

	return largest
}

func main() {
	duration := helpers.MeasureRuntime(func() {
		fileBytes, _ := os.ReadFile("input.txt")
		connections := map[string]map[string]bool{}
		compIndices := map[string]int{}
		currentIndex := 0

		for _, l := range strings.Split(string(fileBytes), "\n") {
			comps := strings.Split(l, "-")
			addConnection(comps[0], comps[1], connections)
			if _, ok := compIndices[comps[0]]; !ok {
				compIndices[comps[0]] = currentIndex
				currentIndex++
			}
			if _, ok := compIndices[comps[1]]; !ok {
				compIndices[comps[1]] = currentIndex
				currentIndex++
			}
		}
		idxComp := make([]string, len(compIndices))
		adjMatrix := make([][]int, len(compIndices))
		for c, i := range compIndices {
			idxComp[i] = c
			adjMatrix[i] = make([]int, len(compIndices))
		}
		for comp1, connsComp := range connections {
			for peer := range connsComp {
				adjMatrix[compIndices[comp1]][compIndices[peer]] = 1
			}
		}

		triangles := map[[3]string]bool{}
		for comp, i := range compIndices {
			if comp[0] == 't' {
				for n := range len(adjMatrix) {
					if adjMatrix[i][n] == 1 {
						for j := n + 1; j < len(adjMatrix); j++ {
							if adjMatrix[n][j] == 1 && adjMatrix[i][j] == 1 {
								triangle := []string{
									idxComp[n],
									idxComp[i],
									idxComp[j],
								}
								slices.Sort(triangle)
								triangles[[3]string(triangle)] = true
							}
						}
					}
				}
			}
		}

		fmt.Println("Part 1:", len(triangles))
		largestGroup := findLargestCliqueIndices(adjMatrix)
		members := []string{}
		for _, l := range largestGroup {
			members = append(members, idxComp[l])
		}
		slices.Sort(members)
		fmt.Println("Part 2:", strings.Join(members, ","))
	})

	fmt.Println("Took", duration.Microseconds(), "μs")
}
