package k8s

import (
	"fmt"
	"testing"
)

func TestServiceList(t *testing.T) {
	ss, err := ServiceList("xindaiquan")
	if err != nil {
		t.Error("ServiceList err", err)
		return
	}
	// b, _ := json.MarshalIndent(ss, "", "  ")
	// fmt.Println(string(b))
	for _, v := range ss {
		fmt.Println(v.GetMetadata().GetNamespace(), v.GetMetadata().GetName())
	}
}

func TestServiceListAll(t *testing.T) {
	ss, err := ServiceListAll()
	if err != nil {
		t.Error("ServiceListAll err", err)
		return
	}
	// b, _ := json.MarshalIndent(ss, "", "  ")
	// fmt.Println(string(b))
	for _, v := range ss {
		fmt.Println(v.GetMetadata().GetNamespace(), v.GetMetadata().GetName())
	}
}
