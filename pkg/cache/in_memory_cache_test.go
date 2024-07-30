package cache

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryCache(t *testing.T) {
	tests := []struct {
		name  string
		cache *InMemoryCache
	}{
		{
			name:  "Test NewInMemoryCache",
			cache: NewInMemoryCache(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.cache)
		})
	}
}

func TestInMemoryCache_Lock_Unlock(t *testing.T) {
	cache := NewInMemoryCache()
	var wg sync.WaitGroup

	expectedResults := make(map[string]string)
	var resultsMutex sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", val)
			value := fmt.Sprintf("val%d", val)

			cache.Lock()
			cache.prefixToRouterIdMap[key] = value
			cache.Unlock()

			// Store the expected result
			resultsMutex.Lock()
			expectedResults[key] = value
			resultsMutex.Unlock()
		}(i)
	}
	wg.Wait()

	for key, expectedValue := range expectedResults {
		cache.Lock()
		if got, exists := cache.prefixToRouterIdMap[key]; !exists || got != expectedValue {
			t.Errorf("For key %s, got %v, want %v", key, got, expectedValue)
		}
		cache.Unlock()
	}
}

func setUpDomainPrefix(key string, igpRouterId string, prefixValue string, prefixLength int32) *domain.DomainPrefix {
	prefix, _ := domain.NewDomainPrefix(&key, &igpRouterId, &prefixValue, &prefixLength)
	return prefix
}
func TestInMemoryCache_StoreClientNetwork(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		prefix domain.Prefix
	}
	tests := []struct {
		name string
		args args
		want domain.Prefix
	}{
		{
			name: "Test StoreClientNetwork - Case 1",
			args: args{
				prefix: setUpDomainPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
			},
			want: setUpDomainPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreClientNetwork(tt.args.prefix)
			clientNetwork := cache.prefixStore[tt.args.prefix.GetKey()]
			if !reflect.DeepEqual(clientNetwork, tt.want) {
				t.Errorf("Got %v, want %v", clientNetwork, tt.want)
			}
			igpRouterId := cache.prefixToRouterIdMap[tt.args.prefix.GetPrefix()]
			if !reflect.DeepEqual(igpRouterId, tt.args.prefix.GetIgpRouterId()) {
				t.Errorf("Got %v, want %v", igpRouterId, tt.args.prefix.GetIgpRouterId())
			}
		})
	}
}

func TestInMemoryCache_RemoveClientNetwork(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		prefix domain.Prefix
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test RemoveClientNetwork - Case 1",
			args: args{
				prefix: setUpDomainPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
			},
		},
		{
			name: "Test RemoveClientNetwork - Case 2",
			args: args{
				prefix: setUpDomainPrefix("2_0_2_0_0_2001:db8:c::_64_0000.0000.000c", "0000.0000.000c", "2001:db8:c::", 64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreClientNetwork(tt.args.prefix)
			cache.RemoveClientNetwork(tt.args.prefix)
			if _, exists := cache.prefixStore[tt.args.prefix.GetKey()]; exists {
				t.Errorf("Prefix %s was not removed from the cache", tt.args.prefix.GetKey())
			}
			if _, exists := cache.prefixToRouterIdMap[tt.args.prefix.GetPrefix()]; exists {
				t.Errorf("Prefix %s was not removed from the prefixToRouterIdMap", tt.args.prefix.GetPrefix())
			}
		})
	}
}

func TestInMemoryCache_GetClientNetworkByKey(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want domain.Prefix
	}{
		{
			name: "Test GetClientNetworkByKey success",
			args: args{
				key: "2_0_2_0_0_2001:db8:b::_64_0000.0000.000b",
			},
			want: setUpDomainPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreClientNetwork(tt.want)
			clientNetwork := cache.GetClientNetworkByKey(tt.args.key)
			if !reflect.DeepEqual(clientNetwork, tt.want) {
				t.Errorf("Got %v, want %v", clientNetwork, tt.want)
			}
		})
	}
}

func setUpDomainSid(key string, igpRouterId string, sidValue string, algorithm uint32) *domain.DomainSid {
	sid, _ := domain.NewDomainSid(&key, &igpRouterId, &sidValue, &algorithm)
	return sid
}

