package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestConfigTypes(t *testing.T) {
	str := "{\"a\":123, \"b\":\"тест\", \"c\":[\"one\", \"two\"], \"d\":{\"e\":\"1\", \"f\": \"2\"}}"

	t.Log("init config")
	cfg := NewConfig()
	err := cfg.InitFromData([]byte(str))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("get string parameter")
	if err = equal("a string", cfg.String("a"), "123"); err != nil {
		t.Error(err)
		return
	}

	cfg.Set("aa", "a123")
	if err = equal("aa string", cfg.String("aa"), "a123"); err != nil {
		t.Error(err)
		return
	}

	t.Log("get cyrillic string parameter")
	if err = equal("b string", cfg.String("b"), "тест"); err != nil {
		t.Error(err)
		return
	}

	cfg.Set("ba", "текст")
	if err = equal("ba string", cfg.String("ba"), "текст"); err != nil {
		t.Error(err)
		return
	}

	t.Log("get slice parameter")
	if err = equal("c slice", cfg.Slice("c"), []string{"one", "two"}); err != nil {
		t.Error(err)
		return
	}

	cfg.Set("ca", []string{"1", "2", "3"})
	if err = equal("ca slice", cfg.Slice("ca"), []string{"1", "2", "3"}); err != nil {
		t.Error(err)
		return
	}

	t.Log("get map parameter")
	if err = equal("d map", cfg.Map("d"), map[string]string{"e": "1", "f": "2"}); err != nil {
		t.Error(err)
		return
	}

	cfg.Set("da", map[string]string{"key": "value"})
	if err = equal("da map", cfg.Map("da"), map[string]string{"key": "value"}); err != nil {
		t.Error(err)
		return
	}

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

func TestChange(t *testing.T) {
	cfg := NewConfig()
	err := cfg.InitFromData([]byte("{}"))
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.HasChange() {
		t.Error("there were no changes")
	}

	cfg.Set("a", "b")

	if !cfg.HasChange() {
		t.Error("there have been changes")
	}

	cfg.hasChange = false

	cfg.Set("a", "b")

	if cfg.HasChange() {
		t.Error("there were no changes")
	}
}

func equal(msgKey string, val interface{}, wait interface{}) error {
	if !reflect.DeepEqual(val, wait) {
		return fmt.Errorf("value of %s must be %s, now is %s",
			msgKey, fmt.Sprintf("%+v", wait), fmt.Sprintf("%+v", val))
	}

	return nil
}
