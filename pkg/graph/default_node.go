package graph

type DefaultNode struct {
	id    interface{}
	edges map[interface{}]Edge
}

func NewDefaultNode(id interface{}) *DefaultNode {
	return &DefaultNode{
		id:    id,
		edges: make(map[interface{}]Edge)}
}

func (node *DefaultNode) GetId() interface{} {
	return node.id
}

func (node *DefaultNode) GetEdges() map[interface{}]Edge {
	return node.edges
}

func (node *DefaultNode) AddEdge(edge Edge) {
	node.edges[edge.GetId()] = edge
}
