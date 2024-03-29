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

type server struct {
	server   string
	interval time.Duration
}

func newserver(grpcAddr string, i time.Duration) server {
	return server{server: grpcAddr, interval: i}
}

// see probe item: http://t.com:9313/status
func (s *server) sync() (err error) {
	client := cloudprober.New(s.server)
	for {
		log.Printf("sleep interval %v...\n", s.interval)
		time.Sleep(s.interval)
		log.Println("start new sync...")
		olditems, err := getOldTargets(client)
		if err != nil {
			log.Printf("get old targets err: %v\n", err)
			continue
		}
		targets, err := getServiceTargets(olditems)
		if err != nil {
			log.Printf("get service err: %v\n", err)
			continue
		}
		noanyupdate := true
		for _, v := range targets {
			if v.skip {
				continue
			}
			if v.delete {
				err := client.RemoveProbe(v.name)
				if err != nil {
					log.Printf("delete %v err: %v\n", v.name, err)
					continue
				}
				log.Printf("deleted %v ok\n", v.name)
				continue
			}
			if v.update {
				err := client.RemoveProbe(v.name)
				if err != nil {
					log.Printf("remove %v first err: %v\n", v.name, err)
					continue
				}
			}
			noanyupdate = false
			err := client.AddProbe(convert(v))
			if err != nil {
				log.Printf("addprobe %v err: %v\n", v.name, err)
				continue
			}
			log.Printf("addprobe %v ok\n", v.name)
		}
		if noanyupdate {
			log.Printf("there's no any update\n")
		}
		// time.Sleep(i)
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
		if v, ok := an[*pathAnnotation]; ok {
			path = v
		}
		// if strings.Contains(name, "prober") {
		// 	spew.Dump("an", an)
		// }
		if !*enableAll {
			var enable string
			if v, ok := an[*enableAnnotation]; ok {
				enable = v
			}
			if enable != "true" {
				// log.Printf("%v not enabled, skip", name)
				continue
			}
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

		url := fmt.Sprintf("%v:%v%v", ip, port, path)
		t := &item{
			name: name,
			url:  url,
		}
		ts[name] = t
	}
	return compare(olditems, ts), nil
}

func convert(v *item) *configpb.ProbeDef {
	interval := int32(*probeInterval)
	timeout := int32(*probeTimeout)
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
	if old == nil {
		return new
	}
	for k, v := range new {
		// exist, see if need update
		if oldv, ok := old[k]; ok {
			if oldv.url != v.url {
				v.update = true
				continue
			} else {
				v.skip = true
				// delete old item
				delete(old, k)
			}
		}
		// default is create
	}
	for k, v := range old {
		// what left need to delete
		v.delete = true
		new[k] = v
	}
	return new
}

type item struct {
	name   string
	url    string
	update bool
	skip   bool
	delete bool
}
