package main

import (
	"context"
	pb "github.com/evilsocket/sum/proto"
	. "github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBalancing(t *testing.T) {
	newTimeout := time.Minute
	timeout = &newTimeout
	pollPeriod = &newTimeout

	dir1, err := ioutil.TempDir("", "")
	Nil(t, err)
	defer os.RemoveAll(dir1)
	dir2, err := ioutil.TempDir("", "")
	Nil(t, err)
	defer os.RemoveAll(dir2)

	for _, baseDir := range []string{dir1, dir2} {
		for _, childDir := range []string{"data", "oracles"} {
			err = os.Mkdir(filepath.Join(baseDir, childDir), 0755)
			Nil(t, err)
		}
	}

	node1, sum1 := spawnNode(t, 12345, dir1)
	defer node1.Stop()
	node2, sum2 := spawnNode(t, 12346, dir2)
	defer node2.Stop()

	testRecord := &pb.Record{Data: []float32{0.1, 0.2, 0.3}, Meta: map[string]string{"name": "test"}}

	for i := 1; i <= 100; i++ {
		rCopy := &pb.Record{Id: uint64(i), Data: testRecord.Data, Meta: testRecord.Meta}
		resp, err := sum1.CreateRecordWithId(context.Background(), rCopy)
		Nil(t, err)
		True(t, resp.Success)
	}

	for i := 101; i <= 120; i++ {
		rCopy := &pb.Record{Id: uint64(i), Data: testRecord.Data, Meta: testRecord.Meta}
		resp, err := sum2.CreateRecordWithId(context.Background(), rCopy)
		Nil(t, err)
		True(t, resp.Success)
	}

	// now we have 100 records on node1, 20 on node2

	master, _ := spawnOrchestrator(t, 12347, "localhost:12345,localhost:12346")
	defer master.Stop()

	Equal(t, sum1.NumRecords(), sum2.NumRecords())
}
