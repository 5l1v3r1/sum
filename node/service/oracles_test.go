package service

import (
	"context"
	"github.com/evilsocket/sum/node/storage"
	"testing"

	pb "github.com/evilsocket/sum/proto"
)

var (
	testOracles  = 5
	byName       = pb.ByName{Name: "findReasonsToLive"}
	brokenOracle = pb.Oracle{
		Id:   123,
		Name: "brokenOracle",
		Code: "lulz i won't compile =)",
	}
	updatedOracle = pb.Oracle{
		Id:   666,
		Name: "myNameHasBeenUpdated",
		Code: "function myBodyToo(){ return 0; }",
	}
)

func sameOracle(a, b pb.Oracle) bool {
	return a.Id == b.Id && a.Name == b.Name && a.Code == b.Code
}

func TestServiceErrOracleResponse(t *testing.T) {
	if r := errOracleResponse("test %d", 123); r.Success {
		t.Fatal("success should be false")
	} else if r.Msg != "test 123" {
		t.Fatalf("unexpected message: %s", r.Msg)
	} else if r.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", r.Oracle)
	}
}

func TestServiceCreateOracle(t *testing.T) {
	setupFolders(t)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.CreateOracle(context.TODO(), &testOracle); err != nil {
		t.Fatal(err)
	} else if !resp.Success {
		t.Fatalf("expected success response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	} else if resp.Msg != "1" {
		t.Fatalf("unexpected response message: %s", resp.Msg)
	}
}

func TestServiceCreateOracleWithInvalidId(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	svc, err := New(testFolder, "", "")
	if err != nil {
		t.Fatal(err)
	}

	svc.oracles.NextID(1)
	if resp, err := svc.CreateOracle(context.TODO(), &testOracle); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Msg != storage.ErrInvalidID.Error() {
		t.Fatalf("unexpected response message: %s", resp.Msg)
	}
}

func TestServiceCreateBrokenOracle(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.CreateOracle(context.TODO(), &brokenOracle); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	}
}

func TestServiceUpdateOracle(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	updatedOracle.Id = 1

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.UpdateOracle(context.TODO(), &updatedOracle); err != nil {
		t.Fatal(err)
	} else if !resp.Success {
		t.Fatalf("expected success response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	} else if stored := svc.oracles.Find(updatedOracle.Id); stored == nil {
		t.Fatal("expected stored oracle with id 1")
	} else if !sameOracle(*stored, updatedOracle) {
		t.Fatal("oracle has not been updated as expected")
	}
}

func TestServiceUpdateOracleWithInvalidId(t *testing.T) {
	setupFolders(t)
	defer teardown(t)

	updatedOracle.Id = 1
	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.UpdateOracle(context.TODO(), &updatedOracle); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Msg != storage.ErrRecordNotFound.Error() {
		t.Fatalf("unexpected response message: %s", resp.Msg)
	}
}

func TestServiceUpdateWithBrokenOracle(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	brokenOracle.Id = 1

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.UpdateOracle(context.TODO(), &brokenOracle); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	}
}

func TestServiceReadOracle(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	byID.Id = 1
	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.ReadOracle(context.TODO(), &byID); err != nil {
		t.Fatal(err)
	} else if !resp.Success {
		t.Fatalf("expected success response: %v", resp)
	} else if resp.Oracle == nil {
		t.Fatal("expected oracles list")
	} else if testOracle.Id = byID.Id; !sameOracle(*resp.Oracle, testOracle) {
		t.Fatalf("oracle does not match: %v", resp.Oracle)
	}
}

func TestServiceReadOracleWithInvalidId(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.ReadOracle(context.TODO(), &pb.ById{Id: 666}); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	} else if resp.Msg != "oracle 666 not found." {
		t.Fatalf("unexpected message: %s", resp.Msg)
	}
}

func TestServiceFindOracle(t *testing.T) {
	bak := testOracles
	testOracles = 1
	defer func() {
		testOracles = bak
	}()

	setup(t, true, true)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.FindOracle(context.TODO(), &byName); err != nil {
		t.Fatal(err)
	} else if !resp.Success {
		t.Fatalf("expected success response: %v", resp)
	} else if resp.Oracle == nil {
		t.Fatal("expected oracles list")
	} else if testOracle.Id = byID.Id; !sameOracle(*resp.Oracle, testOracle) {
		t.Fatalf("oracle does not match: %v", resp.Oracle)
	}
}

