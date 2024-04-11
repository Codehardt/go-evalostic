package evalostic

import (
	"fmt"
	"testing"
)

func TestElasticSearchQueryExport(t *testing.T) {
	for _, condition := range []string{
		`"foo"`,
		`"foo" OR "bar"`,
		`"foo" AND "bar"`,
		`"foo" AND NOT "bar"`,
		`"foo" OR ("bar" AND NOT "baz")`,
		`"foo" OR NOT ("bar" AND NOT "baz")`,
		`("foo" OR "bar") AND NOT ("bar" AND ("baz" OR "qux"))`,
	} {
		ev, err := New([]string{condition})
		if err != nil {
			t.Fatalf("could not parse %s: %s", condition, err)
		}
		query := ev.ExportElasticSearchQuery("raw", false)
		t.Logf("%s: %s", condition, query)
		query = ev.ExportElasticSearchQuery("raw", true)
		t.Logf("%s: %s", condition, query)
	}
}

func TestElasticSearchQueryExport2(t *testing.T) {
	for _, condition := range []string{
		`"A" AND NOT ("B" OR "C" OR "D")`,
	} {
		fmt.Println(condition)
		fmt.Println("-----")
		ev, err := New([]string{condition})
		if err != nil {
			t.Fatalf("could not parse %s: %s", condition, err)
		}
		query := ev.ExportElasticSearchQuery("raw", false)
		t.Logf("%s: %s", condition, query)
	}
}
