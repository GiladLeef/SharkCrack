//+build wireinject

package main

import (
	"github.com/GiladLeef/SharkCrack/api"
	"github.com/GiladLeef/SharkCrack/apiclient"
	"github.com/GiladLeef/SharkCrack/client"
	"github.com/GiladLeef/SharkCrack/encoder"
	"github.com/GiladLeef/SharkCrack/flusher"
	"github.com/GiladLeef/SharkCrack/logger"
	"github.com/GiladLeef/SharkCrack/observer"
	"github.com/GiladLeef/SharkCrack/queue"
	"github.com/GiladLeef/SharkCrack/reader"
	"github.com/GiladLeef/SharkCrack/requester"
	"github.com/GiladLeef/SharkCrack/server"
	"github.com/GiladLeef/SharkCrack/submitter"
	"github.com/GiladLeef/SharkCrack/tracker"
	"github.com/GiladLeef/SharkCrack/userinput"
	"github.com/GiladLeef/SharkCrack/verifier"
	"github.com/GiladLeef/SharkCrack/waiter"
	"github.com/google/wire"
)

func InitializeClient() client.Client {
	wire.Build(client.NewClient, encoder.NewHasherFactory, requester.NewPasswordRequester, submitter.NewHashSubmitter, apiclient.NewHashApiClient, flusher.NewClientQueueFlusher, reader.NewClientBackupReader, queue.NewHashingRequestQueue, queue.NewHashingSubmissionQueue, waiter.NewSleeper, logger.NewConcurrentLogger, queue.NewClientStopReasonQueue, userinput.NewCmdLineConfigProvider)
	return client.Client{}
}

func InitializeServer() server.Server {
	wire.Build(server.NewServer, api.NewHashApi, verifier.NewHashVerifier, reader.NewHashlistReader, reader.NewWordlistReader, queue.NewPasswordQueue, queue.NewHashQueue, observer.NewStatsObserver, logger.NewConcurrentLogger, tracker.NewStatsTracker, userinput.NewCmdLineConfigProvider)
	return server.Server{}
}
