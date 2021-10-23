package lib

import (
	"time"
)

type testCase struct {
	key      string
	value    *FeatureDefinition
	duration time.Duration
	want     *FeatureDefinition
	found    bool
}

func getCommonScenarios() []testCase {

	var feDef FeatureDefinition

	feDef.Name = "name"
	feDef.Meta.Description = "test description"
	feDef.Default.Value.Data = "value123"
	feDef.Default.Value.Type = ValueText
	feDef.Default.TTL = "3s"

	var tests = []testCase{
		{"key1", &feDef, 1 * time.Second, &feDef, true},      // retrieve value with corret TTL
		{"key2", &feDef, 400 * time.Millisecond, nil, false}, //retrieve value with TTL completed
	}

	return tests
}
