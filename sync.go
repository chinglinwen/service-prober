package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"wen/service-prober/pkg/cloudprober"
	"wen/service-prober/pkg/k8s"

	// "github.com/davecgh/go-spew/spew"
	// "github.com/getlantern/deepcopy"

	"github.com/davecgh/go-spew/spew"
	configpb "github.com/google/cloudprober/probes/proto"
	target "github.com/google/cloudprober/targets/proto"
)

const (
	pathAnnotation   = "prober.haodai.net/path"
	enableAnnotation = "prober.haodai.net/enable"
)

// see probe item: http://t.com:9313/status
func sync(server string, i time.Duration) (err error) {

	for {
		log.Println("start new sync...")
		s := cloudprober.New(server)
		olditems, err := getOldTargets(s)
		if err != nil {
			log.Printf("get old targets err: %v\n", err)
			continue
		}
		ts, err := getServiceTargets(olditems)
		if err != nil {
			log.Printf("get service err: %v\n", err)
			continue
		}
		noanyupdate := true
		for _, v := range ts {
			if v.skip {
				continue
			}
			if v.update {
				err := s.RemoveProbe(v.name)
				if err != nil {
					log.Printf("remove %v first err: %v\n", v.name, err)
					continue
				}
			}
			noanyupdate = false
			err := s.AddProbe(convert(v))
			if err != nil {
				log.Printf("addprobe %v err: %v\n", v.name, err)
				continue
			}
			log.Printf("addprobe %v ok\n", v.name)
		}
		if noanyupdate {
			log.Printf("there's no any update\n")
		}
		time.Sleep(i)
	}
	log.Printf("exit sync\n")
	return
}

// var olditems map[string]*item

func getOldTargets(s *cloudprober.Client) (ts map[string]*item, err error) {
	ps, err := s.ListProbe()
	if err != nil {
		return
	}
	log.Printf("got %v exist probes", len(ps))
	ts = make(map[string]*item)
	for _, v := range ps {
		name := v.GetConfig().GetName()
		url := v.GetConfig().GetTargets().GetHostNames()
		if name == "" {
			continue
		}
		// log.Printf("add %v,%v to olditems", name, url)
		t := &item{
			name: name,
			url:  url,
		}
		ts[name] = t
	}
	return
}

func getServiceTargets(olditems map[string]*item) (ts map[string]*item, err error) {
	ts = make(map[string]*item)
	// keys := make(map[string]bool)s
	ss, err := k8s.ServiceListAll()
	if err != nil {
		return
	}
	for _, v := range ss {
		name := fmt.Sprintf("%v/%v", v.GetMetadata().GetNamespace(), v.GetMetadata().GetName())
		if _, ok := ts[name]; ok {
			log.Printf("ignore %v, already exist\n", name)
			continue
		}

		var path string
		an := v.GetMetadata().GetAnnotations()
		if v, ok := an[pathAnnotation]; ok {
			path = v
		}
		if strings.Contains(name, "prober") {
			spew.Dump("an", an)
		}
		var enable string
		if v, ok := an[enableAnnotation]; ok {
			enable = v
		}
		if enable != "true" {
			// log.Printf("%v not enabled, skip", name)
			continue
		}

		ip := v.GetSpec().GetClusterIP()
		if ip == "None" {
			// log.Printf("ignore %v, ip is None\n", name)
			continue
		}
		ports := v.GetSpec().GetPorts()
		if len(ports) == 0 {
			continue
		}
		port := ports[0].GetPort()

		url := fmt.Sprintf("http://%v:%v%v", ip, port, path)
		t := &item{
			name: name,
			url:  url,
		}
		ts[name] = t
	}
	return compare(olditems, ts), nil
}

func convert(v *item) *configpb.ProbeDef {
	interval := int32(5000)
	timeout := int32(1000)
	ptype := configpb.ProbeDef_HTTP
	return &configpb.ProbeDef{
		Name: &v.name,
		Targets: &target.TargetsDef{
			Type: &target.TargetsDef_HostNames{
				HostNames: v.url,
			},
		},
		Type:         &ptype,
		IntervalMsec: &interval,
		TimeoutMsec:  &timeout,
	}

}
func compare(old, new map[string]*item) map[string]*item {
	// spew.Dump("old", old)
	// spew.Dump("new", new)
	if old == nil {
		// olditems = new
		// deepcopy.Copy(olditems, new)
		return new
	}
	// i := 0
	for k, v := range new {
		// if i == 0 {
		// 	spew.Dump(k, v)
		// 	i++
		// }
		if oldv, ok := old[k]; ok {
			if oldv.url != v.url {
				v.update = true
				continue
			} else {
				v.skip = true
			}
		}
	}
	// olditems = new
	// deepcopy.Copy(olditems, new)
	return new
}

type item struct {
	name   string
	url    string
	update bool
	skip   bool
}
