package services

import (
	dfclient "github.com/mlabouardy/dialogflow-go-client"
	dfmodels "github.com/mlabouardy/dialogflow-go-client/models"
	pers "github.com/nicovillanueva/beardbot/persistence"
	log "github.com/sirupsen/logrus"
)

var client *dfclient.DialogFlowClient
var s *pers.Storage

func buildDialogFlowClient(ctx SettingsContext) *dfclient.DialogFlowClient {
	if client == nil {
		var err error
		err, client = dfclient.NewDialogFlowClient(dfmodels.Options{
			AccessToken: GetProvider(ctx).APIKeys.DialogflowAPI,
		})
		if err != nil {
			log.Fatalln("Could not create Dialogflow client")
		}
		log.Info("Created Dialogflow client")
	}
	return client
}

// ReplyToQuery goes out to Dialogflow with a given string, and returns whatever
// speech DF gave back.
// TODO: Actually configure and use fulfillments
func ReplyToQuery(ctx SettingsContext, queryText string) string {
	c := buildDialogFlowClient(ctx)
	query := dfmodels.Query{
		Query: queryText,
	}
	log.Infof("Built query: %+v", query)
	resp, err := c.QueryFindRequest(query)
	if err != nil {
		log.Errorf("Could not query DF: %+v", err)
	}
	log.Infof("DF response: %+v", resp.Result)
	// process response
	if resp.Result.Action == "input.unknown" && resp.Result.ActionIncomplete == false {
		// save unknown action (no context, out of fucking nowhere)
		go func() {
			log.Warnf("Didnt know what to do with: %s", resp.Result.ResolvedQuery)
			s.SaveUnknownResolvedQuery(resp.Result.ResolvedQuery)
		}()
	}
	/*
	   unknown response but in a context'd conversation
	   resp.Result.Action == ""
	   resp.Result.ActionIncomplete == true
	   save resp.Result.Parameters, resp.Result.Contexts, resp.Result.Metadata.IntentName, resp.Result.ResolvedQuery
	*/
	return resp.Result.Fulfillment.Speech
}
