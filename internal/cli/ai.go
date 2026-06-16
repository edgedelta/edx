package cli

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

func newAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI Teammate: issues, threads, channels, teammates and connectors",
		Long: `Manage the AI Teammate product.

Most AI Teammate data lives on dedicated service hosts, not the main API:
  issues, threads, channels   chat service   (chat.ai.edgedelta.com)
  teammates (agents)          agent service  (agent.ai.edgedelta.com)
  connectors, activity        main API       (api.edgedelta.com)

The active environment (prod, staging or local) selects which hosts are used,
so "edx ai issues list --env staging" hits the staging chat service. The issues
API is rolling out to production; use "--env staging" until it lands in prod.`,
	}
	cmd.AddCommand(
		newAIIssuesCmd(),
		newAIChannelsCmd(),
		newAIThreadsCmd(),
		newAIAgentsCmd(),
		newAIConnectorsCmd(),
		newAIActivityCmd(),
	)
	return cmd
}

// aiPageFlags are the pagination parameters accepted by the AI Teammate
// services. Unlike the observability endpoints, order is "ascending" or
// "descending"; an empty value lets the server choose.
type aiPageFlags struct {
	limit  int
	cursor string
	order  string
}

func (p *aiPageFlags) register(cmd *cobra.Command) {
	cmd.Flags().IntVar(&p.limit, "limit", 0, "maximum number of results (0 = server default)")
	cmd.Flags().StringVar(&p.cursor, "cursor", "", "pagination cursor (nextCursor from a previous response)")
	cmd.Flags().StringVar(&p.order, "order", "", "sort order: ascending or descending")
}

func (p *aiPageFlags) apply(q url.Values) {
	if p.limit > 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.cursor != "" {
		q.Set("cursor", p.cursor)
	}
	if p.order != "" {
		q.Set("order", p.order)
	}
}

// --- Issues (chat service) ---

func newAIIssuesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "List and manage issues raised by AI Teammates",
		Long: `Issues are problems AI Teammates have surfaced for your organization.
Served by the chat service; the API is rolling out to production, so use
"--env staging" until it is live in prod.`,
	}
	cmd.AddCommand(
		newAIIssuesListCmd(),
		newAIIssuesGetCmd(),
		newAIIssuesThreadsCmd(),
		newAIIssuesCriticalCmd(),
		newAIIssuesCloseCmd(),
		newAIIssuesDeleteCmd(),
		newAIHealthScoreCmd(),
		newAIHealthTimelineCmd(),
	)
	return cmd
}

func newAIIssuesListCmd() *cobra.Command {
	var withThreads, includeClosed bool
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues for the organization (open issues by default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			if withThreads {
				q.Set("withThreads", "true")
			}
			if includeClosed {
				q.Set("onlyOpenIssues", "false")
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().BoolVar(&withThreads, "with-threads", false, "include the threads attached to each issue")
	cmd.Flags().BoolVar(&includeClosed, "include-closed", false, "include closed issues (default: open only)")
	page.register(cmd)
	return cmd
}

func newAIIssuesGetCmd() *cobra.Command {
	var withThreads bool
	cmd := &cobra.Command{
		Use:   "get <issue-id>",
		Short: "Get a single issue by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if withThreads {
				q.Set("withThreads", "true")
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues/"+url.PathEscape(args[0]), q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().BoolVar(&withThreads, "with-threads", false, "include the issue's threads")
	return cmd
}

func newAIIssuesThreadsCmd() *cobra.Command {
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "threads <issue-id>",
		Short: "List the threads attached to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues/"+url.PathEscape(args[0])+"/threads", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	page.register(cmd)
	return cmd
}

func newAIIssuesCriticalCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "critical",
		Short: "List the most critical open issues, ranked",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues/most-critical", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "number of issues to return (1-100)")
	return cmd
}

func newAIIssuesCloseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close <issue-id>",
		Short: "Close an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.PostFrom(cmdContext(cmd), api.ServiceChat, "/issues/"+url.PathEscape(args[0])+"/close", nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newAIIssuesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <issue-id>",
		Short: "Delete an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete issue %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.DeleteFrom(cmdContext(cmd), api.ServiceChat, "/issues/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newAIHealthScoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Show the current health score (derived from open issues)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues/health/current-score", nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newAIHealthTimelineCmd() *cobra.Command {
	var lookback string
	cmd := &cobra.Command{
		Use:   "timeline",
		Short: "Show the health-score timeline over a lookback window",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if lookback != "" {
				q.Set("lookback", lookback)
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/issues/health/timeline", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&lookback, "lookback", "1h", "lookback window (e.g. 1h, 24h, 7d)")
	return cmd
}

// --- Channels (chat service) ---

func newAIChannelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channels",
		Short: "List AI Teammate channels and their messages",
	}
	cmd.AddCommand(newAIChannelsListCmd(), newAIChannelsGetCmd(), newAIChannelsMessagesCmd())
	return cmd
}

func newAIChannelsListCmd() *cobra.Command {
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List channels for the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/channels", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	page.register(cmd)
	return cmd
}

func newAIChannelsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <channel-id>",
		Short: "Get a single channel by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/channels/"+url.PathEscape(args[0]), nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
}

func newAIChannelsMessagesCmd() *cobra.Command {
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "messages <channel-id>",
		Short: "List messages in a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, "/channels/"+url.PathEscape(args[0])+"/messages", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	page.register(cmd)
	return cmd
}

