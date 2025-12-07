package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出待处理的 Issue",
	Long: `列出带有 cursortoolset 标签的 Issue。

示例:
  github-issue list --repo owner/repo
  github-issue list --repo owner/repo --status pending
  github-issue list --repo owner/repo --type feature-request --format json`,
	RunE: runList,
}

var (
	listRepo   string
	listStatus string
	listType   string
	listLimit  int
	listFormat string
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listRepo, "repo", "", "目标仓库 (owner/repo)")
	listCmd.Flags().StringVar(&listStatus, "status", "pending", "状态过滤 (pending/processing/processed/all)")
	listCmd.Flags().StringVar(&listType, "type", "", "类型过滤")
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "数量限制")
	listCmd.Flags().StringVar(&listFormat, "format", "table", "输出格式 (table/json)")

	listCmd.MarkFlagRequired("repo")
}

func runList(cmd *cobra.Command, args []string) error {
	token := getToken(cmd)

	svc := service.NewIssueService(token)
	issues, err := svc.List(service.ListOptions{
		Repo:   listRepo,
		Status: listStatus,
		Type:   listType,
		Limit:  listLimit,
	})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		fmt.Println("没有找到符合条件的 Issue")
		return nil
	}

	if listFormat == "json" {
		data, _ := json.MarshalIndent(issues, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// 表格输出
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tType\tStatus\tTitle\tCreated")
	fmt.Fprintln(w, "---\t----\t------\t-----\t-------")
	for _, issue := range issues {
		title := issue.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			issue.Number, issue.Type, issue.Status, title, issue.CreatedAt)
	}
	w.Flush()

	return nil
}
