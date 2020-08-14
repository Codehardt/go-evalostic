package evalostic

func parseCondition(s string) (Node, error) {
	t, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	return parse(t)
}
