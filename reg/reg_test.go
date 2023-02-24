package reg

import "testing"

func TestGetNamedMap(t *testing.T) {
	pattern := `^\/user\/(?P<id>[0-9]{1,})?\/career\/(?P<position>.*$)`
	target := `/user/123456/career/officer`
	res := GetNamedMap(pattern, target)
	if res != nil && res["id"] == "123456" && res["position"] == "officer" {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

func BenchmarkGetNamedMap(b *testing.B) {
	pattern := `^\/user\/(?P<id>[0-9]{1,})?\/career\/(?P<position>.*$)`
	target := `/user/123456/career/officer`
	for i := 0; i < b.N; i++ {
		_ = GetNamedMap(pattern, target)
	}
}
