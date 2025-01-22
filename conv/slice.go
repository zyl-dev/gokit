package conv

import "strconv"

func RemoveDuplicateElement(list []string) []string {
	result := make([]string, 0, len(list))
	temp := map[string]struct{}{}
	for _, item := range list {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// IntsToStrings converts an int slice to string slice
func IntsToStrings(input []int) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = strconv.Itoa(v)
	}
	return output
}
