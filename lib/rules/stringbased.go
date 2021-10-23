package rules

import (
	"container/list"
)

type StringArrayBased struct {
}

func (me *StringArrayBased) CheckRule(potentialValues []string, currentValue string, matches *list.List, total *int, label string) (bool, bool) {

	if len(potentialValues) > 0 {
		ok := me.evaluate(potentialValues, currentValue)
		if ok {
			matches.PushBack(label)
			*total++
		}
		return ok, false
	}
	return false, true
}

func (me *StringArrayBased) evaluate(potentialValues []string, currentValue string) bool {
	return Include(potentialValues, currentValue)
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
