package main

import (
	"sync"
	"time"
)

type Project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
}

type Task struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Title     string    `json:"title"`
	Deadline  time.Time `json:"deadline"`
	Tags      []string  `json:"tags"`
	Priority  int       `json:"priority"`
	Done      bool      `json:"done"`
}

var (
	projects   = make(map[int]Project)
	tasks      = make(map[int]Task)
	mutex      sync.Mutex
) 