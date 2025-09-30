package managers

func (s *SessionManager) QueueAction(action IAction) error {
	return nil
}

func (s *SessionManager) RunActionService() {

	for {
		// time.Sleep(100)
		currAction, err := s.DequeueAction() // Needs to pick up an action from redis message queue in redis
		if err != nil {
			continue
		}

		actionErr := currAction.ExecuteAction() // Execute the action on the copy in the redis database
		if actionErr != nil {
			continue
		}

		s.PublishAction(currAction) // Publish that action to all clients part of that bill using connection manager

	}
}

func (s *SessionManager) DequeueAction() (IAction, error) {

	return nil, nil
}

func (s *SessionManager) PublishAction(action IAction) error {

	return nil
}

func (a *Action) ExecuteAction() error {
	return nil
}
