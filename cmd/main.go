package main

import (
	"fmt"
	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

func main() {
	fmt.Println("GO RUNS")
	internal.Loader()
}