func TestInMemoryCache_StoreSid(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		sid domain.Sid
	}
	tests := []struct {
		name string
		args args
		want domain.Sid
	}{
		{
			name: "Test StoreSid - Case 1",
			args: args{
				sid: setUpDomainSid("0_0000.0000.0007_fc00:0:7:0:1::", "0000.0000.0007", "fc00:0:7:0:1::", 0),
			},
			want: setUpDomainSid("0_0000.0000.0007_fc00:0:7:0:1::", "0000.0000.0007", "fc00:0:7:0:1::", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreSid(tt.args.sid)
			sid := cache.sidStore[tt.args.sid.GetKey()]
			if !reflect.DeepEqual(sid, tt.want) {
				t.Errorf("Got %v, want %v", sid, tt.want)
			}
			igpRouterId := cache.igpRouterIdToSrAlgoToSidMap[tt.args.sid.GetIgpRouterId()][tt.args.sid.GetAlgorithm()]
			if !reflect.DeepEqual(igpRouterId, tt.args.sid.GetKey()) {
				t.Errorf("Got %v, want %v", igpRouterId, tt.args.sid.GetKey())
			}
		})
	}
}
func TestInMemoryCache_RemoveSid(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		sid domain.Sid
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test RemoveSid - Case 1",
			args: args{
				sid: setUpDomainSid("0_0000.0000.0007_fc00:0:7:0:1::", "0000.0000.0007", "fc00:0:7:0:1::", 0),
			},
		},
		{
			name: "Test RemoveSid - Case 2",
			args: args{
				sid: setUpDomainSid("1_0000.0000.0008_fc00:0:8:0:1::", "0000.0000.0008", "fc00:0:8:0:1::", 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreSid(tt.args.sid)
			cache.RemoveSid(tt.args.sid)
			if _, exists := cache.sidStore[tt.args.sid.GetKey()]; exists {
				t.Errorf("Sid %s was not removed from the cache", tt.args.sid.GetKey())
			}
			if _, exists := cache.igpRouterIdToSrAlgoToSidMap[tt.args.sid.GetIgpRouterId()][tt.args.sid.GetAlgorithm()]; exists {
				t.Errorf("Sid %s was not removed from the igpRouterIdToSrAlgoToSidMap", tt.args.sid.GetKey())
			}
		})
	}
}
func TestInMemoryCache_GetSidByKey(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want domain.Sid
	}{
		{
			name: "Test GetSidByKey - Case 1",
			args: args{
				key: "0_0000.0000.0007_fc00:0:7:0:1::",
			},
			want: setUpDomainSid("0_0000.0000.0007_fc00:0:7:0:1::", "0000.0000.0007", "fc00:0:7:0:1::", 0),
		},
		{
			name: "Test GetSidByKey - Case 2",
			args: args{
				key: "1_0000.0000.0008_fc00:0:8:0:1::",
			},
			want: setUpDomainSid("1_0000.0000.0008_fc00:0:8:0:1::", "0000.0000.0008", "fc00:0:8:0:1::", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreSid(tt.want)
			sid := cache.GetSidByKey(tt.args.key)
			if !reflect.DeepEqual(sid, tt.want) {
				t.Errorf("Got %v, want %v", sid, tt.want)
			}
		})
	}
}

func TestInMemoryCache_GetRouterIdFromNetworkAddress(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		prefix domain.Prefix
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GetRouterIdFromNetworkAddress",
			args: args{
				prefix: setUpDomainPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
			},
			want: "0000.0000.000b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreClientNetwork(tt.args.prefix)
			routerId := cache.GetRouterIdFromNetworkAddress(tt.args.prefix.GetPrefix())
			if !reflect.DeepEqual(routerId, tt.want) {
				t.Errorf("Got %v, want %v", routerId, tt.want)
			}
		})
	}
}