func TestServiceFindOracleWithInvalidName(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.FindOracle(context.TODO(), &pb.ByName{Name: "no way i'm an oracle name :D"}); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected success response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("expected oracle: %v", resp.Oracle)
	}
}

func TestServiceListOraclesSinglePage(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	list := pb.ListRequest{
		Page:    1,
		PerPage: uint64(testOracles),
	}

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.ListOracles(context.TODO(), &list); err != nil {
		t.Fatal(err)
	} else if resp.Total != uint64(testOracles) {
		t.Fatalf("expected %d total oracles, got %d", testOracles, resp.Total)
	} else if resp.Pages != 1 {
		t.Fatalf("expected 3 pages got %d", resp.Pages)
	} else if len(resp.Oracles) != testOracles {
		t.Fatalf("expected %d total oracles, got %d", testOracles, len(resp.Oracles))
	} else {
		for _, r := range resp.Oracles {
			if testOracle.Id = r.Id; !sameOracle(*r, testOracle) {
				t.Fatalf("unexpected oracle: %v", r)
			}
		}
	}
}

func TestServiceListOraclesMultiPage(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	list := pb.ListRequest{
		Page:    1,
		PerPage: 2,
	}

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.ListOracles(context.TODO(), &list); err != nil {
		t.Fatal(err)
	} else if resp.Total != uint64(testOracles) {
		t.Fatalf("expected %d total oracles, got %d", testOracles, resp.Total)
	} else if resp.Pages != 3 {
		t.Fatalf("expected 3 pages got %d", resp.Pages)
	} else if len(resp.Oracles) != 2 {
		t.Fatalf("expected %d total oracles, got %d", 2, len(resp.Oracles))
	} else {
		for _, r := range resp.Oracles {
			if testOracle.Id = r.Id; !sameOracle(*r, testOracle) {
				t.Fatalf("unexpected oracle: %v", r)
			}
		}
	}
}

func TestServiceListOraclesInvalidPage(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	list := pb.ListRequest{
		Page:    100000,
		PerPage: 2,
	}

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.ListOracles(context.TODO(), &list); err != nil {
		t.Fatal(err)
	} else if resp.Total != uint64(testOracles) {
		t.Fatalf("expected %d total oracles, got %d", testOracles, resp.Total)
	} else if resp.Pages != 3 {
		t.Fatalf("expected 3 pages got %d", resp.Pages)
	} else if len(resp.Oracles) != 0 {
		t.Fatalf("expected %d total oracles, got %d", 0, len(resp.Oracles))
	}
}

func TestServiceDeleteOracle(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	svc, err := New(testFolder, "", "")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < testOracles; i++ {
		id := uint64(i + 1)
		if resp, err := svc.DeleteOracle(context.TODO(), &pb.ById{Id: id}); err != nil {
			t.Fatal(err)
		} else if !resp.Success {
			t.Fatalf("expected success response: %v", resp)
		} else if resp.Oracle != nil {
			t.Fatalf("unexpected oracle: %v", resp.Oracle)
		} else if svc.NumOracles() != testOracles-int(id) {
			t.Fatalf("inconsistent oracles storage size of %d", svc.NumOracles())
		}
	}

	if svc.NumOracles() != 0 {
		t.Fatalf("expected empty oracles storage, found %d instead", svc.NumOracles())
	} else if doublecheck, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if doublecheck.NumOracles() != 0 {
		t.Fatalf("%d dat files left on disk", doublecheck.NumOracles())
	}
}

func TestServiceDeleteOracleWithInvalidId(t *testing.T) {
	setup(t, true, true)
	defer teardown(t)

	if svc, err := New(testFolder, "", ""); err != nil {
		t.Fatal(err)
	} else if resp, err := svc.DeleteOracle(context.TODO(), &pb.ById{Id: 666}); err != nil {
		t.Fatal(err)
	} else if resp.Success {
		t.Fatalf("expected error response: %v", resp)
	} else if resp.Oracle != nil {
		t.Fatalf("unexpected oracle: %v", resp.Oracle)
	} else if resp.Msg != "Oracle 666 not found." {
		t.Fatalf("unexpected message: %s", resp.Msg)
	}
}