// --- Threads (chat service, scoped to a channel) ---

func newAIThreadsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threads",
		Short: "List and inspect threads within a channel",
		Long: `Threads live inside a channel, so these commands take --channel <channel-id>.
Find channel IDs with "edx ai channels list".`,
	}
	cmd.AddCommand(newAIThreadsListCmd(), newAIThreadsGetCmd(), newAIThreadsMessagesCmd(), newAIThreadsDeleteCmd())
	return cmd
}

func threadsBase(channelID string) string {
	return "/channels/" + url.PathEscape(channelID) + "/threads"
}

func newAIThreadsListCmd() *cobra.Command {
	var channel, participant string
	var messageLimit int
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "list --channel <channel-id>",
		Short: "List threads in a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			if messageLimit > 0 {
				q.Set("messageLimit", strconv.Itoa(messageLimit))
			}
			if participant != "" {
				q.Set("participant", participant)
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, threadsBase(channel), q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&channel, "channel", "", "channel ID (required)")
	cmd.Flags().StringVar(&participant, "participant", "", "filter to threads with this participant")
	cmd.Flags().IntVar(&messageLimit, "message-limit", 0, "include up to N recent messages per thread")
	page.register(cmd)
	_ = cmd.MarkFlagRequired("channel")
	return cmd
}

func newAIThreadsGetCmd() *cobra.Command {
	var channel string
	var messageLimit int
	cmd := &cobra.Command{
		Use:   "get --channel <channel-id> <thread-id>",
		Short: "Get a single thread by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if messageLimit > 0 {
				q.Set("messageLimit", strconv.Itoa(messageLimit))
			}
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, threadsBase(channel)+"/"+url.PathEscape(args[0]), q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&channel, "channel", "", "channel ID (required)")
	cmd.Flags().IntVar(&messageLimit, "message-limit", 0, "include up to N recent messages")
	_ = cmd.MarkFlagRequired("channel")
	return cmd
}

func newAIThreadsMessagesCmd() *cobra.Command {
	var channel string
	var page aiPageFlags
	cmd := &cobra.Command{
		Use:   "messages --channel <channel-id> <thread-id>",
		Short: "List messages in a thread",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			page.apply(q)
			data, err := c.GetFrom(cmdContext(cmd), api.ServiceChat, threadsBase(channel)+"/"+url.PathEscape(args[0])+"/messages", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&channel, "channel", "", "channel ID (required)")
	page.register(cmd)
	_ = cmd.MarkFlagRequired("channel")
	return cmd
}

func newAIThreadsDeleteCmd() *cobra.Command {
	var channel string
	cmd := &cobra.Command{
		Use:   "delete --channel <channel-id> <thread-id>",
		Short: "Delete a thread",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm(fmt.Sprintf("Delete thread %s?", args[0])) {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.DeleteFrom(cmdContext(cmd), api.ServiceChat, threadsBase(channel)+"/"+url.PathEscape(args[0]), nil, nil)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVar(&channel, "channel", "", "channel ID (required)")
	_ = cmd.MarkFlagRequired("channel")
	return cmd
}

// --- Teammates / agents (agent service) ---

func newAIAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agents",
		Aliases: []string{"teammates"},
		Short:   "List and inspect AI Teammates (agents)",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List AI Teammates for the organization",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.GetFrom(cmdContext(cmd), api.ServiceAgent, "/agents", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "get <agent-id>",
			Short: "Get a single AI Teammate by ID",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.GetFrom(cmdContext(cmd), api.ServiceAgent, "/agents/"+url.PathEscape(args[0]), nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
	)
	return cmd
}

// --- Connectors & activity (main API) ---

func newAIConnectorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectors",
		Short: "Manage AI Teammate connectors",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List configured AI connectors",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "specs",
			Short: "List available connector types and their configuration specs",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors/specs", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		&cobra.Command{
			Use:   "environments",
			Short: "List environments available for connector deployment",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				data, err := c.Get(cmdContext(cmd), "/ai/connectors/environments", nil)
				if err != nil {
					return err
				}
				return printResult(data)
			},
		},
		newAIConnectorsUpdateCmd(),
		newAIConnectorsDeleteCmd(),
	)
	return cmd
}

func newAIConnectorsUpdateCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "update --file connector.json",
		Short: "Create or update a connector from a JSON request",
		Long: `Create or update a connector. The request body follows the connector spec;
use "edx ai connectors specs" to see required fields per connector type.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Post(cmdContext(cmd), "/ai/connectors", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `connector request JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAIConnectorsDeleteCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "delete --file connector.json",
		Short: "Delete a connector (request body identifies the connector)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm("Delete this AI connector?") {
				return errAborted
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readFileOrStdin(file)
			if err != nil {
				return err
			}
			data, err := c.Delete(cmdContext(cmd), "/ai/connectors", nil, body)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", `connector request JSON file ("-" for stdin)`)
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAIActivityCmd() *cobra.Command {
	var ep extraParams
	var tf timeFlags
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Show AI Teammate activity metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			tf.apply(q)
			if err := ep.apply(q); err != nil {
				return err
			}
			data, err := c.Get(cmdContext(cmd), "/ai/activity_metrics", q)
			if err != nil {
				return err
			}
			return printResult(data)
		},
	}
	tf.register(cmd, "24h")
	ep.register(cmd)
	return cmd
}
