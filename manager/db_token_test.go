package manager

import (
	"testing"
)

func TestTokenGenerater(t *testing.T) {
	tokenInit()
	TokenNew(1)
	TokenNew(2)
	t.Log(TokenList())
}
