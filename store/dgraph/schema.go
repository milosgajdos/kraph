package dgraph

// Schema is dgraph schema
var Schema = `
	type Object {
		xid: string
		name: string
		kind: string
		namespace: string
		created_at: datetime
		link: Object
	}

	xid: string @index(exact) .
	name: string @index(exact) .
	kind: string @index(exact) .
	namespace: string @index(exact) .
	created_at : datetime @index(hour) .
	link: [uid] @reverse .
`
