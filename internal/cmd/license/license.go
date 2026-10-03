package license

import (
	"fmt"

	"github.com/alecthomas/kong"
	"go.vnbr.de/thumbctl"
)

type LicenseCmd struct {
	Project    bool `negatable:"" default:"1" help:"Show the project license."`
	ThirdParty bool `negatable:"" default:"0" aliases:"3p" help:"Show the third party licenses."`
}

func (cmd *LicenseCmd) Run(app *kong.Kong) error {
	if cmd.Project {
		fmt.Println(thumbctl.ProjectLicense)
	}
	if cmd.ThirdParty {
		fmt.Println(thumbctl.ThirdPartyLicenses)
	}
	return nil
}
