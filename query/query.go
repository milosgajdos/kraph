package query

type Query struct {
	matchers map[string]*matcher
}

func Build() *Query {
	q := &Query{
		matchers: make(map[string]*matcher),
	}

	return q.MatchAny()
}

func (q *Query) updateQuery(prop string, val interface{}, fn MatchFunc) *Query {
	q.matchers[prop] = newMatcher(val, fn)
	return q
}

func (q *Query) UID(uid interface{}, fn MatchFunc) *Query {
	return q.updateQuery("uid", uid, fn)
}

func (q *Query) Namespace(ns interface{}, fn MatchFunc) *Query {
	return q.updateQuery("ns", ns, fn)
}

func (q *Query) Kind(k interface{}, fn MatchFunc) *Query {
	return q.updateQuery("kind", k, fn)
}

func (q *Query) Name(n interface{}, fn MatchFunc) *Query {
	return q.updateQuery("name", n, fn)
}

func (q *Query) Version(v interface{}, fn MatchFunc) *Query {
	return q.updateQuery("version", v, fn)
}

func (q *Query) Group(g interface{}, fn MatchFunc) *Query {
	return q.updateQuery("group", g, fn)
}

func (q *Query) Entity(e interface{}, fn MatchFunc) *Query {
	return q.updateQuery("entity", e, fn)
}

func (q *Query) Weight(w interface{}, fn MatchFunc) *Query {
	return q.updateQuery("weight", w, fn)
}

func (q *Query) Attrs(a interface{}, fn MatchFunc) *Query {
	return q.updateQuery("attrs", a, fn)
}

func (q *Query) Metadata(m interface{}, fn MatchFunc) *Query {
	return q.updateQuery("metadata", m, fn)
}

func (q *Query) Matcher() *match {
	return &match{
		q: q,
	}
}

func (q *Query) MatchAny() *Query {
	for _, s := range []string{
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
		q.matchers[s] = newMatcher(MatchAny, AnyFunc)
	}

	return q
}
