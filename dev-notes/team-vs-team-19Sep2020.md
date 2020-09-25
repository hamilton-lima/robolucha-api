# Team vs Team implementation plan 

# Notes from July2020

## API
- Preload teamdefinition on gamedefinition - OK
- change /play to receive team information - OK 
- Add endpoint to join match with teaminformation - OK

## Game 
- allow the user to choose the team - OK 
- when calling /play add team information - OK 
- If only one team dont ask to choose team - OK

- show page to wait while the team is formed
	/wait
- change score page to group results by team

## Runner
- verify min number of team participants to start match 
- verify max number of team participants to join match 
- verify friendly fire 
- send group participation info to match participant or score

# Notes
- Preload teamdefinition - OK

Optimization of play to allow multiple API running
- remove the mutex control - OK

join match with team information

go test -timeout 30s gitlab.com/robolucha/robolucha-api/routes/play -run ^TestLeaveTutorial$

## Removed the indexes, keep the syntax for reference
```
	AvailableMatchID   uint            `gorm:"unique_index:idx_available_match_state" json:"availableMatchID"`
	State              string          `gorm:"unique_index:idx_available_match_state;default:'created'" json:"state"`
```
Will handle concurrency when the match manager is created

## TestLeaveTutorial failing

when another player tries to join a tutorial match from from the same time
the game is shutdown for the original player? - fixed!

```
{"level":"info","matchID":1,"msg":"Play","status":"match found","time":"2020-09-20T13:20:54-04:00"}
{"level":"info","matchID":1,"msg":"Play","status":"Match is an tutorial, end and create again","time":"2020-09-20T13:20:54-04:00"}
{"level":"info","matchID":3,"msg":"Play","status":"Tutorial recreated","time":"2020-09-20T13:20:54-04:00"}
```

# add team to the matchParticipant 
