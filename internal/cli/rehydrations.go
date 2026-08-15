package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// rehydrationTimeFormat is the timestamp format the rehydration API expects
// (core.URLTimeFormat in the backend).
const rehydrationTimeFormat = "2006-01-02T15:04:05.000Z"

func newRehydrationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rehydrations",
		Short: "Rehydrate archived data back into pipelines",
		Long: `List, create, cancel and delete rehydration jobs that replay archived data
(S3, GCS, ...) back through a pipeline to a destination.

Creating a rehydration mirrors the web UI: the CQL filter and time range are
first resolved into the eligible archive-source/destination job combinations
("available rehydrations"), which are shown for confirmation before the jobs
are submitted as a batch.`,
	}
	cmd.AddCommand(
		newRehydrationsListCmd(),
		newRehydrationsGetCmd(),
		newRehydrationsValidateCmd(),
		newRehydrationsAnalyzeCmd(),
		newRehydrationsCreateCmd(),
		newRehydrationsCancelCmd(),
		newRehydrationsDeleteCmd(),
	)
	return cmd
}

func newRehydrationsListCmd() *cobra.Command {
	var query string
	var tf timeFlags
	var pg pageFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List rehydration jobs",
		Example: `  edx rehydrations list --output table --columns rehydration_id,status,percentage,from,to
  edx rehydrations list --query 'rehydration.status:"in-progress"' --lookback 24h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if query != "" {
				q.Set("query", query)
			}
			tf.apply(q)
			pg.apply(q)
			data, err := c.Get(cmdContext(cmd), "/rehydration_v2", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression over rehydration facets")
	tf.register(cmd, "")
	pg.register(cmd, 20)
	return cmd
}

func newRehydrationsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <rehydration-id>",
		Short: "Get a rehydration job with its progress percentage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/rehydration/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newRehydrationsValidateCmd() *cobra.Command {
	var query, source, destination string
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Show available rehydrations for a filter and time range",
		Long: `Resolve a CQL filter and time range into the eligible archive-source /
destination job combinations, with an efficiency rating per job. This is the
"Available Rehydrations" panel in the web UI and the dry-run half of
"edx rehydrations create".`,
		Example: `  edx rehydrations validate --query 'service.name:"api"' --lookback 1h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("query", query)
			if source != "" {
				q.Set("source", source)
			}
			if destination != "" {
				q.Set("destination", destination)
			}
			tf.apply(q)
			data, err := c.Get(cmdContext(cmd), "/rehydration/validate", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression (required)")
	_ = cmd.MarkFlagRequired("query")
	cmd.Flags().StringVar(&source, "source", "", "restrict to one archive source integration")
	cmd.Flags().StringVar(&destination, "destination", "", "override the mapped destination")
	tf.register(cmd, "1h")
	return cmd
}

func newRehydrationsAnalyzeCmd() *cobra.Command {
	var query, tag, source, bucket, file string
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Estimate the size and event count of a rehydration",
		Long: `Estimate how much data a rehydration would replay (approximate bytes and
event count) and list existing rehydrations that overlap the same range.

The API only supports archives backed by a legacy archive integration
(--source must name one). For pipeline-archive sources (archive_source, as
returned by "edx rehydrations validate") estimate instead with the metrics the
web UI uses, e.g.:
  edx metrics query --name ed.rehydration.bytes --agg sum --query '<cql>'`,
		Example: `  edx rehydrations analyze --query 'service.name:"api"' --tag prod --source my-s3 --lookback 1h
  edx rehydrations analyze --file rehydration.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			var body []byte
			if file != "" {
				if body, err = readFileOrStdin(file); err != nil {
					return err
				}
			} else {
				from, to, err := resolveTimeRange(tf, time.Now())
				if err != nil {
					return err
				}
				payload := map[string]any{"from": from, "to": to}
				for k, v := range map[string]string{
					"cql_query":          query,
					"tag":                tag,
					"source_integration": source,
					"bucket":             bucket,
				} {
					if v != "" {
						payload[k] = v
					}
				}
				if body, err = json.Marshal(payload); err != nil {
					return err
				}
			}
			data, err := c.Post(cmdContext(cmd), "/rehydration/analyze", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression")
	cmd.Flags().StringVar(&tag, "tag", "", "pipeline tag the archive belongs to")
	cmd.Flags().StringVar(&source, "source", "", "archive source integration")
	cmd.Flags().StringVar(&bucket, "bucket", "", "archive bucket name")
	cmd.Flags().StringVarP(&file, "file", "f", "", `full rehydration JSON body ("-" for stdin); overrides the other flags`)
	tf.register(cmd, "1h")
	return cmd
}

// potentialRehydration is one eligible job combo from /rehydration/validate.
type potentialRehydration struct {
	Source             string          `json:"source"`
	Bucket             string          `json:"bucket"`
	Destination        string          `json:"destination"`
	ArchiveSource      json.RawMessage `json:"archive_source,omitempty"`
	FilterErrorMessage string          `json:"filter_error_message,omitempty"`
	EfficiencyLevel    string          `json:"efficiency_level,omitempty"`
}

func newRehydrationsCreateCmd() *cobra.Command {
	var query string
	var excludeOverlap, dryRun bool
	var source, destination string
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create rehydration jobs for a filter and time range",
		Long: `Create rehydration jobs the same way the web UI does: the CQL filter and
time range are validated into the eligible archive-source/destination combos,
which are listed for confirmation and then submitted as a batch.

Use --dry-run to stop after the validation step, and --source/--destination to
narrow which of the eligible combos are created.`,
		Example: `  edx rehydrations create --query 'service.name:"api"' --from 2026-08-15T08:00:00.000Z --to 2026-08-15T09:00:00.000Z
  edx rehydrations create --query 'ed.tag:"prod"' --lookback 30m --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			ctx := cmdContext(cmd)
			from, to, err := resolveTimeRange(tf, time.Now())
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("query", query)
			q.Set("from", from)
			q.Set("to", to)
			if source != "" {
				q.Set("source", source)
			}
			if destination != "" {
				q.Set("destination", destination)
			}
			raw, err := c.Get(ctx, "/rehydration/validate", q)
			if err != nil {
				return err
			}
			if dryRun {
				return printResult(raw)
			}

			var validation struct {
				TotalJobs             int                    `json:"total_jobs"`
				Errors                []string               `json:"errors"`
				PotentialRehydrations []potentialRehydration `json:"potential_rehydrations"`
			}
			if err := json.Unmarshal(raw, &validation); err != nil {
				return fmt.Errorf("could not parse validation response: %w", err)
			}
			var jobs []potentialRehydration
			for _, p := range validation.PotentialRehydrations {
				if p.FilterErrorMessage != "" {
					fmt.Fprintf(os.Stderr, "skipping %s -> %s: %s\n", p.Source, p.Destination, p.FilterErrorMessage)
					continue
				}
				jobs = append(jobs, p)
			}
			if len(jobs) == 0 {
				if len(validation.Errors) > 0 {
					return fmt.Errorf("no eligible rehydrations: %s", strings.Join(validation.Errors, "; "))
				}
				return fmt.Errorf("no eligible rehydrations for this filter and time range")
			}

			fmt.Fprintf(os.Stderr, "%d rehydration job(s) for %s .. %s:\n", len(jobs), from, to)
			for _, p := range jobs {
				fmt.Fprintf(os.Stderr, "  %s\n", describePotentialRehydration(p))
			}
			if !confirm(fmt.Sprintf("Create %d rehydration job(s)?", len(jobs))) {
				return errAborted
			}

			batch := make([]map[string]any, 0, len(jobs))
			for _, p := range jobs {
				job := map[string]any{
					"source_integration": p.Source,
					"bucket":             p.Bucket,
					"destination":        p.Destination,
					"cql_query":          query,
					"from":               from,
					"to":                 to,
					"exclude_overlap":    excludeOverlap,
				}
				if len(p.ArchiveSource) > 0 {
					job["archive_source"] = p.ArchiveSource
				}
				batch = append(batch, job)
			}
			body, err := json.Marshal(batch)
			if err != nil {
				return err
			}
			bq := url.Values{}
			bq.Set("query", query)
			data, err := c.Post(ctx, "/rehydration/rehydration_batch", bq, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&query, "query", "q", "", "CQL filter expression (required)")
	_ = cmd.MarkFlagRequired("query")
	cmd.Flags().BoolVar(&excludeOverlap, "exclude-overlap", true, "skip data already rehydrated by previous jobs")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "only show the eligible jobs, do not create anything")
	cmd.Flags().StringVar(&source, "source", "", "restrict to one archive source integration")
	cmd.Flags().StringVar(&destination, "destination", "", "override the mapped destination")
	tf.register(cmd, "1h")
	return cmd
}

func newRehydrationsCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <rehydration-id>",
		Short: "Cancel a rehydration job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Cancel rehydration %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Put(cmdContext(cmd), "/rehydration/"+url.PathEscape(args[0]), nil, []byte(`{"status":"cancelled"}`))
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newRehydrationsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <rehydration-id>",
		Short: "Delete a rehydration job (running jobs are marked for deletion)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete rehydration %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/rehydration/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

// describePotentialRehydration renders one eligible job combo for the create
// confirmation prompt. Pipeline-archive entries carry their identity inside
// archive_source rather than the top-level source/bucket fields.
func describePotentialRehydration(p potentialRehydration) string {
	source, bucket := p.Source, p.Bucket
	if source == "" || bucket == "" {
		var as struct {
			Tag      string `json:"tag"`
			NodeName string `json:"node_name"`
			Bucket   string `json:"bucket"`
		}
		_ = json.Unmarshal(p.ArchiveSource, &as)
		if source == "" {
			source = strings.TrimPrefix(as.Tag+"/"+as.NodeName, "/")
		}
		if bucket == "" {
			bucket = as.Bucket
		}
	}
	eff := p.EfficiencyLevel
	if eff == "" {
		eff = "unknown"
	}
	return fmt.Sprintf("%s (bucket %s) -> %s [efficiency: %s]", source, bucket, p.Destination, eff)
}

// resolveTimeRange turns the shared time flags into the concrete from/to pair
// the rehydration create/analyze payloads require. Explicit --from wins over
// --lookback; a missing --to defaults to now.
func resolveTimeRange(t timeFlags, now time.Time) (from, to string, err error) {
	to = t.to
	if to == "" {
		to = now.UTC().Format(rehydrationTimeFormat)
	}
	if t.from != "" {
		return t.from, to, nil
	}
	d, err := time.ParseDuration(t.lookback)
	if err != nil {
		return "", "", fmt.Errorf("invalid --lookback %q: %w", t.lookback, err)
	}
	return now.UTC().Add(-d).Format(rehydrationTimeFormat), to, nil
}
