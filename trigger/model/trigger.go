package mdmp

import (
	"github.com/qingcloudhx/core/data/metadata"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/trigger"
)

/**
* @Author: hexing
* @Date: 19-6-27 上午10:54
 */
var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Factory is a kafka trigger factory
type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

// Trigger is a mdmp trigger
type Trigger struct {
	settings *Settings
	handlers []trigger.Handler
	logger   log.Logger
	//Consumer client.Client //consumer data
}

// Initialize initializes the trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()
	t.logger = ctx.Logger()

	return nil
}

// Start starts the kafka trigger
func (t *Trigger) Start() error {
	return nil
}

// Start starts the kafka trigger
func (t *Trigger) Stop() error {
	return nil
}
