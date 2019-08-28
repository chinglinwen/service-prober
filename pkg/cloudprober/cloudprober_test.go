package cloudprober

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	configpb "github.com/google/cloudprober/probes/proto"
	proto1 "github.com/google/cloudprober/targets/proto"
)

var testc = New("localhost:9314")

/*
	/home/wen/gocode/src/wen/service-prober/pkg/cloudprober/cloudprober_test.go:23: add err add probe err: rpc error:
	code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection closed
*/
func TestAdd(t *testing.T) {
	name := "test"
	interval := int32(5000)
	timeout := int32(1000)
	ptype := configpb.ProbeDef_HTTP
	err := testc.AddProbe(&configpb.ProbeDef{
		Name: &name,
		Targets: &proto1.TargetsDef{
			Type: &proto1.TargetsDef_HostNames{
				HostNames: "baidu.com",
			},
		},
		Type:         &ptype,
		IntervalMsec: &interval,
		TimeoutMsec:  &timeout,
	})
	if err != nil {
		t.Errorf("add err: %v", err)
		return
	}
}
func TestList(t *testing.T) {
	ts, err := testc.ListProbe()
	if err != nil {
		t.Errorf("list err: %v", err)
		return
	}
	fmt.Printf("len %v\n", len(ts))
	if len(ts) > 0 {
		// fmt.Printf("%#v\n", ts[0])
		spew.Dump(ts[0])
	}

	v := ts[0]

	fmt.Printf("ok %#v,%#v\n,", v.GetName(), v.GetConfig().GetTargets().GetHostNames())

}
