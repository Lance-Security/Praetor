/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/

package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	useColour  = true
	useBold    = true
	envNoColor = color.NoColor
)

// ConfigureTerminal applies terminal formatting preferences.
// If the environment already disables colour (e.g., NO_COLOR), we respect that.
func ConfigureTerminal(useColourCfg bool, useBoldCfg bool) {
	useColour = useColourCfg
	useBold = useBoldCfg
	color.NoColor = envNoColor || !useColourCfg
}

// Palette

func Default(a ...any) string { return fmt.Sprint(a...) }
func Defaultf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func Primary(a ...any) string { return sprintStyled(primaryAttrs(), a...) }
func Primaryf(format string, a ...any) string {
	return sprintStyledf(primaryAttrs(), format, a...)
}

func Muted(a ...any) string { return sprintStyled(mutedAttrs(), a...) }
func Mutedf(format string, a ...any) string {
	return sprintStyledf(mutedAttrs(), format, a...)
}

func Accept(a ...any) string { return sprintStyled(acceptAttrs(), a...) }
func Acceptf(format string, a ...any) string {
	return sprintStyledf(acceptAttrs(), format, a...)
}

func Warning(a ...any) string { return sprintStyled(warningAttrs(), a...) }
func Warningf(format string, a ...any) string {
	return sprintStyledf(warningAttrs(), format, a...)
}

func Error(a ...any) string { return sprintStyled(errorAttrs(), a...) }
func Errorf(format string, a ...any) string {
	return sprintStyledf(errorAttrs(), format, a...)
}

func primaryAttrs() []color.Attribute {
	attrs := []color.Attribute{color.FgCyan}
	if useBold {
		attrs = append(attrs, color.Bold)
	}
	if !useColour {
		attrs = nil
	}
	return attrs
}

func mutedAttrs() []color.Attribute {
	if useColour {
		return []color.Attribute{color.FgHiBlack}
	}
	return nil
}

func acceptAttrs() []color.Attribute {
	attrs := []color.Attribute{color.FgGreen}
	if useBold {
		attrs = append(attrs, color.Bold)
	}
	if !useColour {
		attrs = nil
	}
	return attrs
}

func warningAttrs() []color.Attribute {
	attrs := []color.Attribute{color.FgYellow}
	if useBold {
		attrs = append(attrs, color.Bold)
	}
	if !useColour {
		attrs = nil
	}
	return attrs
}

func errorAttrs() []color.Attribute {
	attrs := []color.Attribute{color.FgRed}
	if useBold {
		attrs = append(attrs, color.Bold)
	}
	if !useColour {
		attrs = nil
	}
	return attrs
}

func sprintStyled(attrs []color.Attribute, a ...any) string {
	if len(attrs) == 0 {
		return fmt.Sprint(a...)
	}
	if useColour {
		return color.New(attrs...).Sprint(a...)
	}
	return fmt.Sprint(a...)
}

func sprintStyledf(attrs []color.Attribute, format string, a ...any) string {
	if len(attrs) == 0 {
		return fmt.Sprintf(format, a...)
	}
	if useColour {
		return color.New(attrs...).Sprintf(format, a...)
	}
	return fmt.Sprintf(format, a...)
}
