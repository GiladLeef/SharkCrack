// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

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
)

// Injectors from wire.go:

func InitializeClient() client.Client {
	configProvider := userinput.NewCmdLineConfigProvider()
	requestQueue := queue.NewHashingRequestQueue()
	submissionQueue := queue.NewHashingSubmissionQueue()
	backupReader := reader.NewClientBackupReader(configProvider, requestQueue, submissionQueue)
	interfacesLogger := logger.NewConcurrentLogger(configProvider)
	clientStopQueue := queue.NewClientStopReasonQueue(configProvider)
	interfacesWaiter := waiter.NewSleeper(configProvider, interfacesLogger)
	encoderFactory := encoder.NewHasherFactory(configProvider, interfacesLogger, requestQueue, submissionQueue, clientStopQueue, interfacesWaiter)
	apiClient := apiclient.NewHashApiClient(configProvider)
	interfacesRequester := requester.NewPasswordRequester(configProvider, apiClient, interfacesLogger, requestQueue, clientStopQueue, interfacesWaiter)
	interfacesSubmitter := submitter.NewHashSubmitter(configProvider, apiClient, interfacesLogger, submissionQueue, clientStopQueue, interfacesWaiter)
	interfacesFlusher := flusher.NewClientQueueFlusher(configProvider, requestQueue, submissionQueue)
	clientClient := client.NewClient(backupReader, configProvider, encoderFactory, interfacesLogger, interfacesRequester, interfacesSubmitter, interfacesFlusher)
	return clientClient
}

func InitializeServer() server.Server {
	configProvider := userinput.NewCmdLineConfigProvider()
	flushingQueue := queue.NewHashQueue(configProvider)
	interfacesQueue := queue.NewPasswordQueue(configProvider)
	interfacesLogger := logger.NewConcurrentLogger(configProvider)
	interfacesTracker := tracker.NewStatsTracker()
	interfacesApi := api.NewHashApi(configProvider, flushingQueue, interfacesQueue, interfacesLogger, interfacesTracker)
	passwordReader := reader.NewWordlistReader(configProvider, interfacesQueue)
	interfacesObserver := observer.NewStatsObserver(interfacesLogger, interfacesTracker, configProvider)
	hashReader := reader.NewHashlistReader(configProvider)
	interfacesVerifier := verifier.NewHashVerifier(flushingQueue, hashReader, interfacesLogger, interfacesTracker)
	serverServer := server.NewServer(interfacesApi, interfacesLogger, passwordReader, interfacesObserver, interfacesVerifier)
	return serverServer
}
