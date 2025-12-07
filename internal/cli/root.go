package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github-issue",
	Short: "标准化的 GitHub Issue 创建与处理工具",
	Long: `GitHub Issue Pack - 标准化的 GitHub Issue 创建与处理工具

支持项目间自动化协作，提供 Issue 打包/解包、状态管理、Gist 存储等功能。

Token 获取优先级:
  1. --token 参数
  2. GITHUB_TOKEN 环境变量
  3. gh CLI 认证信息 (自动获取)`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub Token (可选，默认使用 gh CLI 认证)")
}

func getToken(cmd *cobra.Command) string {
	// 1. 命令行参数
	token, _ := cmd.Flags().GetString("token")
	if token != "" {
		return token
	}

	// 2. 环境变量
	token = os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	// 3. 从 gh CLI 获取
	token = getGhToken()
	if token != "" {
		return token
	}

	fmt.Fprintln(os.Stderr, "错误: 无法获取 GitHub Token")
	fmt.Fprintln(os.Stderr, "请使用以下任一方式提供认证:")
	fmt.Fprintln(os.Stderr, "  1. --token 参数")
	fmt.Fprintln(os.Stderr, "  2. 设置 GITHUB_TOKEN 环境变量")
	fmt.Fprintln(os.Stderr, "  3. 运行 'gh auth login' 进行认证")
	os.Exit(1)
	return ""
}

// getGhToken 从 gh CLI 获取 token
func getGhToken() string {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
