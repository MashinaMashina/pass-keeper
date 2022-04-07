package test

import (
	"encoding/json"
	"fmt"
	"pass-keeper/internal/config"
	"reflect"
	"testing"
)

func TestConfigTypes(t *testing.T) {
	str := "{\"a\":123, \"b\":\"тест\", \"c\":[\"one\", \"two\"], \"d\":{\"e\":\"1\", \"f\": \"2\"}}"

	t.Log("init config")
	cfg := config.NewConfig()
	err := cfg.InitFromData([]byte(str))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("get string parameter")
	if err = equal("a", cfg.String("a"), "123"); err != nil {
		t.Error(err)
		return
	}

	t.Log("get cyrillic string parameter")
	if err = equal("b", cfg.String("b"), "тест"); err != nil {
		t.Error(err)
		return
	}

	t.Log("get slice parameter")
	if err = equal("c", cfg.Slice("c"), []string{"one", "two"}); err != nil {
		t.Error(err)
		return
	}

	t.Log("get map parameter")
	if err = equal("d", cfg.Map("d"), map[string]string{"e": "1", "f": "2"}); err != nil {
		t.Error(err)
		return
	}

	{
		t.Log("check virtual parameters")

		cfg.SetVirtual([]string{"virt"})
		cfg.Set("virt", "test")

		jsonBytes, err := cfg.SaveToData()
		if err != nil {
			t.Error(err)
			return
		}

		var res map[string]interface{}
		err = json.Unmarshal(jsonBytes, &res)
		if err != nil {
			t.Error(err)
			return
		}

		if _, exists := res["virt"]; exists {
			t.Error("virtual fields not working")
			return
		}
	}
}

func equal(key string, val interface{}, wait interface{}) error {
	if !reflect.DeepEqual(val, wait) {
		return fmt.Errorf("value of %s must be %s, now is %s",
			key, fmt.Sprintf("%+v", wait), fmt.Sprintf("%+v", val))
	}

	return nil
}
