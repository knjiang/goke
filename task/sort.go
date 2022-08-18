package task

import (
	"fmt"
	"strings"
)

func sortTasksToRun(allTasks []Task, requiredTaskNames []string) ([][]Task, error) {
	graph, err := buildGraph(allTasks, requiredTaskNames)
	if err != nil {
		return nil, err
	}

	result, err := findTaskGroups(graph)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildGraph(allTasks []Task, requiredTaskNames []string) ([]*graphNode, error) {
	allTasksMap := make(map[string]Task)
	for _, t := range allTasks {
		allTasksMap[strings.ToLower(t.Name())] = t
	}

	var g []*graphNode
	seenTasks := make(map[string]struct{})
	for len(requiredTaskNames) > 0 {
		taskName := requiredTaskNames[0]
		requiredTaskNames = requiredTaskNames[1:]

		task, ok := allTasksMap[strings.ToLower(taskName)]
		if !ok {
			return nil, fmt.Errorf("unknown task '%s'", taskName)
		}

		if _, ok := seenTasks[task.Name()]; !ok {
			seenTasks[task.Name()] = struct{}{}
			g = append(g, &graphNode{task: task, edges: task.Dependencies()})
			requiredTaskNames = append(requiredTaskNames, task.Dependencies()...)
		}
	}

	return g, nil
}

type graphNode struct {
	task  Task
	edges []string
}

func toposort(g []*graphNode) ([]Task, error) {
	var queue []*graphNode
	for _, n := range g {
		if len(n.edges) == 0 {
			queue = append(queue, n)
		}
	}

	var sorted []Task
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		sorted = append(sorted, n.task)
		for _, m := range g {
			for i := range m.edges {
				if m.edges[i] == n.task.Name() {
					m.edges = append(m.edges[:i], m.edges[i+1:]...)
					if len(m.edges) == 0 {
						queue = append(queue, m)
					}
					break
				}
			}
		}
	}

	err := checkForCycles(g)
	if err != nil {
		return nil, err
	}

	return sorted, nil
}

func findTaskGroups(g []*graphNode) ([][]Task, error) {
	taskLevels := make(map[*graphNode]int)
	var queue []*graphNode
	for _, n := range g {
		if len(n.edges) == 0 {
			queue = append(queue, n)
			taskLevels[n] = 0
		}
	}

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, m := range g {
			for i := range m.edges {
				if m.edges[i] == n.task.Name() {
					m.edges = append(m.edges[:i], m.edges[i+1:]...)
					if len(m.edges) == 0 {
						queue = append(queue, m)
					}
					if taskLevels[m] < taskLevels[n]+1 {
						taskLevels[m] = taskLevels[n] + 1
					}
					break
				}
			}
		}
	}

	err := checkForCycles(g)
	if err != nil {
		return nil, err
	}

	return createTaskGroups(taskLevels), nil
}

func createTaskGroups(taskLevels map[*graphNode]int) [][]Task {
	taskMap := make(map[int][]Task)
	maxLevel := 0
	for t, level := range taskLevels {
		taskMap[level] = append(taskMap[level], t.task)
		if maxLevel < level {
			maxLevel = level
		}
	}
	taskGroups := make([][]Task, maxLevel+1)
	for i := 0; i < maxLevel+1; i++ {
		taskGroups[i] = taskMap[i]
	}

	return taskGroups
}

func checkForCycles(g []*graphNode) error {
	for _, n := range g {
		if len(n.edges) > 0 {
			return fmt.Errorf("a cycle exists")
		}
	}
	return nil
}
