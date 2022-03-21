package pkg

import (
	"regexp"
	"strings"
)

// Returns filtered test objects
func FilterTestObjects(testObjects []TestObject, start, end uint64, pattern string, services string, envs string) []TestObject {

	var resultTestObjects []TestObject

	for _, testObject := range testObjects {
		matches, _ := regexp.MatchString(pattern, testObject.FullName)
		service, _ := GetServiceAndClass(testObject.FullName)
		meetEnv := (len(testObject.ParameterValues) > 0 && contains(envs, getEnv(strings.Join(testObject.ParameterValues, ""))))
		if (len(services) == 0 || contains(services, service)) && (pattern == "" || matches) && meetEnv {
			testObject.Extra.History.Items = filterItems(testObject.Extra.History.Items, start, end)
			resultTestObjects = append(resultTestObjects, testObject)
		}
	}
	return resultTestObjects

}

func contains(s string, str string) bool {
	if s == "" {
		return true
	}
	sArr := strings.Split(s, ",")
	for _, v := range sArr {
		if v == str {
			return true
		}
	}

	return false
}

// Filter items by start/end time
func filterItems(items []TestRunItem, start, end uint64) []TestRunItem {
	var result []TestRunItem
	for _, item := range items {
		if item.Time.Start >= start && (item.Time.End <= end || end == 0) {
			result = append(result, item)
		}
	}
	return result
}

func getEnv(parameterValues string) string {
	r := regexp.MustCompile(`.+env=(?P<environmet>[a-z]+).+`)
	matches := r.FindStringSubmatch(parameterValues)
	envIndex := r.SubexpIndex("environmet")
	return matches[envIndex]
}
