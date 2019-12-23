package typrails_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
)

func TestModule(t *testing.T) {
	m := typrails.Module()
	require.True(t, typcore.IsBuildCommander(m))
}