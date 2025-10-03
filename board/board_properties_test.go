package board

import (
	"testing"
	"testing/quick"
	"zenmojo/sharing"
)

// TestBoardPropertiesQuick uses testing/quick to verify board properties
func TestBoardPropertiesQuick(t *testing.T) {
	// Define a function that generates a board and checks its properties
	f := func() bool {
		board := New()
		props := checkBoardProperties(board.Grid())

		// If any property is false, log it for debugging
		if !props.HasNoSingleStones || !props.HasValidGroupSizes ||
			!props.UsesValidColors || !props.IsFull {
			// Get the share code for the failing board
			shareCode, err := sharing.Encode(board.Grid())
			if err != nil {
				t.Logf("Error generating share code: %v", err)
			}

			t.Logf("Found invalid board with properties: %+v\nBoard share code: %s", props, shareCode)
			return false
		}
		return true
	}

	// Configure quick.Check
	config := &quick.Config{
		MaxCount: 1000, // Number of iterations
	}

	// Run the quick check
	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}
