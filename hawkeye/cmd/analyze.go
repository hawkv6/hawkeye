package cmd

import (
	"os"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/analyze"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/jagw"
	"github.com/hawkv6/hawkeye/pkg/normalizer"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/spf13/cobra"
)

type NormalizerCall struct {
	normalizer.Normalizer
	folderName string
	plotTitle  string
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Retrieve and analyze the network metric data",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.NewBaseConfig(jagwServiceAddress, jagwRequestPort)
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		log.Infoln("Config created successfully")

		adapter := adapter.NewDomainAdapter()
		defaultHelper := helper.NewDefaultHelper()

		latencyQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)
		jitterQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)
		packetLossQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)

		normalizers := []NormalizerCall{
			{
				normalizer.NewRobustNormalizer(),
				"robust",
				"Robust Normalization",
			},
			{
				normalizer.NewIQRMinMaxNormalizer(latencyQueue, jitterQueue, packetLossQueue),
				"iqr-based-minmax",
				"IQR Min Max Normalization",
			},
			{
				normalizer.NewMinMaxNormalizer(),
				"minmax",
				"Min Max Normalization",
			},
			{
				normalizer.NewZScoreNormalizer(),
				"zscore",
				"Z-Score Normalization",
			},
		}

		for _, normalizer := range normalizers {
			processor := processor.NewAnalyzeProcessor(normalizer.Normalizer)

			requestService := jagw.NewJagwRequestService(config, adapter, processor, defaultHelper)
			if err := requestService.Init(); err != nil {
				log.Fatalf("Error initializing JAGW Request Service: %v", err)
			}
			if err := requestService.Start(); err != nil {
				log.Fatalf("Error starting JAGW Request Service: %v", err)
			}
			analyzer := analyze.NewMetricAnalyzer(normalizer, normalizer.folderName, normalizer.plotTitle)
			analyzer.Analyze()
			processor.Stop()
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&jagwServiceAddress, "jagw-service-address", "j", os.Getenv("HAWKEYE_JAGW_SERVICE_ADDRESS"), "JAGW Service Address e.g. localhost or 127.0.0.1")
	analyzeCmd.Flags().StringVarP(&jagwRequestPort, "jagw-request-port", "r", os.Getenv("HAWKEYE_JAGW_REQUEST_PORT"), "JAGW Request Port e.g. 9903")
	markRequiredFlags(analyzeCmd, []string{"jagw-service-address", "jagw-request-port"})
}
