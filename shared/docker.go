package shared

// ContainerStats struct
type ContainerStats struct {
	RAMUsage float64 `json:"ram_usage"`
	MaxRAM   float64 `json:"max_ram"`
	CPU      float64 `json:"cpu_usage"`
}
