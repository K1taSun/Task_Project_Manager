package app

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
)

// MetricsResponse zawiera komplet informacji dla dashboardu metryk.
type MetricsResponse struct {
	GeneratedAt          time.Time        `json:"generated_at"`
	Summary              SummaryMetrics   `json:"summary"`
	Projects             []ProjectMetrics `json:"projects"`
	Assignees            []AssigneeMetric `json:"assignees"`
	PriorityDistribution []PriorityBucket `json:"priority_distribution"`
	TagUsage             []TagUsage       `json:"tag_usage"`
	StatusBreakdown      map[string]int   `json:"status_breakdown"`
	Throughput           []Throughput     `json:"throughput"`
}

// SummaryMetrics prezentuje najważniejsze agregaty.
type SummaryMetrics struct {
	TotalProjects       int     `json:"total_projects"`
	ActiveProjects      int     `json:"active_projects"`
	CompletedProjects   int     `json:"completed_projects"`
	TotalTasks          int     `json:"total_tasks"`
	CompletedTasks      int     `json:"completed_tasks"`
	OpenTasks           int     `json:"open_tasks"`
	CompletionRate      float64 `json:"completion_rate"`
	OverdueTasks        int     `json:"overdue_tasks"`
	UpcomingTasks       int     `json:"upcoming_tasks"`
	TotalEstimatedHours float64 `json:"total_estimated_hours"`
	TotalActualHours    float64 `json:"total_actual_hours"`
	EstimateAccuracy    float64 `json:"estimate_accuracy"`
	AveragePriority     float64 `json:"average_priority"`
	Velocity7d          int     `json:"velocity_7d"`
}

// ProjectMetrics agreguje dane dla pojedynczego projektu.
type ProjectMetrics struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	Color            string  `json:"color,omitempty"`
	TasksTotal       int     `json:"tasks_total"`
	TasksCompleted   int     `json:"tasks_completed"`
	CompletionRate   float64 `json:"completion_rate"`
	OverdueTasks     int     `json:"overdue_tasks"`
	UpcomingTasks    int     `json:"upcoming_tasks"`
	AveragePriority  float64 `json:"average_priority"`
	EstimatedHours   float64 `json:"estimated_hours"`
	ActualHours      float64 `json:"actual_hours"`
	LastUpdated      string  `json:"last_updated"`
	CompletionHealth string  `json:"completion_health"`
}

