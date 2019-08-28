package cloudprober

import (
	"testing"

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
