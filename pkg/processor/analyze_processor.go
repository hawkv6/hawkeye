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

func (processor *AnalyzeProcessor) ProcessNodes(nodes []domain.Node) {}

func (processor *AnalyzeProcessor) ProcessLinks(links []domain.Link) error {
	processor.normalizer.Normalize(links)
	return nil
}

func (processor *AnalyzeProcessor) ProcessPrefixes(prefixes []domain.Prefix) {}

func (processor *AnalyzeProcessor) ProcessSids(sids []domain.Sid) {}

func (processor *AnalyzeProcessor) Start() {
	processor.log.Infoln("Starting analyze processor")
}

func (processor *AnalyzeProcessor) Stop() {
	processor.log.Infoln("Stopping analyze processor")
}
