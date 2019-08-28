package main

import (
	"fmt"
	"log"
	"time"
	"wen/service-prober/pkg/cloudprober"
	"wen/service-prober/pkg/k8s"

	configpb "github.com/google/cloudprober/probes/proto"
	target "github.com/google/cloudprober/targets/proto"
)

// see probe item: http://t.com:9313/status
func sync(server string) (err error) {

	for {
		s := cloudprober.New(server)
		ts, err := gettargets()
		if err != nil {
			continue
		}
		for _, v := range ts {
			err := s.AddProbe(v)
			if err != nil {
				log.Printf("addprobe %v err: %v\n", v.GetName(), err)
				continue
			}
			log.Printf("addprobe %v ok\n", v.GetName())
		}
		time.Sleep(1 * time.Minute)
	}
	log.Printf("exit sync\n")
	return
}

func gettargets() (ts []*configpb.ProbeDef, err error) {
	keys := make(map[string]bool)

	ss, err := k8s.ServiceListAll()
	if err != nil {
		return
	}
	interval := int32(5000)
	timeout := int32(1000)
	ptype := configpb.ProbeDef_HTTP
	for _, v := range ss {
		name := fmt.Sprintf("%v/%v", v.GetMetadata().GetName(), v.GetMetadata().GetNamespace())
		if _, ok := keys[name]; ok {
			log.Printf("ignore %v, already exist\n", name)
			continue
		}
		ip := v.GetSpec().GetClusterIP()
		if ip == "None" {
			log.Printf("ignore %v, ip is None\n", name)
			continue
		}
		ports := v.GetSpec().GetPorts()
		if len(ports) == 0 {
			continue
		}
		port := ports[0].GetPort()
		url := fmt.Sprintf("http://%v:%v/healthz", ip, port)
		t := &configpb.ProbeDef{
			Name: &name,
			Targets: &target.TargetsDef{
				Type: &target.TargetsDef_HostNames{
					HostNames: url,
				},
			},
			Type:         &ptype,
			IntervalMsec: &interval,
			TimeoutMsec:  &timeout,
		}
		ts = append(ts, t)
		keys[name] = true
	}
	return
}
