// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/json"
	"fmt"
	"strings"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

type documentRef struct {
	Kind  string
	Token string
}

func parseDocumentRef(input string) (documentRef, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return documentRef{}, output.ErrValidation("--doc cannot be empty")
	}

	if token, ok := extractDocumentToken(raw, "/wiki/"); ok {
		return documentRef{Kind: "wiki", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/docx/"); ok {
		return documentRef{Kind: "docx", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/doc/"); ok {
		return documentRef{Kind: "doc", Token: token}, nil
	}
	if strings.Contains(raw, "://") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx URL/token or a wiki URL that resolves to docx", raw)
	}
	if strings.ContainsAny(raw, "/?#") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx token or a wiki URL", raw)
	}

	return documentRef{Kind: "docx", Token: raw}, nil
}

func extractDocumentToken(raw, marker string) (string, bool) {
	idx := strings.Index(raw, marker)
	if idx < 0 {
		return "", false
	}
	token := raw[idx+len(marker):]
	if end := strings.IndexAny(token, "/?#"); end >= 0 {
		token = token[:end]
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}

// doDocAPI executes an OpenAPI request against the docs_ai endpoints and returns
// the parsed "data" field from the standard Lark response envelope {code, msg, data}.
func doDocAPI(runtime *common.RuntimeContext, method, apiPath string, body interface{}) (map[string]interface{}, error) {
	resp, err := runtime.DoAPI(&larkcore.ApiReq{
		HttpMethod: method,
		ApiPath:    apiPath,
		Body:       body,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		if len(resp.RawBody) > 0 {
			var errEnv struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}
			if json.Unmarshal(resp.RawBody, &errEnv) == nil && errEnv.Msg != "" {
				return nil, output.ErrAPI(errEnv.Code, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, errEnv.Msg), nil)
			}
		}
		return nil, output.ErrAPI(resp.StatusCode, fmt.Sprintf("HTTP %d", resp.StatusCode), nil)
	}
	if len(resp.RawBody) == 0 {
		return nil, fmt.Errorf("empty response body")
	}
	var envelope struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(resp.RawBody, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if envelope.Code != 0 {
		return nil, output.ErrAPI(envelope.Code, envelope.Msg, nil)
	}
	return envelope.Data, nil
}

func buildDriveRouteExtra(docID string) (string, error) {
	extra, err := json.Marshal(map[string]string{"drive_route_token": docID})
	if err != nil {
		return "", output.Errorf(output.ExitInternal, "internal_error", "failed to marshal upload extra data: %v", err)
	}
	return string(extra), nil
}
