package shared

// ContainerStats struct
type ContainerStats struct {
	RAMUsage int64   `json:"ram_usage"`
	MaxRAM   int64   `json:"max_ram"`
	CPU      float64 `json:"cpu_usage"`
}
