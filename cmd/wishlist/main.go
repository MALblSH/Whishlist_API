package main

import (
	"log/slog"
	"os"

	"github.com/MALblSH/Wishlist_API/internal/application/runner"
)

func main() {
	err := runner.RunWishlist()
	if err != nil {
		slog.Error("wishlist application error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
