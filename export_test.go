package evalostic

import "testing"

func TestElasticSearchQueryExport(t *testing.T) {
	for _, condition := range []string{
		`"foo"`,
		`"foo" OR "bar"`,
		`"foo" AND "bar"`,
		`"foo" AND NOT "bar"`,
		`("foo" OR "bar") AND NOT ("bar" AND ("baz" OR "qux"))`,
	} {
		ev, err := New([]string{condition})
		if err != nil {
			t.Fatalf("could not parse %s: %s", condition, err)
		}
		query := ev.ExportElasticSearchQuery("raw")
		t.Logf("%s: %s", condition, query)
	}
}
