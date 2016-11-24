package utils

// Cartesian gets all combinations of values for keys
type Cartesian struct {
}

// GetCombinations returns an array of maps with all combinations built by a map with values
// "key1": ["1", "2"], "key2":["3", "4"] -> [{"key1":"1", "key2":"3"}, {"key1":"1", "key2":"4"}, {"key1":"2", "key2":"3"}, {"key1":"2", "key2":"4"}]
func (c Cartesian) GetCombinations(inputMap map[string][]string) []map[string]string {
	if len(inputMap) == 0 {
		return []map[string]string{}
	}
	inputArray := [][]string{}
	keys := []string{}
	for key, value := range inputMap {
		keys = append(keys, key)
		inputArray = append(inputArray, value)
	}

	arrayResult := c.CartesianDistribute(inputArray)

	results := []map[string]string{}

	for _, values := range arrayResult {
		currentMap := make(map[string]string)
		for index, value := range values {
			key := keys[index]
			currentMap[key] = value
		}
		results = append(results, currentMap)
	}

	return results
}

// CartesianDistribute getting full combinations distribution using cartesian combination algorithm - http://stackoverflow.com/a/15310051/613113
func (c Cartesian) CartesianDistribute(input [][]string) [][]string {
	result := [][]string{}

	c.cartesianHelper(&result, input, []string{}, 0)

	return result
}

func (c Cartesian) cartesianHelper(output *[][]string, input [][]string, current []string, i int) {
	max := len(input) - 1
	for j := 0; j < len(input[i]); j++ {
		a := append(current, input[i][j])
		if i == max {
			*output = append(*output, a)
		} else {
			c.cartesianHelper(output, input, a, i+1)
		}
	}
}
