package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/spf13/pflag"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"github.com/MovieStoreGuy/datagen/pkg/dice"
	"github.com/MovieStoreGuy/datagen/pkg/generic"
	"github.com/MovieStoreGuy/datagen/pkg/randomise"
	"github.com/MovieStoreGuy/datagen/pkg/simulation"
)

var (
	flagDices = pflag.IntSlice("dice", []int{2, 4, 6, 10, 12, 20}, "The selection of dice to use when trying to perfom the simulation")
	flagCount = pflag.Int("simulation-count", 10000, "Defines the number of times the simulation is run")
	flagSaved = pflag.String("excell-named", "dice-simulation.xlsx", "The name of the excell document created once finished")
)

func main() {
	pflag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	var (
		sheets = excelize.NewFile()
		dices  = generic.StreamMap(func(in int) dice.Die {
			return dice.NewDie(in)
		}, (*flagDices)...)
		sampling = map[string]rand.Source{
			"fixed":      randomise.NewFixedSource(ctx),     // Fixed uses the same value each simulation
			"biased":     randomise.NewRepeatingSource(ctx), // biased uses the same set of values
			"randomised": randomise.NewTimeSource(ctx),      // randomised is different simulation
		}
	)

	for k, s := range sampling {
		var (
			log  = log.Named(k)
			rCtx = randomise.WithRandom(ctx, rand.New(s))
		)
		sheets.NewSheet(k)

		if err := sheets.SetCellStr(k, "A1", "#Sides"); err != nil {
			log.Panic("Issue writing dice header to sheets", zap.Error(err))
		}
		for i, d := range dices {
			if err := sheets.SetCellInt(k, ToExcellAxis(i+1, 0), d.Sides()); err != nil {
				log.Panic("Issue writing dice header to sheets", zap.Error(err))
			}
		}

		log.Info("Starting simulation")
		sim := simulation.NewDiceSimulation(dices)

		for i := 0; i < *flagCount; i++ {
			for j, v := range sim.Next(rCtx) {
				if err := sheets.SetCellInt(k, ToExcellAxis(j+1, i+1), v); err != nil {
					log.Error("Isse trying to write cell data", zap.Error(err))
				}
			}
		}

		log.Info("Completed simulation")
	}

	if err := sheets.SaveAs(*flagSaved); err != nil {
		log.Panic("Issue saving excell document", zap.Error(err))
	}
	log.Info("Saved excell document", zap.Stringp("location", flagSaved))
}

// ToExcellAxis is not fit for any use, so this will need to be addressed in future.
// BUG(MovieStoreGuy): This will break if col is above 26 or uses a negative value
func ToExcellAxis(col, row int) string {
	collums := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	return fmt.Sprintf("%s%d", string(collums[col]), row+1)
}
