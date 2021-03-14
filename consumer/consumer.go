package consumer

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/common"
	"fooddlv/common/asyncjob"
	"fooddlv/pubsub"
	"fooddlv/skio"
	"log"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	appCtx         appctx.AppContext
	realtimeEngine skio.RealtimeEngine
}

func NewEngine(appContext appctx.AppContext, realtimeEngine skio.RealtimeEngine) *consumerEngine {
	return &consumerEngine{appCtx: appContext, realtimeEngine: realtimeEngine}
}

func (engine *consumerEngine) Start() error {
	//ps := engine.appCtx.GetPubsub()

	//engine.startSubTopic(common.ChanNoteCreated, asyncjob.NewGroup(
	//	false,
	//	asyncjob.NewJob(SendNotificationAfterCreateNote(engine.appCtx, context.Background(), nil))),
	//)

	engine.startSubTopic(
		common.ChannelNoteCreated,
		true,
		SendNotificationAfterCreateNote(engine.appCtx),
		DeleteImageRecordAfterCreateNote(engine.appCtx),
		EmitRealtimeAfterCreateNote(engine.appCtx, engine.realtimeEngine),
	)
	// Many sub on a topic

	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Channel, isParallel bool, hdls ...consumerJob) error {
	c, _ := engine.appCtx.GetPubsub().Subscribe(context.Background(), topic)

	for _, item := range hdls {
		log.Println("Setup consumer for:", item.Title)
	}

	getJobHandler := func(job *consumerJob, message *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Println("running job for ", job.Title, ". Value: ", message.Data())
			return job.Hld(ctx, message)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHdlArr := make([]asyncjob.Job, len(hdls))

			for i := range hdls {
				//_ = hdls[i].Hld(context.Background(), msg)

				// capture msg & hlds[i], WRONG CODE
				//jobHdl := func(ctx context.Context) error {
				//	log.Println("running job for ", hdls[i].Title, ". Value: ", msg.Data())
				//	return hdls[i].Hld(ctx, msg)
				//}

				jobHdlArr[i] = asyncjob.NewJob(getJobHandler(&hdls[i], msg))
			}

			group := asyncjob.NewGroup(isParallel, jobHdlArr...)

			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}

//
//func Setup(appCtx appctx.AppContext, ctx context.Context) {
//	// setup all consumer / subscriber
//	RunDeleteImageRecordAfterCreateNote(appCtx, context.Background())
//	RunPushNotificationAfterCreateNote(appCtx, context.Background())
//}
