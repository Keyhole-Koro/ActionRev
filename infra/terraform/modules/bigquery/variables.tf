variable "project_id"  { type = string }
variable "region"      { type = string }
variable "dataset_id"  { type = string; default = "graph" }
variable "labels"      { type = map(string) }
