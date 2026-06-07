package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFTSQueryUsesPrefixTerms(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `"lutan"*`, buildFTSQuery("lutan"))
	assert.Equal(t, `"mem"* AND "muscu"*`, buildFTSQuery("mem muscu"))
	assert.Empty(t, buildFTSQuery("   "))
}
