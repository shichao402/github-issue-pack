package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shichao402/github-issue-pack/internal/models"
	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建标准化的 GitHub Issue",
	Long: `创建标准化的 GitHub Issue，自动打包内容到 Gist。

示例:
  github-issue create --repo owner/repo --type feature-request --title "添加新功能"
  github-issue create --repo owner/repo --type bug-report --title "修复问题" --payload request.json`,
	RunE: runCreate,
}

var (
	createRepo    string
	createType    string
	createTitle   string
	createPayload string
	createAttach  []string
	createDryRun  bool
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createRepo, "repo", "", "目标仓库 (owner/repo)")
	createCmd.Flags().StringVar(&createType, "type", "", "Issue 类型 (feature-request/bug-report/pack-register/pack-sync)")
	createCmd.Flags().StringVar(&createTitle, "title", "", "Issue 标题")
	createCmd.Flags().StringVar(&createPayload, "payload", "", "详细内容文件路径 (JSON)")
	createCmd.Flags().StringSliceVar(&createAttach, "attach", nil, "附件文件路径")
	createCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "预览模式，不实际创建")

	createCmd.MarkFlagRequired("repo")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("title")
}

func runCreate(cmd *cobra.Command, args []string) error {
	token := getToken(cmd)

	// 验证 Issue 类型
	issueType := models.IssueType(createType)
	switch issueType {
	case models.TypeFeatureRequest, models.TypeBugReport,
		models.TypePackRegister, models.TypePackSync,
		models.TypeQuestion, models.TypeCustom:
		// valid
	default:
		return fmt.Errorf("无效的 Issue 类型: %s", createType)
	}

	// 读取 payload
	var payload interface{}
	if createPayload != "" {
		data, err := os.ReadFile(createPayload)
		if err != nil {
			return fmt.Errorf("读取 payload 文件失败: %w", err)
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return fmt.Errorf("解析 payload JSON 失败: %w", err)
		}
	} else {
		// 使用默认 payload
		payload = map[string]string{
			"title":       createTitle,
			"description": "",
		}
	}

	// 读取附件
	var attachments []models.Attachment
	for _, path := range createAttach {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取附件失败 %s: %w", path, err)
		}
		attachments = append(attachments, models.Attachment{
			Name:    path,
			Content: string(data),
		})
	}

	// 创建 Issue
	svc := service.NewIssueService(token)
	result, err := svc.Create(service.CreateIssueOptions{
		Repo:        createRepo,
		Type:        issueType,
		Title:       createTitle,
		Payload:     payload,
		Attachments: attachments,
		DryRun:      createDryRun,
	})
	if err != nil {
		return err
	}

	if !createDryRun {
		fmt.Println("✅ Issue 创建成功!")
		fmt.Printf("   Issue: %s\n", result.IssueURL)
		fmt.Printf("   Gist:  %s\n", result.GistURL)
	}

	return nil
}
