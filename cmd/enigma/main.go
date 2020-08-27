package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mkideal/cli"
	enigma "github.com/thomas-chastaingt/Enigmatic"
)

type CLIOpts struct {
	Help      bool `cli:"!h,help" usage:"Show help."`
	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

	Rotors    []string `cli:"rotors" name:"I II III" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     []int    `cli:"rings" name:"1 1 1" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
	Position  []string `cli:"position" name:"A A A" usage:"Starting position of the rotors: from A (default) to Z for each."`
	Plugboard []string `cli:"plugboard" name:"[]" usage:"Optional plugboard pairs to scramble the message further."`

	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-Thin, C-Thin."`
}

var CLIDefaults = struct {
	Reflector string
	Ring      int
	Position  string
	Rotors    []string
}{
	Reflector: "B",
	Ring:      1,
	Position:  "A",
	Rotors:    []string{"I", "II", "III"},
}

func SetDefaults(argv *CLIOpts) {
	if argv.Reflector == "" {
		argv.Reflector = CLIDefaults.Reflector
	}
	if len(argv.Rotors) == 0 {
		argv.Rotors = CLIDefaults.Rotors
	}
	loadRings := (len(argv.Rings) == 0)
	loadPosition := (len(argv.Position) == 0)
	if loadRings || loadPosition {
		for range argv.Rotors {
			if loadRings {
				argv.Rings = append(argv.Rings, CLIDefaults.Ring)
			}
			if loadPosition {
				argv.Position = append(argv.Position, CLIDefaults.Position)
			}
		}
	}
}

func main() {

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		originalPlaintext := strings.Join(ctx.Args(), " ")
		plaintext := enigma.SanitizePlaintext(originalPlaintext)

		if argv.Help || len(plaintext) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}

		config := make([]enigma.RotorConfig, len(argv.Rotors))
		for index, rotor := range argv.Rotors {
			ring := argv.Rings[index]
			value := argv.Position[index][0]
			config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
		}

		e := enigma.NewEnigma(config, argv.Reflector, argv.Plugboard)
		encoded := e.EncodeString(plaintext)

		if argv.Condensed {
			fmt.Print(encoded)
			return nil
		}

		tmpl, _ := template.New("cli").Parse(OutputTemplate)
		err := tmpl.Execute(os.Stdout, struct {
			Original, Plain, Encoded string
			Args                     *CLIOpts
			Ctx                      *cli.Context
		}{originalPlaintext, plaintext, encoded, argv, ctx})
		return err

	})

}
