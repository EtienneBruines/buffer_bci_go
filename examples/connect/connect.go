package main

import (
	"log"

	"github.com/EtienneBruines/go_buffer_bci_client/bufferbci"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

func main() {
	// Connecting
	log.Println("Trying to connect to localhost:1972 ...")

	conn, err := bufferbci.Connect("localhost:1972")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connection established")

	// Getting header information
	header, err := conn.GetHeader()
	if err != nil {
		log.Fatal(err)
	}

	// Getting samples
	var amountOfSamples uint32 = 100
	if header.NSamples < amountOfSamples {
		log.Fatal("Not enough samples avialable")
	}

	log.Println("Requesting data ...")
	samples, err := conn.GetData(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Visualizing the channels
	channels := make([]channelXYer, header.NChannels)
	for _, sample := range samples {
		for i := uint32(0); i < header.NChannels; i++ {
			channels[i].Values = append(channels[i].Values, sample[i])
		}
	}

	log.Println("Plotting ...")
	plt, err := plot.New()
	plotutil.AddLinePoints(plt,
		"CH0", plotter.XYer(channels[0]),
		"CH1", plotter.XYer(channels[1]),
		"CH2", plotter.XYer(channels[2]))

	log.Println("Saving plot to output.jpg ...")
	plt.Save(10*vg.Inch, 5*vg.Inch, "output.jpg")
}

type channelXYer struct {
	Values []float64
}

func (c channelXYer) Len() int {
	return len(c.Values)
}

func (c channelXYer) XY(index int) (x, y float64) {
	return float64(index), c.Values[index]
}
