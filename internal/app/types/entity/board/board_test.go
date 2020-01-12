package board

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	a := &Entity{
		Subjects: []Subject{
			{
				ThreadKey:    "0123456789",
				ThreadTitle:  "たいとる",
				MessageCount: 123,
				LastModified: time.Now(),
			},
		},
		WriteCount: 100,
	}

	s := fmt.Sprintf("%v", a)
	if strings.HasPrefix(s, "0x") { // 0x40c138
		t.Errorf("fmt.Sprintf(\"v\",a) = %v", s)
	}
}