func TestInMemoryCache_GetSrAlgorithmSid(t *testing.T) {
	type args struct {
		sid         domain.Sid
		routerIgpId string
		algorithm   uint32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GetSrAlgorithmSid successfully",
			args: args{
				sid:         setUpDomainSid("0_0000.0000.0007_fc00:0:7:0:1::", "0000.0000.0007", "fc00:0:7:0:1::", 0),
				routerIgpId: "0000.0000.0007",
				algorithm:   0,
			},
			want: "fc00:0:7:0:1::",
		},
		{
			name: "Test GetSrAlgorithmSid routerId not in cache",
			args: args{
				sid:         setUpDomainSid("1_0000.0000.0008_fc00:0:8:0:1::", "0000.0000.0008", "fc00:0:8:0:1::", 1),
				routerIgpId: "0000.0000.0007",
				algorithm:   0,
			},
			want: "",
		},
		{
			name: "Test GetSrAlgorithmSid wrong algorithm",
			args: args{
				sid:         setUpDomainSid("1_0000.0000.0008_fc00:0:8:0:1::", "0000.0000.0008", "fc00:0:8:0:1::", 1),
				routerIgpId: "0000.0000.0008",
				algorithm:   0,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewInMemoryCache()
			cache.StoreSid(tt.args.sid)
			sid := cache.GetSrAlgorithmSid(tt.args.routerIgpId, tt.args.algorithm)
			if !reflect.DeepEqual(sid, tt.want) {
				t.Errorf("Got %v, want %v", sid, tt.want)
			}
		})
	}
}

func setUpDomainNode(key string, igpRouterId string, name string, srAlgorithm []uint32) *domain.DomainNode {
	node, _ := domain.NewDomainNode(&key, &igpRouterId, &name, srAlgorithm)
	return node
}

