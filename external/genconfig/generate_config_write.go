package genconfig

import (
	"fmt"
	"io"
	"os"

	"hashicorp/terraform/external/plans"
	"hashicorp/terraform/external/tfdiags"
)

func ValidateTargetFile(out string) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Target generated file already exists",
			"Terraform can only write generated config into a new file. Either choose a different target location or move all existing configuration out of the target file, delete it and try again."))

	}
	return diags
}

func MaybeWriteGeneratedConfig(plan *plans.Plan, out string) tfdiags.Diagnostics {
	if len(out) == 0 {
		// No specified out file, so don't write anything.
		return nil
	}

	diags := ValidateTargetFile(out)
	if diags.HasErrors() {
		return diags
	}

	var writer io.Writer

	for _, change := range plan.Changes.Resources {
		if len(change.GeneratedConfig) > 0 {

			if writer == nil {
				// Lazily create the generated file, in case we have no
				// generated config to create.
				var err error
				if writer, err = os.Create(out); err != nil {
					if os.IsPermission(err) {
						diags = diags.Append(tfdiags.Sourceless(
							tfdiags.Error,
							"Failed to create target generated file",
							fmt.Sprintf("Terraform did not have permission to create the generated file (%s) in the target directory. Please modify permissions over the target directory, and try again.", out)))
						return diags
					}

					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Failed to create target generated file",
						fmt.Sprintf("Terraform could not create the generated file (%s) in the target directory: %v. Depending on the error message, this may be a bug in Terraform itself. If so, please report it!", out, err)))
					return diags
				}

				header := "# __generated__ by Terraform\n# Please review these resources and move them into your main configuration files.\n"
				// Missing the header from the file, isn't the end of the world
				// so if this did return an error, then we will just ignore it.
				_, _ = writer.Write([]byte(header))
			}

			header := "\n# __generated__ by Terraform"
			if change.Importing != nil && len(change.Importing.ID) > 0 {
				header += fmt.Sprintf(" from %q", change.Importing.ID)
			}
			header += "\n"
			if _, err := writer.Write([]byte(fmt.Sprintf("%s%s\n", header, change.GeneratedConfig))); err != nil {
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Warning,
					"Failed to save generated config",
					fmt.Sprintf("Terraform encountered an error while writing generated config: %v. The config for %s must be created manually before applying. Depending on the error message, this may be a bug in Terraform itself. If so, please report it!", err, change.Addr.String())))
			}
		}
	}

	return diags
}
