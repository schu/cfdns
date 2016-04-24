package cmd

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

var (
	recordTTL   int
	recordName  string
	recordType  string
	recordValue string
	recordZone  string
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Create, read, update and delete records",
}

var recordUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a record",
	Run:   recordUpdate,
}

func recordUpdate(cmd *cobra.Command, args []string) {
	if recordZone == "" || recordName == "" || recordValue == "" {
		fmt.Fprint(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}
	record := cloudflare.DNSRecord{
		Name: recordName,
		Type: recordType,
	}
	records, err := api.DNSRecords(recordZone, record)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered error while fetching records: %s\n", err)
		os.Exit(1)
	}
	if len(records) == 0 {
		fmt.Fprintf(os.Stderr, "No such record: %s (%s)\n", recordName, recordType)
		os.Exit(1)
	}
	if len(records) > 1 {
		fmt.Fprintf(os.Stderr, "Record name must be unique: %s\n", recordName)
		os.Exit(1)
	}
	record = records[0]
	record.Content = recordValue
	record.TTL = recordTTL
	err = api.UpdateDNSRecord(recordZone, record.ID, record)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update record: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(recordUpdateCmd)

	recordUpdateCmd.Flags().StringVarP(&recordZone, "zone", "z", "", "zone (required)")
	recordUpdateCmd.Flags().StringVarP(&recordName, "name", "n", "", "record name (required)")
	recordUpdateCmd.Flags().StringVarP(&recordType, "type", "t", "A", "record type")
	recordUpdateCmd.Flags().StringVarP(&recordValue, "value", "v", "", "record value (required)")
	recordUpdateCmd.Flags().IntVar(&recordTTL, "ttl", 120, "ttl in seconds")
}
