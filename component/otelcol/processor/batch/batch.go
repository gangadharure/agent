// Package batch provides an otelcol.processor.batch component.
package batch

import (
	"fmt"
	"time"

	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/otelcol"
	"github.com/grafana/agent/component/otelcol/processor"
	otelcomponent "go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

func init() {
	component.Register(component.Registration{
		Name:    "otelcol.processor.batch",
		Args:    Arguments{},
		Exports: otelcol.ConsumerExports{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			fact := batchprocessor.NewFactory()
			return processor.New(opts, fact, args.(Arguments))
		},
	})
}

// Arguments configures the otelcol.processor.batch component.
type Arguments struct {
	Timeout          time.Duration `river:"timeout,attr,optional"`
	SendBatchSize    uint32        `river:"send_batch_size,attr,optional"`
	SendBatchMaxSize uint32        `river:"send_batch_max_size,attr,optional"`

	// Output configures where to send processed data. Required.
	Output *otelcol.ConsumerArguments `river:"output,block"`
}

var (
	_ processor.Arguments = Arguments{}
)

// DefaultArguments holds default settings for Arguments.
var DefaultArguments = Arguments{
	Timeout:       200 * time.Millisecond,
	SendBatchSize: 8192,
}

// SetToDefault implements river.Defaulter.
func (args *Arguments) SetToDefault() {
	*args = DefaultArguments
}

// Validate implements river.Validator.
func (args *Arguments) Validate() error {
	if args.SendBatchMaxSize > 0 && args.SendBatchMaxSize < args.SendBatchSize {
		return fmt.Errorf("send_batch_max_size must be greater or equal to send_batch_size when not 0")
	}
	return nil
}

// Convert implements processor.Arguments.
func (args Arguments) Convert() (otelconfig.Processor, error) {
	return &batchprocessor.Config{
		ProcessorSettings: otelconfig.NewProcessorSettings(otelconfig.NewComponentID("batch")),
		Timeout:           args.Timeout,
		SendBatchSize:     args.SendBatchSize,
		SendBatchMaxSize:  args.SendBatchMaxSize,
	}, nil
}

// Extensions implements processor.Arguments.
func (args Arguments) Extensions() map[otelconfig.ComponentID]otelcomponent.Extension {
	return nil
}

// Exporters implements processor.Arguments.
func (args Arguments) Exporters() map[otelconfig.DataType]map[otelconfig.ComponentID]otelcomponent.Exporter {
	return nil
}

// NextConsumers implements processor.Arguments.
func (args Arguments) NextConsumers() *otelcol.ConsumerArguments {
	return args.Output
}
