//go:build !darwin

package printing

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func printWindow(ctx context.Context) error {
	runtime.WindowPrint(ctx)
	return nil
}
