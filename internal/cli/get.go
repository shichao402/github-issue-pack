package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <issue-number>",
	Short: "获取并解析指定 Issue",
	Long: `获取指定 Issue 并解析其中的结构化数据。

示例:
  github-issue get 123 --repo owner/repo
  github-issue get 123 --repo owner/repo --format json
  github-issue get 123 --repo owner/repo --output issue.json`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

var (
	getRepo   string
	getFormat string
	getOutput string
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getRepo, "repo", "", "目标仓库 (owner/repo)")
	getCmd.Flags().StringVar(&getFormat, "format", "json", "输出格式 (json/text)")
	getCmd.Flags().StringVar(&getOutput, "output", "", "输出到文件")

	getCmd.MarkFlagRequired("repo")
}

func runGet(cmd *cobra.Command, args []string) error {
	token := getToken(cmd)

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("无效的 Issue 编号: %s", args[0])
	}

	svc := service.NewIssueService(token)
	result, err := svc.Get(getRepo, number)
	if err != nil {
		return err
	}

	// 构建输出
	output := map[string]interface{}{
		"issue": map[string]interface{}{
			"number":     result.Issue.Number,
			"title":      result.Issue.Title,
			"state":      result.Issue.State,
			"url":        result.Issue.HTMLURL,
			"created_at": result.Issue.CreatedAt,
			"updated_at": result.Issue.UpdatedAt,
		},
	}

	if result.Package != nil {
		output["package"] = result.Package
	}

	var outputData []byte
	if getFormat == "text" {
		fmt.Printf("Issue #%d: %s\n", result.Issue.Number, result.Issue.Title)
		fmt.Printf("State: %s\n", result.Issue.State)
		fmt.Printf("URL: %s\n", result.Issue.HTMLURL)
		fmt.Printf("Created: %s\n", result.Issue.CreatedAt)
		if result.Package != nil {
			fmt.Printf("\n--- Package ---\n")
			fmt.Printf("Type: %s\n", result.Package.Type)
			fmt.Printf("Schema: %s\n", result.Package.Schema)
			pkgData, _ := json.MarshalIndent(result.Package, "", "  ")
			fmt.Println(string(pkgData))
		}
		return nil
	}

	outputData, _ = json.MarshalIndent(output, "", "  ")

	if getOutput != "" {
		if err := os.WriteFile(getOutput, outputData, 0644); err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
		fmt.Printf("已保存到 %s\n", getOutput)
	} else {
		fmt.Println(string(outputData))
	}

	return nil
}
