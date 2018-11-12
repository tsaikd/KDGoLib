package slackutil

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/nlopes/slack"
	"github.com/tsaikd/KDGoLib/errutil"
	"golang.org/x/sync/errgroup"
)

// LoggerOptions options for create logger
type LoggerOptions struct {
	// default log channel, default: "#general"
	ChannelDefault string
	// channel for debug, fallback to Default channel if empty
	ChannelDebug string
	// channel for error, fallback to Default channel if empty
	ChannelError string

	// format log
	Formatter Formatter

	// fallback logger when slack log failed, default: errutil.Logger()
	FallbackLogger errutil.LoggerType

	// context for async api call, default: context.Background()
	Context context.Context
	// message queue size for async api call, default: 10
	MessageQueueSize int64
	// worker count for async api call, default: 1
	Worker int
	// do not use async api call
	Sync bool
	// flush timeout before close, default 3 seconds
	FlushTimeout time.Duration
}

func (t *LoggerOptions) check() {
	if t.ChannelDefault == "" {
		t.ChannelDefault = "#general"
	}
	if t.ChannelDebug == "" {
		t.ChannelDebug = t.ChannelDefault
	}
	if t.ChannelError == "" {
		t.ChannelError = t.ChannelDefault
	}
	if t.FallbackLogger == nil {
		t.FallbackLogger = errutil.Logger()
	}
	if t.Context == nil {
		t.Context = context.Background()
	}
	if t.MessageQueueSize < 1 {
		t.MessageQueueSize = 10
	}
	if t.Worker < 1 {
		t.Worker = 1
	}
	if t.FlushTimeout < 1 {
		t.FlushTimeout = 3 * time.Second
	}
}

// NewLogger create new Logger instance
func NewLogger(token string, opt LoggerOptions) (instance *LoggerType, err error) {
	opt.check()
	client := slack.New(token)
	if _, err = client.AuthTest(); err != nil {
		err = ErrSlackAuthTestFailed.New(err)
		return
	}

	ctx, cancel := context.WithCancel(opt.Context)
	eg, ctx := errgroup.WithContext(ctx)
	opt.Context = ctx

	logger := &LoggerType{
		token:        token,
		options:      opt,
		client:       client,
		messageQueue: make(chan messageType, opt.MessageQueueSize),
		eg:           eg,
		cancel:       cancel,
	}

	if !opt.Sync {
		for i := 0; i < opt.Worker; i++ {
			eg.Go(func() error {
				logger.workerLoop()
				return nil
			})
		}
	}

	return logger, nil
}

// LoggerType logger main type
type LoggerType struct {
	token   string
	options LoggerOptions

	client       *slack.Client
	messageQueue chan messageType
	eg           *errgroup.Group
	cancel       context.CancelFunc
	sendingCount int64
}

type messageType struct {
	ChannelID string
	Options   []slack.MsgOption
}

func (t *LoggerType) workerLoop() {
	ctx := t.options.Context
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-t.messageQueue:
			t.sendMessageSync(ctx, msg.ChannelID, msg.Options...)
		}
	}
}

func (t *LoggerType) flush(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(t.options.Context, timeout)
	defer cancel()
	for len(t.messageQueue) > 0 || t.sendingCount > 0 {
		fmt.Println(time.Now().UnixNano())
		select {
		case <-ctx.Done():
			return
		case msg := <-t.messageQueue:
			t.sendMessageSync(ctx, msg.ChannelID, msg.Options...)
		}
	}
}

// Close cancel slack context and wait for worker loop done
func (t *LoggerType) Close() (err error) {
	t.flush(t.options.FlushTimeout)
	t.cancel()
	return t.eg.Wait()
}

// Debug output message to slack debug channel
func (t LoggerType) Debug(v ...interface{}) {
	text := fmt.Sprint(v...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, false)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelDebug, msgopts...)
}

// Debugf output message to slack debug channel
func (t LoggerType) Debugf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, false)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelDebug, msgopts...)
}

// Print output message to slack default channel
func (t LoggerType) Print(v ...interface{}) {
	text := fmt.Sprint(v...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, false)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelDefault, msgopts...)
}

// Printf output message to slack default channel
func (t LoggerType) Printf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, false)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelDefault, msgopts...)
}

// Error output message to slack error channel
func (t LoggerType) Error(v ...interface{}) {
	text := fmt.Sprint(v...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, true)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelError, msgopts...)
}

// Errorf output message to slack error channel
func (t LoggerType) Errorf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)

	msgopts, err := t.options.Formatter.formatSkip(text, 1, true)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, text))
		return
	}

	t.sendMessage(t.options.ChannelError, msgopts...)
}

// Trace output message to slack error channel if errin is not nil
func (t LoggerType) Trace(errin error) {
	if errin == nil {
		return
	}

	msgopts, err := t.options.Formatter.formatSkip(errin.Error(), 1, true)
	if err != nil {
		t.options.FallbackLogger.Trace(ErrSlackFormatFailed1.New(err, errin.Error()))
		t.options.FallbackLogger.Trace(errin)
		return
	}

	t.sendMessage(t.options.ChannelError, msgopts...)
}

func (t LoggerType) sendMessage(channelID string, options ...slack.MsgOption) {
	ctx := t.options.Context

	if t.options.Sync {
		t.sendMessageSync(ctx, channelID, options...)
		return
	}

	t.messageQueue <- messageType{
		ChannelID: channelID,
		Options:   options,
	}
}

func (t *LoggerType) sendMessageSync(ctx context.Context, channelID string, options ...slack.MsgOption) {
	atomic.AddInt64(&t.sendingCount, 1)
	defer atomic.AddInt64(&t.sendingCount, -1)
	_, _, _, err := t.client.SendMessageContext(ctx, channelID, options...)
	t.options.FallbackLogger.Trace(err)
	return
}
