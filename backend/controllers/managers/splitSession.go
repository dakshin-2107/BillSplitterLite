package managers

// run a separate go routine that will publish all actions that come through a broadcast channel
func (s *SplitSession) RunPublisherService() {

}

// run a separate go routine that will tally all shares of the takers and publish the result to everyone at periodic intervals say every 5 seconds
// get the bill from the database and calculate the shares of the takers apply the actions
func (s *SplitSession) RunTallyService() {

}
