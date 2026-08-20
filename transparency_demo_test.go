package consolizer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/supercom32/consolizer/constants"
)

/*
TestTransparencyVisualDemo is a test which provides a demonstration of all transparency strategies. It initializes the
terminal in headless simulation mode and animates the alpha value from 1.0 down to 0.0 for each algorithm, rendering
every frame through the full compositing pipeline without requiring a real terminal display. In addition, the
following should be noted:

  - The test runs headless by enabling debug mode before terminal initialization, exactly like the other tests in this
    suite. To watch the demo interactively on a real terminal, remove the debug mode line and run the test on its own.

  - Frame delays are only applied when a real terminal screen exists, so the headless run completes as fast as the
    rendering pipeline allows.

Example:

	Expected Inputs:
	    A terminal session with background and foreground images, animated across various transparency strategies.
	Expected Outputs:
	    Every frame of the CurtainWipe, Blinds, and Interlaced transitions renders through the compositing pipeline,
	    and the final alpha zero screen is written to an ansi file for inspection.
*/
func TestTransparencyVisualDemo(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	demoSleep := func(duration time.Duration) {
		if commonResource.screen != nil {
			time.Sleep(duration)
		}
	}

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
	forgroundLayer := AddLayer(0, 0, width, height, 2, nil)

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
		forgroundLayer.SetTransparencyStrategy(s.strategy)

		// Draw foreground image for this strategy
		forgroundLayer.Clear()
		forgroundLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)

		// Set Alpha and Transition Progress
		forgroundLayer.SetAlphaValue(1.0)
		forgroundLayer.SetTransitionProgress(1.0)
		UpdateDisplay(false)

		// Animate Transition Progress from 1.0 to 0.0
		// We'll use a diagonal wipe for all rendering strategies to show how they combine.
		transitionStyle := NewTransitionStyle()
		transitionStyle.TransitionType = constants.TransitionTypeCurtainWipe
		transitionStyle.Direction = constants.TransitionDirectionBottomRightToTopLeft
		transitionStyle.SoftEdgeWidth = 0.2
		forgroundLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			forgroundLayer.SetTransitionProgress(currentProgress)
			forgroundLayer.Locate(0, 0)
			forgroundLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			forgroundLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: CurtainWipe | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)

			if currentProgress == 0 {
				result := commonResource.screenLayer.GetBasicAnsiString()
				_ = os.WriteFile(fmt.Sprintf("alpha_zero_%s.ansi", s.name), []byte(result), 0644)
			}
			demoSleep(100 * time.Millisecond)
		}
		demoSleep(1000 * time.Millisecond)

		// Demo Blinds (Horizontal, Forward: TopToBottom)
		transitionStyle.TransitionType = constants.TransitionTypeBlinds
		transitionStyle.BlindCount = 8
		transitionStyle.Direction = constants.TransitionDirectionTopToBottom
		forgroundLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			forgroundLayer.SetTransitionProgress(currentProgress)
			forgroundLayer.Locate(0, 0)
			forgroundLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			forgroundLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Horiz Blinds (Down) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			demoSleep(100 * time.Millisecond)
		}
		demoSleep(1000 * time.Millisecond)

		// Reset foreground for reverse test
		forgroundLayer.Clear()
		forgroundLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		forgroundLayer.SetAlphaValue(1.0)
		forgroundLayer.SetTransitionProgress(1.0)
		UpdateDisplay(false)

		// Demo Blinds (Horizontal, Reverse: BottomToTop)
		transitionStyle.Direction = constants.TransitionDirectionBottomToTop
		forgroundLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			forgroundLayer.SetTransitionProgress(currentProgress)
			forgroundLayer.Locate(0, 0)
			forgroundLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			forgroundLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Horiz Blinds (Up) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			demoSleep(100 * time.Millisecond)
		}
		demoSleep(1000 * time.Millisecond)

		// Demo Interlaced (Horizontal: Left and Right)
		forgroundLayer.Clear()
		forgroundLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		transitionStyle.TransitionType = constants.TransitionTypeInterlaced
		transitionStyle.Direction = constants.TransitionDirectionLeftToRight
		forgroundLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			forgroundLayer.SetTransitionProgress(currentProgress)
			forgroundLayer.Locate(0, 0)
			forgroundLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			forgroundLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Interlaced (H) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			demoSleep(100 * time.Millisecond)
		}
		demoSleep(1000 * time.Millisecond)

		// Demo Interlaced (Vertical: Top and Bottom)
		forgroundLayer.Clear()
		forgroundLayer.DrawImage(fgImage, imageStyle, 0, 0, width, height, 0)
		transitionStyle.Direction = constants.TransitionDirectionTopToBottom
		forgroundLayer.SetTransitionStyle(transitionStyle)

		for a := float32(1.0); a >= -0.001; a -= 0.05 {
			currentProgress := a
			if currentProgress < 0 {
				currentProgress = 0
			}
			forgroundLayer.SetTransitionProgress(currentProgress)
			forgroundLayer.Locate(0, 0)
			forgroundLayer.Color24Bit(constants.AnsiColorByIndex[constants.ColorRed], constants.AnsiColorByIndex[constants.ColorWhite])
			forgroundLayer.Print(fmt.Sprintf(" Strategy: %-12s | Transition: Interlaced (V) | Progress: %0.2f ", s.name, currentProgress))
			UpdateDisplay(false)
			demoSleep(100 * time.Millisecond)
		}
		demoSleep(2000 * time.Millisecond)
	}

	DeleteAllLayers()
	if commonResource.screen != nil {
		commonResource.screen.Fini()
	}
}
