package models

type Worker struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`
	Usage  int      `json:"usage"`
}
