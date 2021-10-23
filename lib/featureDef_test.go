package lib

import (
	"os"
	"testing"
)

func TestFeatureDef_example1(t *testing.T) {

	t.Run("test config1", func(t *testing.T) {
		dat, err := os.ReadFile("../examples/config1.yml")
		if err != nil {
			t.Errorf("error opening file1")
		}
		var fdef FeatureDefinition

		err = fdef.LoadFromYaml(string(dat))
		if err != nil {
			t.Errorf("parsing file1")
		}

	})

}
