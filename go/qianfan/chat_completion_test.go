package qianfan

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatCompletion(t *testing.T) {
	client, err := NewClientFromEnv()
	assert.NoError(t, err)
	for model, endpoint := range ChatModelEndpoint {
		chat := client.ChatCompletionFromModel(model)
		resp, err := chat.Do(context.Background(), &ChatCompletionRequest{
			Messages: []ChatCompletionMessage{
				ChatCompletionUserMessage("你好"),
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.RawResponse.StatusCode, 200)
		assert.NotEqual(t, resp.Id, nil)
		assert.Equal(t, resp.Object, "chat.completion")
		assert.True(t, strings.Contains(resp.RawResponse.Request.URL.Path, endpoint))
		assert.True(t, strings.Contains(resp.Result, "你好"))
	}
}

func TestChatCompletionStream(t *testing.T) {
	client, err := NewClientFromEnv()
	assert.NoError(t, err)
	for model, endpoint := range ChatModelEndpoint {
		chat := client.ChatCompletionFromModel(model)
		resp, err := chat.DoStream(context.Background(), &ChatCompletionRequest{
			Messages: []ChatCompletionMessage{
				ChatCompletionUserMessage("你好"),
			},
		})
		assert.NoError(t, err)
		turn_count := 0
		for {
			resp, err := resp.Recv()
			assert.NoError(t, err)
			if resp.IsEnd {
				break
			}
			turn_count++
			assert.Equal(t, resp.RawResponse.StatusCode, 200)
			assert.NotEqual(t, resp.Id, nil)
			assert.Equal(t, resp.Object, "chat.completion")
			assert.True(t, strings.Contains(resp.RawResponse.Request.URL.Path, endpoint))
			assert.True(t, strings.Contains(resp.Result, "你好"))
		}
		assert.True(t, turn_count > 1)
	}
}

func TestMain(m *testing.M) {
	// fmt.Println("TestMain")
	os.Setenv("QIANFAN_BASE_URL", "http://127.0.0.1:8866")
	os.Setenv("QIANFAN_ACCESS_KEY", "test_access_key")
	os.Setenv("QIANFAN_SECRET_KEY", "test_secret_key")

	os.Exit(m.Run())
}