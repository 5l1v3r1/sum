package orchestrator

import (
	"fmt"
	. "github.com/evilsocket/sum/proto"
	"github.com/evilsocket/islazy/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
)

// information to work with a specific node
type NodeInfo struct {
	sync.RWMutex
	ID uint
	// Name of this node ( it's address )
	Name string
	// Certificate file ( TODO: optionally load it in base64 format )
	CertFile string
	// GRPC client to the node's sum service
	Client SumServiceClient
	// GRPC client to the node's sum internal service
	InternalClient SumInternalServiceClient
	// node's status
	status ServerInfo
	// records stored on this node
	RecordIds map[uint64]bool
}

// update node's status
func (n *NodeInfo) UpdateStatus() {
	ctx, cf := newCommContext()
	defer cf()
	srvInfo, err := n.Client.Info(ctx, &Empty{})

	if err != nil {
		log.Error("Unable to update node '%s' status: %v", n.Name, err)
		return
	}

	n.Lock()
	defer n.Unlock()

	n.status = *srvInfo
}

// get currently available node's status
func (n *NodeInfo) Status() ServerInfo {
	n.RLock()
	defer n.RUnlock()
	return n.status
}

// create node from a connection string and a certificate file
// this method connect to the node and create it's respective GRPC clients.
// it verify the connection by retrieving node's status using the aforementioned clients.
func CreateNode(node, certFile string) (*NodeInfo, error) {
	var dialOptions grpc.DialOption

	if certFile != "" {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return nil, fmt.Errorf("cannot load certificate file '%s': %v", certFile, err)
		}
		dialOptions = grpc.WithTransportCredentials(creds)
	} else {
		dialOptions = grpc.WithInsecure()
	}
	conn, err := grpc.Dial(node, dialOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to dial service at '%s': %v", node, err)
	}
	client := NewSumServiceClient(conn)
	internalClient := NewSumInternalServiceClient(conn)
	ctx, cancelFn := newCommContext()
	defer cancelFn()

	// check service availability
	svcInfo, err := client.Info(ctx, &Empty{})
	if err != nil {
		return nil, fmt.Errorf("unable to get service info from node '%s': %v", node, err)
	}

	ni := &NodeInfo{
		status:         *svcInfo,
		Name:           node,
		CertFile:       certFile,
		Client:         client,
		InternalClient: internalClient,
		RecordIds:      make(map[uint64]bool),
	}

	// get stored records
	pages := int(svcInfo.Records / 1024)
	if svcInfo.Records%1024 > 0 {
		pages++
	}
	for i := 0; i < pages; i++ {
		resp, err := client.ListRecords(ctx, &ListRequest{Page: uint64(i + 1), PerPage: 1024})
		if err != nil {
			return nil, fmt.Errorf("unable to get list of records from node '%s': %v", node, err)
		}

		for _, r := range resp.Records {
			ni.RecordIds[r.Id] = true
		}
	}

	return ni, nil
}