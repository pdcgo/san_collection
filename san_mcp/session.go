package san_mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type McpSessionManager struct {
	client *redis.Client
}

type McpSession struct {
	Token     string `json:"token"`
	PdcSource string `json:"pdc_source"`
	UserId    uint64 `json:"user_id"`
	TeamId    uint64 `json:"team_id"`
}

func NewMcpSessionManager(client *redis.Client) *McpSessionManager {
	return &McpSessionManager{client}
}

func (m *McpSessionManager) GetSession(ctx context.Context, id string) (*McpSession, error) {
	var session McpSession
	sessionId := fmt.Sprintf("mcp_session:%s", id)
	val, err := m.client.Get(ctx, sessionId).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *McpSessionManager) SetSession(ctx context.Context, id string, session *McpSession) error {
	sessionId := fmt.Sprintf("mcp_session:%s", id)
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, sessionId, string(b), time.Hour*24).Err()
}
