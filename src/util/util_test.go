package util_test

import (
	"fmt"
	"lightOA-end/src/util"
	"testing"
)

func TestSha(t *testing.T) {
	fmt.Print(util.Sha256("123456"))
}
