package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"os"
	"testing"
	"time"
)

/*
TestTransparencyVisualDemo is a test which provides a real-time interactive demonstration of all
transparency strategies. It initializes the terminal and animates the alpha value from 1.0 down to
0.0 for each algorithm, allowing you to see the dissolve effects in action.

Example:
    Expected Inputs:
        A terminal session with background and foreground images, animated across various transparency strategies.
    Expected Outputs:
        An animated terminal display showing smooth transitions between layers using CurtainWipe, Blinds, and Interlaced effects.
*/
func TestTransparencyVisualDemo(test *testing.T) {
	InitializeTerminal(80, 25)

	// Setup background and foreground images
	bgImage := "./test_data/image/complex_image.png"
	fgImage := "./test_data/image/complex_geometry.png"
	_ = LoadImage(bgImage)
	_ = LoadImage(fgImage)

	imageStyle := NewImageStyle()
	imageStyle.DrawingStyle = constants.ImageStyleBlockElementsAccurate
	imageStyle.DitheringStyle = constants.DitheringStyle4x4BayerMatrix

	// 1. Setup Background Layer (Complex Image)
	width, height := 80, 25
	bgLayer := AddLayer(0, 0, width, height, 1, nil)
	bgLayer.DrawImage(bgImage, imageStyle, 0, 0, width, height, 0)

	// 2. Setup Foreground Layer (Another Complex Image)
	fgLayer := AddLayer(0, 0, width, height, 2, nil)

	// 3. Define the strategies to test
	strategies := []struct {
		name     string
		strategy constants.TransparencyStrategy
	}{
		//{"Stochastic", constants.TransparencyStrategyStochastic},
		//{"Bayer2x2", constants.TransparencyStrategy2x2Bayer},
		//{"Bayer4x4", constants.TransparencyStrategy4x4Bayer},
		//{"Bayer8x8", constants.TransparencyStrategy8x8Bayer},
		{"Dissolve", constants.TransparencyStrategyDissolve},
	}

	for _, s := range strategies {
		fgLayer.SetTransparencyStrategy(s.strategy)

		// Draw foreground image for this strategy
		fgLayer.Clear()
		fgLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)

		// Set Alpha and Transition Progress
		fgLayer.SetAlphaValue(1.0)
		fgLayer.SetTransitionProgress(1.0)
		UpdateDisplay(false)

		// Animate Transition Progress from 1.0 to 0.0
		// We'll use a diagonal wipe for all rendering strategies to show how they combine.
		transitionStyle := NewTransitionStyle()
		transitionStyle.TransitionType = constants.TransitionTypeCurtainWipe
		transitionStyle.Direction = constants.TransitionDirectionBottomRightToTopLeft
		transitionStyle.SoftEdgeWidth = 0.2
		fgLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			fgLayer.SetTransitionProgress(currentProgress)
			fgLayer.Locate(0, 0)
			fgLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			fgLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: CurtainWipe | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)

			if currentProgress == 0 {
				result := commonResource.screenLayer.GetBasicAnsiString()
				_ = os.WriteFile(fmt.Sprintf("alpha_zero_%s.ansi", s.name), []byte(result), 0644)
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)

		// Demo Blinds (Horizontal, Forward: TopToBottom)
		transitionStyle.TransitionType = constants.TransitionTypeBlinds
		transitionStyle.BlindCount = 8
		transitionStyle.Direction = constants.TransitionDirectionTopToBottom
		fgLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			fgLayer.SetTransitionProgress(currentProgress)
			fgLayer.Locate(0, 0)
			fgLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			fgLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Horiz Blinds (Down) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)

		// Reset foreground for reverse test
		fgLayer.Clear()
		fgLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		fgLayer.SetAlphaValue(1.0)
		fgLayer.SetTransitionProgress(1.0)
		UpdateDisplay(false)

		// Demo Blinds (Horizontal, Reverse: BottomToTop)
		transitionStyle.Direction = constants.TransitionDirectionBottomToTop
		fgLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			fgLayer.SetTransitionProgress(currentProgress)
			fgLayer.Locate(0, 0)
			fgLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			fgLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Horiz Blinds (Up) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)

		// Demo Interlaced (Horizontal: Left and Right)
		fgLayer.Clear()
		fgLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		transitionStyle.TransitionType = constants.TransitionTypeInterlaced
		transitionStyle.Direction = constants.TransitionDirectionLeftToRight
		fgLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			fgLayer.SetTransitionProgress(currentProgress)
			fgLayer.Locate(0, 0)
			fgLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			fgLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Interlaced (H) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)

		// Demo Interlaced (Vertical: Top and Bottom)
		fgLayer.Clear()
		fgLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		transitionStyle.Direction = constants.TransitionDirectionTopToBottom
		fgLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			fgLayer.SetTransitionProgress(currentProgress)
			fgLayer.Locate(0, 0)
			fgLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			fgLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Interlaced (V) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(2000 * time.Millisecond)
	}

	DeleteAllLayers()
	commonResource.screen.Fini()
}
