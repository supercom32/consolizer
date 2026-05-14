package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/math"
)

/*
GetTransitionedColor is a method which allows you to obtain a color that has been transitioned to another color by a
specific percent. For example, if your source color is red (255, 0, 0) and your target color is green (0, 255, 0),
transitioning by 0.5 (fifty percent) will yield the color (128, 128, 0). In addition, the following should be noted:

- If your percent change yields color indexes which are not evenly divisible, then the color index will be rounded up.

- If you pass in a percent change of less than 0.0 or greater than 1.0, you are simply specifying that you want to.

- If the resultant transitioned color falls outside of the RGB range of Black (0, 0, 0) or White (255, 255, 255), it.

:param sourceColor: The color you wish to transition from.
:param targetColor: The color you wish to transition to.
:param percentChange: The percentage of the transition to apply (0.0 to 1.0).

:return: The transitioned color as a ColorType.

Example:

	transitionedColor := GetTransitionedColor(GetRGBColor(255, 0, 0), GetRGBColor(0, 255, 0), 0.5)
*/
func GetTransitionedColor(sourceColor constants.ColorType, targetColor constants.ColorType, percentChange float32) constants.ColorType {
	var sourceColorIndex [3]int32
	var targetColorIndex [3]int32
	var newColorIndex [3]int32
	sourceColorIndex[0], sourceColorIndex[1], sourceColorIndex[2] = GetRGBColorComponents(sourceColor)
	targetColorIndex[0], targetColorIndex[1], targetColorIndex[2] = GetRGBColorComponents(targetColor)
	for currentColorIndex := 0; currentColorIndex < 3; currentColorIndex++ {
		colorDifference := targetColorIndex[currentColorIndex] - sourceColorIndex[currentColorIndex]
		colorDifference = int32(math.RoundToWholeNumber(float32(colorDifference) * percentChange))
		if colorDifference < 0 {
			colorDifference = int32(math.GetAbsoluteValueAsFloat64(colorDifference))
			newColorIndex[currentColorIndex] = sourceColorIndex[currentColorIndex] - colorDifference
		} else {
			newColorIndex[currentColorIndex] = sourceColorIndex[currentColorIndex] + colorDifference
		}
		if newColorIndex[currentColorIndex] > 255 {
			newColorIndex[currentColorIndex] = 255
		}
		if newColorIndex[currentColorIndex] < 0 {
			newColorIndex[currentColorIndex] = 0
		}
	}
	return constants.ColorType(tcell.NewRGBColor(newColorIndex[0], newColorIndex[1], newColorIndex[2]))
}

/*
GetRGBColorComponents is a method which allows you to obtain RGB color component indexes for red, green, an blue color
channels.

:param color: The color for which you want to retrieve the RGB components.

:return: The red, green, and blue color component indexes (0-255).

Example:

	red, green, blue := GetRGBColorComponents(constants.ColorRed)
*/
func GetRGBColorComponents(color constants.ColorType) (int32, int32, int32) {
	var redColorIndex int32
	var greenColorIndex int32
	var blueColorIndex int32
	redColorIndex, greenColorIndex, blueColorIndex = tcell.Color.RGB(tcell.Color(color))
	return redColorIndex, greenColorIndex, blueColorIndex
}
