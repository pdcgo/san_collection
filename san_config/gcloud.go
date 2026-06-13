package san_config

import "fmt"

type ProjectConfig struct {
	ProjectID string `env:"GOOGLE_CLOUD_PROJECT"`
}

func (projectCfg *ProjectConfig) PubsubTopicPath(topic string) string {
	return fmt.Sprintf("projects/%s/topics/%s", projectCfg.ProjectID, topic)
}

func (projectCfg *ProjectConfig) PubsubSubscriberPath(sub string) string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", projectCfg.ProjectID, sub)
}
