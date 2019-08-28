package main

import (
	"fmt"
	"log"
	"testing"
	"wen/service-prober/pkg/cloudprober"

	"github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
)

func TestUpdate(t *testing.T) {
	// var old map[string]*item

	// new := map[string]*item{
	// 	"a": &item{
	// 		name: "a",
	// 		url:  "url1",
	// 	}}

	// _ = compare(olditems, new)
	// // new["a"].url = "url2"
	// new1 := map[string]*item{
	// 	"a": &item{
	// 		name: "a",
	// 		url:  "url2",
	// 	}}

	s := cloudprober.New(*server)
	olditems, err := getOldTargets(s)
	if err != nil {
		log.Printf("get old targets err: %v\n", err)
		return
	}
	for k, v := range olditems {
		fmt.Printf("%v: %v\n", k, v)
	}
	ts, err := getServiceTargets(olditems)
	if err != nil {
		log.Printf("get service err: %v\n", err)
		return
	}

	spew.Dump(olditems["netshoot/t"], ts["netshoot/t"])

	// a := compare(olditems, new1)
	// if !a["a"].update {
	// 	t.Error("update should be true")
	// }
}
