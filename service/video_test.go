package service

import (
	"fmt"
	"testing"
	"time"
)

func TestActionService(t *testing.T) {
	format := time.Now().Format("20060102150405")
	fmt.Println(format)
}