// AssigneeMetric przechowuje obciążenie i realizację dla wykonawcy.
type AssigneeMetric struct {
	Name            string  `json:"name"`
	TotalTasks      int     `json:"total_tasks"`
	OpenTasks       int     `json:"open_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	EstimatedHours  float64 `json:"estimated_hours"`
	ActualHours     float64 `json:"actual_hours"`
	AveragePriority float64 `json:"average_priority"`
}

// PriorityBucket reprezentuje liczebność zadań w danym priorytecie.
type PriorityBucket struct {
	Level int `json:"level"`
	Count int `json:"count"`
}

// TagUsage określa popularność tagów.
type TagUsage struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// Throughput reprezentuje tempo pracy w horyzoncie czasowym.
type Throughput struct {
	Date      string `json:"date"`
	Created   int    `json:"created"`
	Completed int    `json:"completed"`
}

// metricsHandler przygotowuje i wysyła dane metryk jako JSON.
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	metrics := computeMetrics(time.Now())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func computeMetrics(now time.Time) MetricsResponse {
	mutex.RLock()
	defer mutex.RUnlock()

	projectList := make([]Project, 0, len(projects))
	for _, p := range projects {
		projectList = append(projectList, p)
	}

	taskList := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, t)
	}

	summary := SummaryMetrics{}
	statusCounts := map[string]int{
		"open":     0,
		"done":     0,
		"overdue":  0,
		"upcoming": 0,
	}

	priorityCounts := map[int]int{}
	tagCounts := map[string]int{}
	assigneeBuckets := map[string]*AssigneeMetric{}
	projectBuckets := map[int]*ProjectMetrics{}

	throughtputWindow := 14
	throughputMap := map[string]*Throughput{}
	for i := throughtputWindow - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		throughputMap[day] = &Throughput{
			Date: day,
		}
	}

	// Initialize project buckets
	for _, p := range projectList {
		projectBuckets[p.ID] = &ProjectMetrics{
			ID:     p.ID,
			Name:   p.Name,
			Status: p.Status,
			Color:  p.Color,
		}
		if strings.EqualFold(p.Status, "completed") {
			summary.CompletedProjects++
		} else if strings.EqualFold(p.Status, "active") || p.Status == "" {
			summary.ActiveProjects++
		}
	}
	summary.TotalProjects = len(projectList)

	var totalPriority float64
	var priorityCount float64

	for _, task := range taskList {
		summary.TotalTasks++
		if task.Done {
			summary.CompletedTasks++
			statusCounts["done"]++
		} else {
			summary.OpenTasks++
			statusCounts["open"]++
		}

		if task.Priority > 0 {
			priorityCounts[task.Priority]++
			totalPriority += float64(task.Priority)
			priorityCount++
		}

		if task.EstimatedHours > 0 {
			summary.TotalEstimatedHours += task.EstimatedHours
		}
		if task.ActualHours > 0 {
			summary.TotalActualHours += task.ActualHours
		}

		if task.Deadline != nil {
			deadline := *task.Deadline
			if !task.Done && deadline.Before(now) {
				summary.OverdueTasks++
				statusCounts["overdue"]++
			} else if !task.Done && deadline.After(now) && deadline.Before(now.AddDate(0, 0, 7)) {
				summary.UpcomingTasks++
				statusCounts["upcoming"]++
			}
		}

		for _, tag := range task.Tags {
			tag = strings.TrimSpace(strings.ToLower(tag))
			if tag == "" {
				continue
			}
			tagCounts[tag]++
		}

		assignee := strings.TrimSpace(task.Assignee)
		if assignee == "" {
			assignee = "Nieprzypisane"
		}
		if _, ok := assigneeBuckets[assignee]; !ok {
			assigneeBuckets[assignee] = &AssigneeMetric{
				Name: assignee,
			}
		}
		assigneeBucket := assigneeBuckets[assignee]
		assigneeBucket.TotalTasks++
		if task.Done {
			assigneeBucket.CompletedTasks++
		} else {
			assigneeBucket.OpenTasks++

		}
		if task.EstimatedHours > 0 {
			assigneeBucket.EstimatedHours += task.EstimatedHours
		}
		if task.ActualHours > 0 {
			assigneeBucket.ActualHours += task.ActualHours
		}
		if task.Priority > 0 {
			assigneeBucket.AveragePriority += float64(task.Priority)
		}

		if bucket, ok := projectBuckets[task.ProjectID]; ok {
			bucket.TasksTotal++
			if task.Done {
				bucket.TasksCompleted++
			}
			if task.Priority > 0 {
				bucket.AveragePriority += float64(task.Priority)
			}
			if task.EstimatedHours > 0 {
				bucket.EstimatedHours += task.EstimatedHours
			}
			if task.ActualHours > 0 {
				bucket.ActualHours += task.ActualHours
			}
			if task.Deadline != nil {
				deadline := *task.Deadline
				if !task.Done && deadline.Before(now) {
					bucket.OverdueTasks++
				} else if !task.Done && deadline.After(now) && deadline.Before(now.AddDate(0, 0, 7)) {
					bucket.UpcomingTasks++
				}
			}
			if bucket.LastUpdated == "" || bucket.LastUpdated < task.UpdatedAt.Format(time.RFC3339) {
				bucket.LastUpdated = task.UpdatedAt.Format(time.RFC3339)
			}
		}

		createdDay := task.CreatedAt.Format("2006-01-02")
		if entry, ok := throughputMap[createdDay]; ok {
			entry.Created++
		}
		if task.Done {
			completedDay := task.UpdatedAt.Format("2006-01-02")
			if entry, ok := throughputMap[completedDay]; ok {
				entry.Completed++
			}
			if task.UpdatedAt.After(now.AddDate(0, 0, -7)) {
				summary.Velocity7d++
			}
		}
	}

	if summary.TotalTasks > 0 {
		summary.CompletionRate = round2(float64(summary.CompletedTasks) / float64(summary.TotalTasks) * 100)
	}

	if summary.TotalEstimatedHours > 0 {
		summary.EstimateAccuracy = round2(summary.TotalActualHours / summary.TotalEstimatedHours * 100)
	} else {
		summary.EstimateAccuracy = 0
	}

	if priorityCount > 0 {
		summary.AveragePriority = round2(totalPriority / priorityCount)
	}

	assignees := make([]AssigneeMetric, 0, len(assigneeBuckets))
	for _, bucket := range assigneeBuckets {
		if bucket.TotalTasks > 0 && bucket.AveragePriority > 0 {
			bucket.AveragePriority = round2(bucket.AveragePriority / float64(bucket.TotalTasks))
		} else {
			bucket.AveragePriority = 0
		}
		bucket.EstimatedHours = round2(bucket.EstimatedHours)
		bucket.ActualHours = round2(bucket.ActualHours)
		assignees = append(assignees, *bucket)
	}
	sort.Slice(assignees, func(i, j int) bool {
		if assignees[i].TotalTasks == assignees[j].TotalTasks {
			return assignees[i].Name < assignees[j].Name
		}
		return assignees[i].TotalTasks > assignees[j].TotalTasks
	})

	projectMetrics := make([]ProjectMetrics, 0, len(projectBuckets))
	for _, bucket := range projectBuckets {
		if bucket.TasksTotal > 0 {
			bucket.CompletionRate = round2(float64(bucket.TasksCompleted) / float64(bucket.TasksTotal) * 100)
			if bucket.AveragePriority > 0 {
				bucket.AveragePriority = round2(bucket.AveragePriority / float64(bucket.TasksTotal))
			}
		}
		bucket.EstimatedHours = round2(bucket.EstimatedHours)
		bucket.ActualHours = round2(bucket.ActualHours)
		bucket.CompletionHealth = completionHealth(bucket.CompletionRate, bucket.OverdueTasks)
		projectMetrics = append(projectMetrics, *bucket)
	}
	sort.Slice(projectMetrics, func(i, j int) bool {
		if projectMetrics[i].TasksTotal == projectMetrics[j].TasksTotal {
			return projectMetrics[i].Name < projectMetrics[j].Name
		}
		return projectMetrics[i].TasksTotal > projectMetrics[j].TasksTotal
	})

	priorityBuckets := make([]PriorityBucket, 0, len(priorityCounts))
	for level := 1; level <= MaxPriority; level++ {
		priorityBuckets = append(priorityBuckets, PriorityBucket{
			Level: level,
			Count: priorityCounts[level],
		})
	}

	tagUsage := make([]TagUsage, 0, len(tagCounts))
	for tag, count := range tagCounts {
		tagUsage = append(tagUsage, TagUsage{
			Tag:   tag,
			Count: count,
		})
	}
	sort.Slice(tagUsage, func(i, j int) bool {
		if tagUsage[i].Count == tagUsage[j].Count {
			return tagUsage[i].Tag < tagUsage[j].Tag
		}
		return tagUsage[i].Count > tagUsage[j].Count
	})
	if len(tagUsage) > 12 {
		tagUsage = tagUsage[:12]
	}

	throughput := make([]Throughput, 0, len(throughputMap))
	keys := make([]string, 0, len(throughputMap))
	for day := range throughputMap {
		keys = append(keys, day)
	}
	sort.Strings(keys)
	for _, day := range keys {
		throughput = append(throughput, *throughputMap[day])
	}

	return MetricsResponse{
		GeneratedAt:          now,
		Summary:              summary,
		Projects:             projectMetrics,
		Assignees:            assignees,
		PriorityDistribution: priorityBuckets,
		TagUsage:             tagUsage,
		StatusBreakdown:      statusCounts,
		Throughput:           throughput,
	}
}

func completionHealth(completionRate float64, overdue int) string {
	switch {
	case overdue > 0:
		return "at_risk"
	case completionRate >= 80:
		return "on_track"
	case completionRate >= 50:
		return "needs_attention"
	default:
		return "behind"
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
