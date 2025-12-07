package models

import (
	"encoding/json"
	"time"
)

const SchemaVersion = "cursortoolset-issue-v1"

// IssueType Issue 类型
type IssueType string

const (
	TypeFeatureRequest IssueType = "feature-request"
	TypeBugReport      IssueType = "bug-report"
	TypePackRegister   IssueType = "pack-register"
	TypePackSync       IssueType = "pack-sync"
	TypeQuestion       IssueType = "question"
	TypeCustom         IssueType = "custom"
)

// IssuePackage Issue 包的完整结构（存储在 Gist 中）
type IssuePackage struct {
	Schema      string            `json:"$schema"`
	Meta        Meta              `json:"meta"`
	Type        IssueType         `json:"type"`
	Target      Target            `json:"target"`
	Payload     json.RawMessage   `json:"payload"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Meta 元数据
type Meta struct {
	CreatedAt           string `json:"created_at"`
	SourceProject       string `json:"source_project,omitempty"`
	CursorToolsetVersion string `json:"cursortoolset_version,omitempty"`
	GitHubIssueVersion  string `json:"github_issue_version"`
}

// Target 目标信息
type Target struct {
	Repo    string `json:"repo"`
	Pack    string `json:"pack,omitempty"`
	Version string `json:"version,omitempty"`
}

// Attachment 附件
type Attachment struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// FeatureRequestPayload 功能请求的 payload
type FeatureRequestPayload struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	UseCase          string `json:"use_case,omitempty"`
	ExpectedBehavior string `json:"expected_behavior,omitempty"`
	Alternatives     string `json:"alternatives,omitempty"`
}

// BugReportPayload Bug 报告的 payload
type BugReportPayload struct {
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	StepsToReproduce []string    `json:"steps_to_reproduce,omitempty"`
	ExpectedBehavior string      `json:"expected_behavior,omitempty"`
	ActualBehavior   string      `json:"actual_behavior,omitempty"`
	Environment      Environment `json:"environment,omitempty"`
}

// Environment 环境信息
type Environment struct {
	OS                   string `json:"os,omitempty"`
	CursorToolsetVersion string `json:"cursortoolset_version,omitempty"`
	PackVersion          string `json:"pack_version,omitempty"`
}

// PackRegisterPayload 包注册的 payload
type PackRegisterPayload struct {
	Repository  string `json:"repository"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// PackSyncPayload 包同步的 payload
type PackSyncPayload struct {
	Repository string `json:"repository"`
	Version    string `json:"version"`
	Changes    string `json:"changes,omitempty"`
}

// NewIssuePackage 创建新的 Issue 包
func NewIssuePackage(issueType IssueType, targetRepo string, payload interface{}) (*IssuePackage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &IssuePackage{
		Schema: SchemaVersion,
		Meta: Meta{
			CreatedAt:          time.Now().UTC().Format(time.RFC3339),
			GitHubIssueVersion: "0.1.0",
		},
		Type: issueType,
		Target: Target{
			Repo: targetRepo,
		},
		Payload: payloadBytes,
	}, nil
}

// GetPayload 解析 payload 到指定类型
func (p *IssuePackage) GetPayload(v interface{}) error {
	return json.Unmarshal(p.Payload, v)
}

// ToJSON 序列化为 JSON
func (p *IssuePackage) ToJSON() (string, error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseIssuePackage 从 JSON 解析 Issue 包
func ParseIssuePackage(data string) (*IssuePackage, error) {
	var pkg IssuePackage
	if err := json.Unmarshal([]byte(data), &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}
