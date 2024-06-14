package processor

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/normalizer"
	"github.com/sirupsen/logrus"
)

type AnalyzeProcessor struct {
	log        *logrus.Entry
	normalizer normalizer.Normalizer
}

func NewAnalyzeProcessor(normalizer normalizer.Normalizer) *AnalyzeProcessor {
	return &AnalyzeProcessor{
		log:        logging.DefaultLogger.WithField("subsystem", Subsystem),
		normalizer: normalizer,
	}
}

func (processor *AnalyzeProcessor) CreateGraphNodes(nodes []domain.Node) error { return nil }

func (processor *AnalyzeProcessor) CreateGraphEdges(links []domain.Link) error {
	processor.normalizer.Normalize(links)
	return nil
}

func (processor *AnalyzeProcessor) CreateClientNetworks(prefixes []domain.Prefix) {}

func (processor *AnalyzeProcessor) CreateSids(sids []domain.Sid) {}

func (processor *AnalyzeProcessor) Start() {
	processor.log.Infoln("Starting analyze processor")
}

func (processor *AnalyzeProcessor) Stop() {
	processor.log.Infoln("Stopping analyze processor")
}
