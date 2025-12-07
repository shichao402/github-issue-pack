package cli

import (
	"fmt"
	"strconv"

	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <issue-number>",
	Short: "关闭并标记 Issue 处理结果",
	Long: `关闭 Issue 并标记处理结果。

示例:
  github-issue close 123 --repo owner/repo --result success
  github-issue close 123 --repo owner/repo --result rejected --comment "不符合规范"`,
	Args: cobra.ExactArgs(1),
	RunE: runClose,
}

var (
	closeRepo    string
	closeResult  string
	closeComment string
)

func init() {
	rootCmd.AddCommand(closeCmd)

	closeCmd.Flags().StringVar(&closeRepo, "repo", "", "目标仓库 (owner/repo)")
	closeCmd.Flags().StringVar(&closeResult, "result", "", "处理结果 (success/rejected)")
	closeCmd.Flags().StringVar(&closeComment, "comment", "", "处理说明")

	closeCmd.MarkFlagRequired("repo")
	closeCmd.MarkFlagRequired("result")
}

func runClose(cmd *cobra.Command, args []string) error {
	token := getToken(cmd)

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("无效的 Issue 编号: %s", args[0])
	}

	// 验证结果
	if closeResult != "success" && closeResult != "rejected" {
		return fmt.Errorf("无效的结果: %s，只能是 success 或 rejected", closeResult)
	}

	svc := service.NewIssueService(token)
	if err := svc.Close(closeRepo, number, closeResult, closeComment); err != nil {
		return err
	}

	status := "processed"
	if closeResult == "rejected" {
		status = "rejected"
	}
	fmt.Printf("✅ Issue #%d 已关闭，状态: %s\n", number, status)
	return nil
}
