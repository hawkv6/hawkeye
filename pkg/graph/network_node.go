package graph

type NetworkNode struct {
	id                 string
	name               string
	edges              map[string]Edge
	flexibleAlgorithms map[uint32]struct{}
}

func translateSrToFlexibleAlgorithm(srAlgorithms []uint32) map[uint32]struct{} {
	flexibleAlgorithms := make(map[uint32]struct{})
	for _, flexAlgoNumber := range srAlgorithms {
		if flexAlgoNumber != 0 && flexAlgoNumber != 1 {
			flexibleAlgorithms[flexAlgoNumber] = struct{}{}
		}
	}
	return flexibleAlgorithms
}

func NewNetworkNode(id, name string, srAlgorithms []uint32) *NetworkNode {
	return &NetworkNode{
		id:                 id,
		name:               name,
		flexibleAlgorithms: translateSrToFlexibleAlgorithm(srAlgorithms),
		edges:              make(map[string]Edge),
	}
}

func (node *NetworkNode) GetName() string {
	return node.name
}

func (node *NetworkNode) SetName(name string) {
	node.name = name
}

func (node *NetworkNode) GetId() string {
	return node.id
}

func (node *NetworkNode) GetEdges() map[string]Edge {
	return node.edges
}

func (node *NetworkNode) AddEdge(edge Edge) {
	node.edges[edge.GetId()] = edge
}

func (node *NetworkNode) DeleteEdge(id string) {
	delete(node.edges, id)
}

func (node *NetworkNode) SetFlexibleAlgorithms(srAlgorithms []uint32) {
	node.flexibleAlgorithms = translateSrToFlexibleAlgorithm(srAlgorithms)
	for _, edge := range node.edges {
		edge.UpdateFlexibleAlgorithms()
	}
}

func (node *NetworkNode) GetFlexibleAlgorithms() map[uint32]struct{} {
	return node.flexibleAlgorithms
}
