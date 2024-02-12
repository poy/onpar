package matchers

type DifferUser interface {
	UseDiffer(Differ)
}

type ChainMatcher struct {
	Children []Matcher
	differ   Differ
}

func Chain(a, b Matcher, ms ...Matcher) *ChainMatcher {
	return &ChainMatcher{
		Children: append(append([]Matcher{a}, b), ms...),
	}
}

func (m *ChainMatcher) UseDiffer(d Differ) {
	m.differ = d
}

func (m ChainMatcher) Match(actual any) (any, error) {
	var err error
	next := actual
	for _, child := range m.Children {
		if d, ok := child.(DifferUser); ok && m.differ != nil {
			d.UseDiffer(m.differ)
		}
		next, err = child.Match(next)
		if err != nil {
			return nil, err
		}
	}
	return next, nil
}
