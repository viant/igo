package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Algorithms(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected map[string]int
	}{
		{
			name: "copy_slice_index_assign",
			code: `
				src := []int{9, 7, 5, 3, 1}
				dst := make([]int, len(src))
				for i := 0; i < len(src); i++ {
					dst[i] = src[i]
				}
				ok := 1
				sum := 0
				for i := 0; i < len(src); i++ {
					sum += dst[i]
					if dst[i] != src[i] {
						ok = 0
					}
				}
			`,
			expected: map[string]int{"ok": 1, "sum": 25},
		},
		{
			name: "insertion_sort",
			code: `
				data := []int{7, 1, 9, 3, 5, 2, 8, 4, 6}
				n := len(data)
				for i := 1; i < n; i++ {
					key := data[i]
					j := i - 1
					for j >= 0 && data[j] > key {
						data[j + 1] = data[j]
						j--
					}
					data[j + 1] = key
				}
				ok := 1
				sum := 0
				for i := 0; i < n; i++ {
					sum += data[i]
					if i > 0 && data[i-1] > data[i] {
						ok = 0
					}
				}
			`,
			expected: map[string]int{"ok": 1, "sum": 45},
		},
		{
			name: "binary_search_iterative",
			code: `
				data := []int{1, 3, 5, 7, 9, 11, 13, 15}
				target := 11
				low := 0
				high := len(data) - 1
				idx := -1
				for low <= high {
					mid := (low + high) / 2
					if data[mid] == target {
						idx = mid
						break
					}
					if data[mid] < target {
						low = mid + 1
					} else {
						high = mid - 1
					}
				}
			`,
			expected: map[string]int{"idx": 5},
		},
		{
			name: "merge_sort_iterative_small",
			code: `
				src := []int{19, 7, 3, 17, 5, 13, 11, 2, 23, 29, 31, 37, 41, 43, 47, 53}
				data := make([]int, len(src))
				tmp := make([]int, len(src))

				n := len(src)
				for i := 0; i < n; i++ {
					data[i] = src[i]
				}

				width := 1
				left := 0
				mid := 0
				right := 0
				i := 0
				j := 0
				k := 0
				p := 0
				prev := 0
				for width < n {
					left = 0
					for left < n {
						mid = left + width
						if mid > n {
							mid = n
						}
						right = mid + width
						if right > n {
							right = n
						}
						i = left
						j = mid
						k = left
						for i < mid && j < right {
							if data[i] <= data[j] {
								tmp[k] = data[i]
								i++
							} else {
								tmp[k] = data[j]
								j++
							}
							k++
						}
						for i < mid {
							tmp[k] = data[i]
							i++
							k++
						}
						for j < right {
							tmp[k] = data[j]
							j++
							k++
						}
						p = left
						for p < right {
							data[p] = tmp[p]
							p++
						}
						left = right
					}
					width += width
				}

				ok := 1
				sum := 0
				if n > 0 {
					prev = data[0]
				}
				i = 0
				for i < n {
					sum += data[i]
					if i > 0 && prev > data[i] {
						ok = 0
					}
					prev = data[i]
					i++
				}
			`,
			expected: map[string]int{"ok": 1, "sum": 381},
		},
		{
			name: "bfs_adjacency_matrix",
			code: `
				n := 6
				adj := []int{
					0,1,1,0,0,0,
					1,0,0,1,0,0,
					1,0,0,0,1,0,
					0,1,0,0,0,0,
					0,0,1,0,0,1,
					0,0,0,0,1,0,
				}
				visited := make([]int, n)
				queue := make([]int, n)
				head := 0
				tail := 0
				visited[0] = 1
				queue[tail] = 0
				tail++
				count := 0
				score := 0
				step := 1
				for head < tail {
					v := queue[head]
					head++
					count++
					score += v * step
					step++
					for u := 0; u < n; u++ {
						if adj[v*n+u] == 1 && visited[u] == 0 {
							visited[u] = 1
							queue[tail] = u
							tail++
						}
					}
				}
			`,
			expected: map[string]int{"count": 6, "score": 70},
		},
		{
			name: "dfs_iterative_stack",
			code: `
				n := 6
				adj := []int{
					0,1,1,0,0,0,
					1,0,0,1,0,0,
					1,0,0,0,1,0,
					0,1,0,0,0,0,
					0,0,1,0,0,1,
					0,0,0,0,1,0,
				}
				visited := make([]int, n)
				stack := make([]int, n*n)
				top := 0
				stack[top] = 0
				top++
				count := 0
				score := 0
				step := 1
				for top > 0 {
					top--
					v := stack[top]
					if visited[v] == 1 {
						continue
					}
					visited[v] = 1
					count++
					score += v * step
					step++
					for u := n - 1; u >= 0; u-- {
						if adj[v*n+u] == 1 && visited[u] == 0 {
							stack[top] = u
							top++
						}
					}
				}
			`,
			expected: map[string]int{"count": 6, "score": 69},
		},
		{
			name: "dijkstra_lite_matrix",
			code: `
				n := 5
				inf := 1000000
				w := []int{
					0,2,5,inf,inf,
					2,0,1,2,inf,
					5,1,0,3,5,
					inf,2,3,0,1,
					inf,inf,5,1,0,
				}
				dist := make([]int, n)
				used := make([]int, n)
				for i := 0; i < n; i++ {
					dist[i] = inf
				}
				dist[0] = 0
				for iter := 0; iter < n; iter++ {
					v := -1
					for i := 0; i < n; i++ {
						if used[i] == 0 {
							if v == -1 || dist[i] < dist[v] {
								v = i
							}
						}
					}
					if v == -1 {
						break
					}
					used[v] = 1
					for u := 0; u < n; u++ {
						wt := w[v*n+u]
						if wt < inf {
							nd := dist[v] + wt
							if nd < dist[u] {
								dist[u] = nd
							}
						}
					}
				}
				d4 := dist[4]
				sum := 0
				for i := 0; i < n; i++ {
					sum += dist[i]
				}
			`,
			expected: map[string]int{"d4": 5, "sum": 14},
		},
		{
			name: "heap_push_pop_min",
			code: `
				data := []int{7,1,9,3,5,2,8,4,6}
				heap := make([]int, 32)
				size := 0
				for i := 0; i < len(data); i++ {
					x := data[i]
					heap[size] = x
					j := size
					size++
					for j > 0 {
						p := (j - 1) / 2
						if heap[p] <= heap[j] {
							break
						}
						tmp := heap[p]
						heap[p] = heap[j]
						heap[j] = tmp
						j = p
					}
				}
				ok := 1
				prev := -1
				sum := 0
				score := 0
				step := 1
				for size > 0 {
					minv := heap[0]
					sum += minv
					score += minv * step
					step++
					if minv < prev {
						ok = 0
					}
					prev = minv
					size--
					if size > 0 {
						heap[0] = heap[size]
						i := 0
						for {
							l := 2*i + 1
							if l >= size {
								break
							}
							r := l + 1
							c := l
							if r < size && heap[r] < heap[l] {
								c = r
							}
							if heap[i] <= heap[c] {
								break
							}
							tmp := heap[i]
							heap[i] = heap[c]
							heap[c] = tmp
							i = c
						}
					}
				}
			`,
			expected: map[string]int{"ok": 1, "sum": 45, "score": 285},
		},
		{
			name: "quickselect_kth",
			code: `
				data := []int{9,1,8,2,7,3,6,4,5}
				k := 4
				left := 0
				right := len(data) - 1
				ans := -1
				for left <= right {
					pivot := data[right]
					store := left
					for i := left; i < right; i++ {
						if data[i] < pivot {
							tmp := data[store]
							data[store] = data[i]
							data[i] = tmp
							store++
						}
					}
					tmp := data[store]
					data[store] = data[right]
					data[right] = tmp
					if store == k {
						ans = data[store]
						break
					}
					if store < k {
						left = store + 1
					} else {
						right = store - 1
					}
				}
			`,
			expected: map[string]int{"ans": 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			execution, err := scope.Compile(tc.code)
			if !assert.NoError(t, err) {
				return
			}
			state := execution.NewState()
			execution.Exec(state)
			for key, expected := range tc.expected {
				actual, err := state.Int(key)
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, expected, actual, key)
			}
		})
	}
}
