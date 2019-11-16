package typdocker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

func TestModule(t *testing.T) {
	m := typdocker.Module()
	require.True(t, typcli.IsBuildCommander(m))
}
