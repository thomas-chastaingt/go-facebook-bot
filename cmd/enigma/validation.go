package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mkideal/cli"
	enigma "github.com/thomas-chastaingt/Enigmatic"
)

func (argv *CLIOpts) Validate(ctx *cli.Context) error {
	SetDefaults(argv)
	validators := [](func(argv *CLIOpts, ctx *cli.Context) error){
		ValidatePlugboard,
		ValidateRotors,
		ValidateReflector,
		ValidatePosition,
		ValidateRings,
		ValidateUniformity,
	}
	for _, validator := range validators {
		if err := validator(argv, ctx); err != nil {
			return err
		}
	}
	return nil
}

func ValidatePlugboard(argv *CLIOpts, ctx *cli.Context) error {
	var plugboard string
	for _, pair := range argv.Plugboard {
		if matched, _ := regexp.MatchString(`^[A-Z]{2}$`, pair); !matched {
			return fmt.Errorf(
				`plugboard should be grouped by letter pairs ("AB CD"), got "%s"`,
				ctx.Color().Yellow(pair))
		}
		if strings.ContainsAny(pair, plugboard) || pair[0] == pair[1] {
			return fmt.Errorf(
				`letters cannot repeat across the plugboard, check "%s"`,
				ctx.Color().Yellow(pair))
		}
		plugboard += pair
	}
	return nil
}

func ValidateRotors(argv *CLIOpts, ctx *cli.Context) error {
	for _, rotor := range argv.Rotors {
		if r := enigma.HistoricRotors.GetByID(rotor); r == nil {
			return fmt.Errorf(`unknown rotor "%s"`, ctx.Color().Yellow(rotor))
		}
	}
	return nil
}

func ValidateReflector(argv *CLIOpts, ctx *cli.Context) error {
	if r := enigma.HistoricReflectors.GetByID(argv.Reflector); r == nil {
		return fmt.Errorf(`unknown reflector "%s"`, ctx.Color().Yellow(argv.Reflector))
	}
	return nil
}

func ValidatePosition(argv *CLIOpts, ctx *cli.Context) error {
	for _, char := range argv.Position {
		if matched, _ := regexp.MatchString(`^[A-Z]$`, char); !matched {
			return fmt.Errorf(
				`rotor positions should be single letters in the A-Z range, got "%s"`,
				ctx.Color().Yellow(char))
		}
	}
	return nil
}

func ValidateRings(argv *CLIOpts, ctx *cli.Context) error {
	for _, ring := range argv.Rings {
		if ring < 1 || ring > 26 {
			return fmt.Errorf(
				`ring out of range: must be 1-26, got "%s"`,
				ctx.Color().Yellow(ring))
		}
	}
	return nil
}

func ValidateUniformity(argv *CLIOpts, ctx *cli.Context) error {
	if !(len(argv.Rotors) == len(argv.Position) && len(argv.Position) == len(argv.Rings)) {
		return fmt.Errorf(
			"number of configured rotors, rings, and position settings should be equal")
	}
	return nil
}
