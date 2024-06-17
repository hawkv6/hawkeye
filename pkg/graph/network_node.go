package graph

type NetworkNode struct {
	id    string
	name  string
	edges map[string]Edge
}

func NewNetworkNode(id string) *NetworkNode {
	return &NetworkNode{
		id:    id,
		edges: make(map[string]Edge),
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
