package runny

import "github.com/fatih/color"

var defaultShell = "/bin/bash"
var primaryColor *color.Color = color.New(color.Bold)
var secondaryColor *color.Color = color.New(color.FgHiBlack)
var errorColor *color.Color = color.New(color.FgRed)
