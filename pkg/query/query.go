package query

type Entity int

const (
	Node Entity = iota
	Edge
)

type Query struct {
	matchers map[string]*matcher
}

func Build() *Query {
	q := &Query{
		matchers: make(map[string]*matcher),
	}

	return q.MatchAny()
}

func (q *Query) updateQuery(prop string, val interface{}, funcs ...MatchFunc) *Query {
	q.matchers[prop] = newMatcher(val, funcs...)
	return q
}

func (q *Query) UID(uid interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("uid", uid, funcs...)
}

func (q *Query) Name(n interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("name", n, funcs...)
}

func (q *Query) Group(g interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("group", g, funcs...)
}

func (q *Query) Version(v interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("version", v, funcs...)
}

func (q *Query) Kind(k interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("kind", k, funcs...)
}

func (q *Query) Namespace(ns interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("ns", ns, funcs...)
}

func (q *Query) Entity(e interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("entity", e, funcs...)
}

func (q *Query) Weight(w interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("weight", w, funcs...)
}

func (q *Query) Attrs(a interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("attrs", a, funcs...)
}

func (q *Query) Metadata(m interface{}, funcs ...MatchFunc) *Query {
	return q.updateQuery("metadata", m, funcs...)
}

func (q *Query) Matcher() *match {
	return &match{
		q: q,
	}
}

func (q *Query) Reset() *Query {
	return Build()
}

func (q *Query) MatchAny() *Query {
	for _, s := range []string{
		"uid",
		"ns",
		"kind",
		"name",
		"version",
		"group",
		"entity",
		"weight",
		"attrs",
		"metadata",
	} {
		q.matchers[s] = newMatcher(MatchAny, IsAnyFunc)
	}

	return q
}
