package authz

type namedRelation struct {
	name string
}

func (r *namedRelation) String() string {
	return r.name
}

var _ Relation = &namedRelation{}

var RelationList = &namedRelation{name: "list"}
var RelationCreate = &namedRelation{name: "create"}
var RelationRetrieve = &namedRelation{name: "retrieve"}
var RelationUpdate = &namedRelation{name: "update"}
var RelationDelete = &namedRelation{name: "delete"}
