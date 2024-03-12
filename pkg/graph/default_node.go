package graph

type DefaultNode struct {
	id    interface{}
	name  string
	edges map[interface{}]Edge
}

func NewDefaultNode(id interface{}) *DefaultNode {
	return &DefaultNode{
		id:    id,
		edges: make(map[interface{}]Edge)}
}

func (node *DefaultNode) GetName() string {
	return node.name
}

func (node *DefaultNode) SetName(name string) {
	node.name = name
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

func (node *DefaultNode) DeleteEdge(id interface{}) {
	delete(node.edges, id)
}
