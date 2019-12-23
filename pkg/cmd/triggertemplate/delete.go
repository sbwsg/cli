// Copyright © 2019 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package triggertemplate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/helper/options"
	"github.com/tektoncd/cli/pkg/helper/validate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cliopts "k8s.io/cli-runtime/pkg/genericclioptions"
)

func deleteCommand(p cli.Params) *cobra.Command {
	opts := &options.DeleteOptions{Resource: "triggertemplate", ForceDelete: false}
	f := cliopts.NewPrintFlags("delete")
	eg := `Delete a TriggerTemplate of name 'foo' in namespace 'bar'

    tkn triggertemplate delete foo -n bar

or

    tkn tt rm foo -n bar
`

	c := &cobra.Command{
		Use:          "delete",
		Aliases:      []string{"rm"},
		Short:        "Delete a triggertemplate in a namespace",
		Example:      eg,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		Annotations: map[string]string{
			"commandType": "main",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := &cli.Stream{
				In:  cmd.InOrStdin(),
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			if err := validate.NamespaceExists(p); err != nil {
				return err
			}

			if err := opts.CheckOptions(s, args[0]); err != nil {
				return err
			}

			return deleteTriggerTemplate(s, p, args[0])
		},
	}
	f.AddFlags(c)
	c.Flags().BoolVarP(&opts.ForceDelete, "force", "f", false, "Whether to force deletion (default: false)")

	_ = c.MarkZshCompPositionalArgumentCustom(1, "__tkn_get_triggertemplate")
	return c
}

func deleteTriggerTemplate(s *cli.Stream, p cli.Params, ttName string) error {
	cs, err := p.Clients()
	if err != nil {
		return fmt.Errorf("failed to create tekton client")
	}

	if err := cs.Triggers.TektonV1alpha1().TriggerTemplates(p.Namespace()).Delete(ttName, &metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete triggertemplate %q: %s", ttName, err)
	}

	fmt.Fprintf(s.Out, "TriggerTemplate deleted: %s\n", ttName)
	return nil
}
