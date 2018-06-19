package stats_test

import (
	"testing"
	"github.com/msales/pkg/stats"
	"github.com/magiconair/properties/assert"
)

func TestTags_With(t *testing.T) {
	t1 := stats.Tags{"k": "v"}

	t2 := t1.With("k2", "v2")

	assert.Equal(t, t2, stats.Tags{"k": "v", "k2": "v2"})
}

func TestTags_With_OverrideKey(t *testing.T) {
	t1 := stats.Tags{"k": "original"}

	t2 := t1.With("k", "overridden")

	assert.Equal(t, t2, stats.Tags{"k": "overridden"})
}

func TestTags_Merge(t *testing.T) {
	t1 := stats.Tags{"k": "v"}
	t2 := stats.Tags{"k2": "v2"}

	t3 := t1.Merge(t2)

	assert.Equal(t, t3, stats.Tags{"k": "v", "k2": "v2"})
}

func TestTags_Merge_OverrideKey(t *testing.T) {
	t1 := stats.Tags{"k": "original"}
	t2 := stats.Tags{"k2": "v2", "k": "overridden"}

	t3 := t1.Merge(t2)

	assert.Equal(t, t3, stats.Tags{"k": "overridden", "k2": "v2"})
}
