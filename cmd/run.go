/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/mauroalderete/weasel/pathfinder"
	"github.com/mauroalderete/weasel/wallet"
	"github.com/mauroalderete/weasel/wallet/repository"
	"github.com/mauroalderete/weasel/wallet/repository/filerepository"
	"github.com/mauroalderete/weasel/wallet/repository/stdoutrepository"
	"github.com/spf13/cobra"
)

const (
	ROJO  = 31
	VERDE = 32
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes weasel to explore accounts",
	Long: `Start weasel service to generate random accounts and explore his activity.
Store accounts with activity.
For example:

weasel run --thread 12

Weasel is a tool to search accounts with activity generating a private key randomly.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Int32P("thread", "t", 1, "Number of threads >0 to execute. Each thread handle his own connection and own search.")
	runCmd.Flags().StringP("gateway", "g", "https://cloudflare-eth.com", "Ethereum gateway to connect.")
	runCmd.Flags().BoolP("stop-search-errors", "e", false, "Stop all process when an error in is detected in any thread.")
	runCmd.Flags().StringP("match-file", "m", "", "Filepath to store in json format all wallets matched.")
	runCmd.Flags().StringP("unmatch-file", "u", "", "Filepath to store in json format all wallets unmatched.")

	runCmd.RunE = runMain
}

func runMain(cmd *cobra.Command, args []string) error {

	log.Printf("Verifing arguments...")
	// reviso los argumentos
	v, err := strconv.ParseInt(cmd.Flag("thread").Value.String(), 10, 32)
	if err != nil {
		log.Printf("[FAIL]")
		return fmt.Errorf("thread argument must be a integer value, currently is %s", cmd.Flag("thread").Value.String())
	}
	threads := int(v)
	gateway := cmd.Flag("gateway").Value.String()
	stopSearchErrors := cmd.Flag("stop-search-errors").Value.String() == "true"
	matchFilename := cmd.Flag("match-file").Value.String()
	unmatchFilename := cmd.Flag("unmatch-file").Value.String()

	// preparo los repositorios
	log.Printf("Preparing repositories...")
	repoHandler := RepositoryHandle{}
	err = repoHandler.Start(matchFilename, unmatchFilename)
	if err != nil {
		log.Printf("[FAIL]")
		return fmt.Errorf("error to prepare repositories: %v", err)
	}

	// preparo los canales
	log.Printf("Preparing channels...")
	termsignal := make(chan os.Signal, 1)
	signal.Notify(termsignal, syscall.SIGINT, syscall.SIGTERM)

	stopsignal := make(chan bool, 1)
	errorsignal := make(chan error, 1)

	var wg sync.WaitGroup

	log.Printf("Starting threads...\n")
	// inicializo los buscadores
	for i := 0; i < threads; i++ {
		log.Printf("\tLaunching thread [%d]...\n", i)
		wg.Add(1)
		go LaunchSearcher(i, gateway, &repoHandler, stopsignal, errorsignal, &wg)
	}

	log.Printf("Ready!")
	// espero hasta que finalice
	var someerror error
mainLoop:
	for {
		select {
		case sig := <-termsignal:
			{
				log.Printf("Signal received: %v\n", sig)
				break mainLoop
			}
		case err := <-errorsignal:
			{
				log.Printf("Something was wrong: %v\n", err)
				if stopSearchErrors {
					someerror = err
					break mainLoop
				}
			}
		}
	}

	log.Printf("Stopping threads...\n")
	for i := 0; i < threads; i++ {
		stopsignal <- true
	}

	log.Printf("Awaiting just all terminate correctly...\n")
	wg.Wait()

	log.Printf("Close repositories...\n")
	repoHandler.Close()

	log.Printf("Exit\n")
	return someerror
}

func LaunchSearcher(idx int, gateway string, repoHandler *RepositoryHandle, stopsignal chan bool, errsignal chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	pf := pathfinder.Pathfinder{}

	err := pf.Connect(gateway)
	if err != nil {
		errsignal <- fmt.Errorf("failed to connect to gateway '%s': %v", gateway, err)
	}
	defer pf.Close()

loop:
	for {
		select {
		case <-stopsignal:
			{
				break loop
			}
		default:
			err := pf.Search()
			if err != nil {
				errsignal <- fmt.Errorf("error in %d thread: %v", idx, err)
			} else {

				err := repoHandler.Save(*pf.Wallet(), pf.Match())
				if err != nil {
					errsignal <- fmt.Errorf("error saving match wallet in %d thread: %v", idx, err)
				}
			}
		}
	}

	log.Printf("thread [%d]: Done\n", idx)
}

type RepositoryHandle struct {
	matchs   []repository.Repository
	unmatchs []repository.Repository
}

func (r *RepositoryHandle) Save(w wallet.Wallet, match bool) error {
	if match {
		for _, v := range r.matchs {
			err := v.Save(w)
			if err != nil {
				return fmt.Errorf("error saving the wallet in match stores: %v", err)
			}
		}
	} else {
		for _, v := range r.unmatchs {
			err := v.Save(w)
			if err != nil {
				return fmt.Errorf("error saving the wallet in match stores: %v", err)
			}
		}
	}

	return nil
}

func (r *RepositoryHandle) Close() {
	for _, v := range r.matchs {
		v.Close()
	}

	for _, v := range r.unmatchs {
		v.Close()
	}
}

func (r *RepositoryHandle) Start(matchFilename string, unmatchFilename string) error {
	//prepare pools
	r.matchs = make([]repository.Repository, 0)
	r.unmatchs = make([]repository.Repository, 0)

	// salida estandar para match
	stdoutMatch, err := stdoutrepository.New(VERDE)
	if err != nil {
		return fmt.Errorf("failed to instance a FileRepository to store match wallets: %v", err)
	}

	stdoutUnmatch, err := stdoutrepository.New(ROJO)
	if err != nil {
		return fmt.Errorf("failed to instance a FileRepository to store match wallets: %v", err)
	}

	r.matchs = append(r.matchs, stdoutMatch)
	r.unmatchs = append(r.unmatchs, stdoutUnmatch)

	// json file to match wallets
	if len(matchFilename) != 0 {
		fmatch, err := filerepository.New(matchFilename)
		if err != nil {
			return fmt.Errorf("failed to instance a FileRepository to store match wallets in %s: %v", matchFilename, err)
		}

		r.matchs = append(r.matchs, fmatch)
	}

	// json file to unmatch wallets
	if len(unmatchFilename) != 0 {
		funmatch, err := filerepository.New(unmatchFilename)
		if err != nil {
			return fmt.Errorf("failed to instance a FileRepository to store unmatch wallets in %s: %v", unmatchFilename, err)
		}
		r.unmatchs = append(r.unmatchs, funmatch)
	}

	return nil
}
