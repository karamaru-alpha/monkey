package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode_String(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	assert.Equal(t, expected, concatted.String())
}

func TestCode_Make(t *testing.T) {
	for _, tt := range []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	} {
		assert.Equal(t, tt.expected, Make(tt.op, tt.operands...))
	}
}

func TestCode_ReadOperands(t *testing.T) {
	for _, tt := range []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	} {
		instructions := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		assert.NoError(t, err)

		operandsRead, n := ReadOperands(def, instructions[1:])
		assert.Equal(t, tt.bytesRead, n)
		assert.Equal(t, tt.operands, operandsRead)
	}
}