func TestInMemoryCache_StoreNode(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		node domain.Node
	}
	tests := []struct {
		name string
		args args
		want domain.Node
	}{
		{
			name: "Test StoreNode ",
			args: args{
				node: setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
			},
			want: setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreNode(tt.args.node)
			storedNode := cache.nodeStore[tt.args.node.GetKey()]
			if !reflect.DeepEqual(storedNode, tt.want) {
				t.Errorf("Got %v, want %v", storedNode, tt.want)
			}
			nodeKey := cache.igpRouterIdToRouterKeyMap[tt.args.node.GetIgpRouterId()]
			if !reflect.DeepEqual(nodeKey, tt.args.node.GetKey()) {
				t.Errorf("Got %v, want %v", nodeKey, tt.args.node.GetKey())
			}
		})
	}
}
func TestInMemoryCache_RemoveNode(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		node domain.Node
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test RemoveNode - Case 1",
			args: args{
				node: setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
			},
		},
		{
			name: "Test RemoveNode - Case 2",
			args: args{
				node: setUpDomainNode("2_0_0_0000.0000.0005", "0000.0000.0005", "XR-5", []uint32{0, 1}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreNode(tt.args.node)
			cache.RemoveNode(tt.args.node)
			if _, exists := cache.nodeStore[tt.args.node.GetKey()]; exists {
				t.Errorf("Node %s was not removed from the cache", tt.args.node.GetKey())
			}
			if _, exists := cache.igpRouterIdToRouterKeyMap[tt.args.node.GetIgpRouterId()]; exists {
				t.Errorf("Node %s was not removed from the igpRouterIdToRouterKeyMap", tt.args.node.GetKey())
			}
		})
	}
}
func TestInMemoryCache_GetNodeByKey(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want domain.Node
	}{
		{
			name: "Test GetNodeByKey - Case 1",
			args: args{
				key: "2_0_0_0000.0000.0004",
			},
			want: setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
		},
		{
			name: "Test GetNodeByKey - Case 2",
			args: args{
				key: "2_0_0_0000.0000.0005",
			},
			want: setUpDomainNode("2_0_0_0000.0000.0005", "0000.0000.0005", "XR-5", []uint32{0, 1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreNode(tt.want)
			node := cache.GetNodeByKey(tt.args.key)
			if !reflect.DeepEqual(node, tt.want) {
				t.Errorf("Got %v, want %v", node, tt.want)
			}
		})
	}
}
func TestInMemoryCache_GetNodeByIgpRouterId(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		igpRouterId string
		node        domain.Node
	}
	tests := []struct {
		name string
		args args
		want domain.Node
	}{
		{
			name: "Test GetNodeByIgpRouterId - Case 1",
			args: args{
				igpRouterId: "0000.0000.0004",
				node:        setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
			},
			want: setUpDomainNode("2_0_0_0000.0000.0004", "0000.0000.0004", "XR-4", []uint32{0, 1}),
		},
		{
			name: "Test GetNodeByIgpRouterId - Case 2",
			args: args{
				igpRouterId: "0000.0000.0005",
				node:        setUpDomainNode("2_0_0_0000.0000.0005", "0000.0000.0005", "XR-5", []uint32{0, 1}),
			},
			want: setUpDomainNode("2_0_0_0000.0000.0005", "0000.0000.0005", "XR-5", []uint32{0, 1}),
		},
		{
			name: "Test GetNodeByIgpRouterId - Case 2",
			args: args{
				igpRouterId: "0000.0000.0006",
				node:        setUpDomainNode("2_0_0_0000.0000.0005", "0000.0000.0005", "XR-5", []uint32{0, 1}),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreNode(tt.args.node)
			node := cache.GetNodeByIgpRouterId(tt.args.igpRouterId)
			if !reflect.DeepEqual(node, tt.want) {
				t.Errorf("Got %v, want %v", node, tt.want)
			}
		})
	}
}
func TestInMemoryCache_StoreServiceSid(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		serviceType      string
		servicePrefixSid string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test StoreServiceSid - fw",
			args: args{
				serviceType:      "fw",
				servicePrefixSid: "fc00:0:2f::",
			},
			want: []string{"fc00:0:2f::"},
		},
		{
			name: "Test StoreServiceSid - ids",
			args: args{
				serviceType:      "ids",
				servicePrefixSid: "fc00:0:6f::",
			},
			want: []string{"fc00:0:6f::"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreServiceSid(tt.args.serviceType, tt.args.servicePrefixSid)
			serviceSids := cache.GetServiceSids(tt.args.serviceType)
			if !reflect.DeepEqual(serviceSids, tt.want) {
				t.Errorf("Got %v, want %v", serviceSids, tt.want)
			}
		})
	}
}
func TestInMemoryCache_RemoveServiceSid(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		serviceType      string
		servicePrefixSid string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test RemoveServiceSid - Case 1",
			args: args{
				serviceType:      "fw",
				servicePrefixSid: "fc00:0:2f::",
			},
		},
		{
			name: "Test RemoveServiceSid - Case 2",
			args: args{
				serviceType:      "ids",
				servicePrefixSid: "fc00:0:6f::",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreServiceSid(tt.args.serviceType, tt.args.servicePrefixSid)
			cache.RemoveServiceSid(tt.args.serviceType, tt.args.servicePrefixSid)
			serviceSids := cache.GetServiceSids(tt.args.serviceType)
			if len(serviceSids) != 0 {
				t.Errorf("Service SIDs still exist in cache after removal")
			}
		})
	}
}
func TestInMemoryCache_GetServiceSids(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		serviceType string
		serviceSid  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test GetServiceSids - add fw service",
			args: args{
				serviceType: "fw",
				serviceSid:  "fc00:0:2f::",
			},
			want: []string{"fc00:0:2f::"},
		},
		{
			name: "Test GetServiceSids - add ids service",
			args: args{
				serviceType: "ids",
				serviceSid:  "fc00:0:6f::",
			},
			want: []string{"fc00:0:6f::"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreServiceSid(tt.args.serviceType, tt.args.serviceSid)
			serviceSids := cache.GetServiceSids(tt.args.serviceType)
			if !reflect.DeepEqual(serviceSids, tt.want) {
				t.Errorf("Got %v, want %v", serviceSids, tt.want)
			}
		})
	}
}
func TestInMemoryCache_DoesServiceSidExist(t *testing.T) {
	cache := NewInMemoryCache()
	type args struct {
		servicePrefixSid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test DoesServiceSidExist - Case 1",
			args: args{
				servicePrefixSid: "fc00:0:2f::",
			},
			want: true,
		},
		{
			name: "Test DoesServiceSidExist - Case 2",
			args: args{
				servicePrefixSid: "fc00:0:6f::",
			},
			want: true,
		},
		{
			name: "Test DoesServiceSidExist - Case 3",
			args: args{
				servicePrefixSid: "fc00:0:af::",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.StoreServiceSid("fw", "fc00:0:2f::")
			cache.StoreServiceSid("ids", "fc00:0:6f::")
			exists := cache.DoesServiceSidExist(tt.args.servicePrefixSid)
			if exists != tt.want {
				t.Errorf("Got %v, want %v", exists, tt.want)
			}
		})
	}
}
