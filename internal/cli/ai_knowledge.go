package cli

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/edgedelta/edx/internal/api"
)

func newAIKnowledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knowledge",
		Short: "Explore the AI Team knowledge graph (read-only)",
		Long: `Explore the organization's knowledge graph: services, repos, teams, people,
channels, incidents and the relationships between them, as discovered from
your connectors (GitHub, PagerDuty, Jira, AWS, Slack, ...).

All commands are read-only. Served by the agent service
(agent.ai.edgedelta.com).

Entity IDs have the form {orgId}::{type}::{externalId}; find them with
"edx ai knowledge search".`,
	}
	cmd.AddCommand(
		newAIKnowledgeStatsCmd(),
		newAIKnowledgeTopologyCmd(),
		newAIKnowledgeSearchCmd(),
		newAIKnowledgeGetCmd(),
		newAIKnowledgeSubgraphCmd(),
		newAIKnowledgeBlastRadiusCmd(),
		newAIKnowledgeCriticalityCmd(),
	)
	return cmd
}

// kgGet performs a GET against a knowledge-graph path on the agent service.
func kgGet(cmd *cobra.Command, path string, q url.Values) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	data, err := c.GetFrom(cmdContext(cmd), api.ServiceAgent, "/knowledge-graph"+path, q)
	if err != nil {
		return err
	}
	return printResult(data)
}

// kgEntityPath builds an /entities/{id}... path. Entity IDs may contain "/"
// (e.g. AWS ARNs), so the ID must be escaped into a single path segment.
func kgEntityPath(entityID, suffix string) string {
	return "/entities/" + url.PathEscape(entityID) + suffix
}

func newAIKnowledgeStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Node and edge counts by type and source, plus last sync time",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return kgGet(cmd, "/stats", nil)
		},
	}
}

func newAIKnowledgeTopologyCmd() *cobra.Command {
	var limit int
	var namespaces string
	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Fetch a slice of the graph: nodes, edges and stats",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if namespaces != "" {
				q.Set("namespaces", namespaces)
			}
			return kgGet(cmd, "/topology", q)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum nodes to return, most-observed first (server default 500, max 5000)")
	cmd.Flags().StringVar(&namespaces, "namespaces", "", "comma-separated namespaces: topology, learned (server default: topology)")
	return cmd
}

func newAIKnowledgeSearchCmd() *cobra.Command {
	var types, source, namespaces, cursor string
	var minConfidence float64
	var limit int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search entities by name or alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("q", args[0])
			if types != "" {
				q.Set("types", types)
			}
			if cmd.Flags().Changed("min-confidence") {
				q.Set("minConfidence", strconv.FormatFloat(minConfidence, 'f', -1, 64))
			}
			if source != "" {
				q.Set("source", source)
			}
			if namespaces != "" {
				q.Set("namespaces", namespaces)
			}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			return kgGet(cmd, "/search", q)
		},
	}
	cmd.Flags().StringVar(&types, "types", "", "comma-separated node types (Org, Integration, Service, Repo, Channel, JiraProject, PagerDutyService, AwsResource, Team, Person, Incident, Document)")
	cmd.Flags().Float64Var(&minConfidence, "min-confidence", 0, "minimum confidence score (0..1)")
	cmd.Flags().StringVar(&source, "source", "", "filter by provenance source")
	cmd.Flags().StringVar(&namespaces, "namespaces", "", "comma-separated namespaces: topology, learned")
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum matches to return (server default 25, max 200)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor (nextCursor from a previous response)")
	return cmd
}

func newAIKnowledgeGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <entity-id>",
		Short: "Get one entity with its neighbors and edges",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return kgGet(cmd, kgEntityPath(args[0], ""), nil)
		},
	}
}

func newAIKnowledgeSubgraphCmd() *cobra.Command {
	var hops int
	cmd := &cobra.Command{
		Use:   "subgraph <entity-id>",
		Short: "Fetch the N-hop neighborhood around an entity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if hops > 0 {
				q.Set("hops", strconv.Itoa(hops))
			}
			return kgGet(cmd, kgEntityPath(args[0], "/subgraph"), q)
		},
	}
	cmd.Flags().IntVar(&hops, "hops", 0, "neighborhood depth (server default 2, max 3)")
	return cmd
}

func newAIKnowledgeBlastRadiusCmd() *cobra.Command {
	var maxHops int
	cmd := &cobra.Command{
		Use:   "blast-radius <entity-id>",
		Short: "What is affected if this entity fails (follows dependency edges)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if maxHops > 0 {
				q.Set("maxHops", strconv.Itoa(maxHops))
			}
			return kgGet(cmd, kgEntityPath(args[0], "/blast-radius"), q)
		},
	}
	cmd.Flags().IntVar(&maxHops, "max-hops", 0, "maximum dependency distance to traverse (server default 3)")
	return cmd
}

func newAIKnowledgeCriticalityCmd() *cobra.Command {
	var limit int
	var namespaces string
	cmd := &cobra.Command{
		Use:   "criticality",
		Short: "Rank entities by how many others depend on them",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if namespaces != "" {
				q.Set("namespaces", namespaces)
			}
			return kgGet(cmd, "/criticality", q)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum entities to return (server default 25)")
	cmd.Flags().StringVar(&namespaces, "namespaces", "", "comma-separated namespaces: topology, learned")
	return cmd
}
