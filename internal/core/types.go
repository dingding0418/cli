// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package core

import (
	"os"

	"github.com/larksuite/cli/internal/envvars"
)

// LarkBrand represents the Lark platform brand.
// "feishu" targets China-mainland, "lark" targets international.
// Any other string is treated as a custom base URL.
type LarkBrand string

const (
	BrandFeishu LarkBrand = "feishu"
	BrandLark   LarkBrand = "lark"
)

// ParseBrand normalizes a brand string to a LarkBrand constant.
// Unrecognized values default to BrandFeishu.
func ParseBrand(value string) LarkBrand {
	if value == "lark" {
		return BrandLark
	}
	return BrandFeishu
}

// Endpoints holds resolved endpoint URLs for different Lark services.
type Endpoints struct {
	Open     string // e.g. "https://open.feishu.cn"
	Accounts string // e.g. "https://accounts.feishu.cn"
	MCP      string // e.g. "https://mcp.feishu.cn"
}

// ResolveEndpoints resolves endpoint URLs based on brand.
func ResolveEndpoints(brand LarkBrand) Endpoints {
	switch brand {
	case BrandLark:
		return Endpoints{
			Open:     "https://open.larksuite.com",
			Accounts: "https://accounts.larksuite.com",
			MCP:      "https://mcp.larksuite.com",
		}
	default:
		return Endpoints{
			Open:     getenvOrDefault(envvars.CliFeishuOpenBaseURL, "https://open.feishu.cn"),
			Accounts: getenvOrDefault(envvars.CliFeishuAccountsBaseURL, "https://accounts.feishu.cn"),
			MCP:      getenvOrDefault(envvars.CliFeishuMCPBaseURL, "https://mcp.feishu.cn"),
		}
	}
}

// ResolveOpenBaseURL returns the Open API base URL for the given brand.
func ResolveOpenBaseURL(brand LarkBrand) string {
	return ResolveEndpoints(brand).Open
}

func getenvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
