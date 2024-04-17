package graph

type NetworkNode struct {
	id    interface{}
	name  string
	edges map[interface{}]Edge
}

func NewNetworkNode(id interface{}) *NetworkNode {
	return &NetworkNode{
		id:    id,
		edges: make(map[interface{}]Edge),
	}
}

func (node *NetworkNode) GetName() string {
	return node.name
}

func (node *NetworkNode) SetName(name string) {
	node.name = name
}

func (node *NetworkNode) GetId() interface{} {
	return node.id
}

func (node *NetworkNode) GetEdges() map[interface{}]Edge {
	return node.edges
}

func (node *NetworkNode) AddEdge(edge Edge) {
	node.edges[edge.GetId()] = edge
}

func (node *NetworkNode) DeleteEdge(id interface{}) {
	delete(node.edges, id)
}
