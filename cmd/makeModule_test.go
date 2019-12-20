package cmd

import "testing"

func TestMakeModuleContent(t *testing.T) {
	result := MakeModuleContent("RedisModule")
	t.Log(result)
}

func TestMakeModuleProviderContent(t *testing.T) {
	result := MakeModuleProviderContent("RedisModule")
	t.Log(result)
}

func TestMakeModule(t *testing.T) {
	err := MakeModule("RedisModule", "")
	if err != nil {
		t.Fatal(err)
	}
}
