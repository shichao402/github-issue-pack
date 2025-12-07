package cli

import (
	"fmt"
	"strconv"

	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <issue-number>",
	Short: "更新 Issue 状态",
	Long: `更新 Issue 的状态标签。

示例:
  github-issue update 123 --repo owner/repo --status processing
  github-issue update 123 --repo owner/repo --status pending --comment "需要更多信息"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

var (
	updateRepo    string
	updateStatus  string
	updateComment string
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateRepo, "repo", "", "目标仓库 (owner/repo)")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "新状态 (processing/pending)")
	updateCmd.Flags().StringVar(&updateComment, "comment", "", "添加评论")

	updateCmd.MarkFlagRequired("repo")
	updateCmd.MarkFlagRequired("status")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	token := getToken(cmd)

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("无效的 Issue 编号: %s", args[0])
	}

	// 验证状态
	if updateStatus != "processing" && updateStatus != "pending" {
		return fmt.Errorf("无效的状态: %s，只能是 processing 或 pending", updateStatus)
	}

	svc := service.NewIssueService(token)
	if err := svc.UpdateStatus(updateRepo, number, updateStatus, updateComment); err != nil {
		return err
	}

	fmt.Printf("✅ Issue #%d 状态已更新为 %s\n", number, updateStatus)
	return nil
}
