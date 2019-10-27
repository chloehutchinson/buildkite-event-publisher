package main

import (
	"log"
	"strings"
)

// build.running build.finished
func handleBuildEvent(event buildkiteEvent) map[string]interface{} {
	if event.Build == nil {
		return nil
	}

	data := map[string]interface{}{
		"BuildID":      event.Build.ID,
		"BuildNumber":  event.Build.Number,
		"PipelineID":   event.Pipeline.ID,
		"PipelineSlug": event.Pipeline.Slug,
		"State":        event.Build.State,
		"Blocked":      event.Build.Blocked,
		"Branch":       event.Build.Branch,
		"Message":      event.Build.Message,
		"CreatedBy":    event.Build.Creator.Name,
		"WebURL":       event.Build.WebURL,
		"RebuiltFrom":  event.Build.RebuiltFrom,
	}
	if event.Build.StartedAt != nil && *event.Build.StartedAt != "" {
		data["StartedAt"] = *event.Build.StartedAt
	}
	if event.Build.FinishedAt != nil && *event.Build.FinishedAt != "" {
		data["FinishedAt"] = *event.Build.FinishedAt
	}

	return data
}

// job.started job.finished job.activated
func handleJobEvent(event buildkiteEvent) map[string]interface{} {
	data := map[string]interface{}{
		"JobID":        event.Job.ID,
		"JobState":     event.Job.State,
		"WebURL":       event.Job.WebURL,
		"BuildID":      event.Build.ID,
		"BuildNumber":  event.Build.Number,
		"BuildState":   event.Build.State,
		"PipelineID":   event.Pipeline.ID,
		"PipelineSlug": event.Pipeline.Slug,
		"Blocked":      event.Build.Blocked,
		"Branch":       event.Build.Branch,
		"Message":      event.Build.Message,
		"CreatedBy":    event.Build.Creator.Name,
	}
	if event.Job.Name != nil {
		data["JobName"] = *event.Job.Name
	}
	if event.Job.Command != nil {
		data["Command"] = *event.Job.Command
	}
	if event.Job.Retried != nil {
		data["Retried"] = *event.Job.Retried
	}
	if event.Job.RetriesCount != nil {
		data["RetriesCount"] = *event.Job.RetriesCount
	}
	if event.Job.ExitStatus != nil {
		data["ExitStatus"] = *event.Job.ExitStatus
	}
	if event.Job.StartedAt != nil && *event.Job.StartedAt != "" {
		data["StartedAt"] = *event.Job.StartedAt
	}
	if event.Job.FinishedAt != nil && *event.Job.FinishedAt != "" {
		data["FinishedAt"] = *event.Job.FinishedAt
	}
	if event.Job.Agent != nil {
		data["AgentID"] = event.Job.Agent.ID
		data["AgentName"] = event.Job.Agent.Name
		data["AgentWebURL"] = event.Job.Agent.WebURL
		data["AgentVersion"] = event.Job.Agent.Version
		data["AgentPriority"] = event.Job.Agent.Priority

		for _, v := range event.Job.Agent.Metadata {
			components := strings.SplitN(v, "=", 2)
			if len(components) < 2 {
				log.Printf("Ignoring agent metadata entry (insufficient data for mapping to a key = value): %s", v)
			}

			metaKey := "AgentMeta_" + components[0]
			data[metaKey] = components[1]
		}
	}

	return data
}
