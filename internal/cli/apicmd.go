package cli

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

// errAborted is returned when the user declines a confirmation prompt.
var errAborted = errors.New("aborted")

func newAPICmd() *cobra.Command {
	var dataArg string
	var ep extraParams
	cmd := &cobra.Command{
		Use:   "api <METHOD> <path>",
		Short: "Call any Edge Delta API endpoint directly",
		Long: `Escape hatch for endpoints not yet covered by a dedicated edx command.

The path may contain {org} (or {org_id}), which is replaced with the active
organization ID. Paths not starting with /v1 or /public are treated as
org-relative and prefixed with /v1/orgs/{org}.

Request bodies: pass --data with inline JSON, or @file / @- to read from a
file or stdin.`,
		Example: `  edx api GET /v1/orgs/{org}/dashboards
  edx api GET /tokens                              # org-relative shorthand
  edx api POST /pipelines/{conf-id}/save --data @save.json
  edx api GET /users --param limit=10`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			method := strings.ToUpper(args[0])
			path := args[1]
			path = strings.ReplaceAll(path, "{org}", c.OrgID)
			path = strings.ReplaceAll(path, "{org_id}", c.OrgID)
			if !strings.HasPrefix(path, "/v1") && !strings.HasPrefix(path, "/public") && !strings.HasPrefix(path, "/v2") {
				path = c.OrgPath(path)
			}

			var body []byte
			if dataArg != "" {
				if strings.HasPrefix(dataArg, "@") {
					body, err = readFileOrStdin(strings.TrimPrefix(dataArg, "@"))
					if err != nil {
						return err
					}
				} else {
					body = []byte(dataArg)
				}
			}

			q := url.Values{}
			if err := ep.apply(q); err != nil {
				return err
			}

			switch method {
			case "GET", "POST", "PUT", "DELETE", "PATCH":
			default:
				return fmt.Errorf("unsupported method %q", method)
			}

			data, err := c.Do(cmdContext(cmd), method, path, q, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&dataArg, "data", "d", "", "request body: inline JSON, @file, or @- for stdin")
	ep.register(cmd)
	return cmd
}
