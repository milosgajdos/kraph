package dgraph

var Schema = `
	type Object {
		xid
		name
		kind
		ns
		created_at
	}

	xid: string @index(exact) .
	name: string @index(exact) .
	kind: string @index(exact) .
	ns: string @index(exact) .
	created_at : datetime @index(hour) .
	is_owned: [uid] @reverse .
`
