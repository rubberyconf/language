package lib

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rubberyconf/language/lib/rules"
	"gopkg.in/yaml.v2"
)

type ValueFeatureType int

const (
	ValueText           ValueFeatureType = 0
	ValueNumber         ValueFeatureType = 1
	ValueBoolean        ValueFeatureType = 2
	ValueJsonFormat     ValueFeatureType = 3
	ValueExternResource ValueFeatureType = 4
)

func (s ValueFeatureType) String() string {
	switch s {
	case ValueText:
		return "text"
	case ValueNumber:
		return "number"
	case ValueBoolean:
		return "boolean"
	case ValueJsonFormat:
		return "jsonFormat"
	case ValueExternResource:
		return "externResource"
	}
	return "unknown"
}

type FeatureDefinition struct {
	Name string `yaml:"name" json:"name"`
	Meta struct {
		Description string   `yaml:"description" json:"description"`
		Owner       string   `yaml:"owner" json:"owner"`
		Tags        []string `yaml:"tags" json:"tags,omitempty"`
	} `yaml:"meta" json:"meta"`

	Default struct {
		Value struct {
			Data interface{}      `yaml:"data" json:"data"`
			Type ValueFeatureType `yaml:"type" json:"type"`
		} `yaml:"value" json:"value"`
		TTL string `yaml:"ttl" json:"ttl"`
	} `yaml:"default" json:"default"`

	Configurations []struct {
		ConfigId       string              `yaml:"id" json:"id"`
		RulesBehaviour string              `yaml:"rulesBehaviour" json:"rulesBehaviour"`
		Rules          []rules.FeatureRule `yaml:"rules" json:"rules,omitempty"`
		Value          interface{}         `yaml:"value" json:"value"`
		Rollout        struct {
			Strategy       string `yaml:"strategy" json:"strategy"`
			EnabledForOnly string `yaml:"enabledForOnly" json:"enabledForOnly"`
			Selector       string `yaml:"selector" json:"selector"`
		} `yaml:"rollout" json:"rollout"`
	} `yaml:"configurations" json:"configurations,omitempty"`
}

func (conf *FeatureDefinition) LoadFromYaml(payload interface{}) error {
	aux := fmt.Sprintf("%v", payload)
	decoder := yaml.NewDecoder(strings.NewReader(aux))
	err := decoder.Decode(conf)
	return err
}

func (conf *FeatureDefinition) LoadFromJsonBinary(b []byte) error {

	err := json.Unmarshal(b, &conf)
	return err
}

func (conf *FeatureDefinition) LoadFromString(text string) error {

	err := yaml.Unmarshal([]byte(text), &conf)
	if err != nil {
		//logs.GetLogs().WriteMessage(logs.ERROR, "error unmarshalling yaml content to featureDefinition", nil)
		return err
	}
	return nil
}
func (conf *FeatureDefinition) ToString() (string, error) {

	b, err := yaml.Marshal(conf)
	if err != nil {
		return "", err
	}
	sb := string(b)
	return sb, nil
}

func (conf *FeatureDefinition) GetFinalValue(vars map[string]string) (interface{}, error) {

	var afterCast interface{}

	data, found, _, _ /*confId, matches*/ := conf.SelectRule(vars)
	if !found {
		data = conf.Default.Value.Data
	}

	switch conf.Default.Value.Type {

	case ValueText:
		afterCast = data.(string)
	case ValueJsonFormat:
		b, err := json.MarshalIndent(data, "", "   ")
		if err != nil {
			return nil, err
		}
		afterCast = string(b)
	case ValueNumber:
		afterCast = data.(int)
	case ValueExternResource:
		afterCast = nil
	}
	return afterCast, nil
}

func (conf *FeatureDefinition) SelectRule(vars map[string]string) (interface{}, bool, string, *list.List) {

	ruleMast := NewRuleMaster()

	for _, c := range conf.Configurations {
		total := 0
		totalMatches := list.New()
		for _, r := range c.Rules {
			matches, labelMatches := ruleMast.CheckRules(r, vars)
			total += matches
			totalMatches.PushBackList(labelMatches)
		}

		logic := c.RulesBehaviour
		if logic == "" {
			logic = "AND"
		}
		if logic == "OR" && total > 1 {
			return c.Value, true, c.ConfigId, totalMatches

		} else if logic == "AND" && total == len(c.Rules) {
			return c.Value, true, c.ConfigId, totalMatches
		}
	}
	return nil, false, "", nil
}
