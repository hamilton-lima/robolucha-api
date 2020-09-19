- call get avaialable match

If there is no team to select
- call play/gamedefinitionID

If there is a team to select 
- call play/gamedefinitionID/TeamID

## play (case 1) 
There is no active match for the gamedefinitionID

- Create match
- send message: match.start : runner is listening
- Overcome potential read/update issue with Match unique index of AvailableMatchID + TimeEnd to keep only one single active match for the avaialableMatchID as the TimeEnd is updated when the match ends.
- If the creation fails, get active Match by GameDefinitionID
- If the creation succeed, return the created Match

## play (case 2) 
There is active match for the gamedefinitionID

- get active Match by GameDefinitionID

## play after both cases
- If limit of participants or limit in a team is reached, send error message
(RISK) Will not check for concurrency here

- send message for the match.join : runner is listening
- get the current state of the match for participants/teams
- return current state to the game

# game
If the number of participants is enough to start the match
redirect to /watch/matchID

If still waiting for the necessary participants
- display current participants and how many are missing

If need to choose a team
- display current participants and how many are missing
(NO OPTION TO CHANGE TEAM)

Display countdown on the screen for refresh Match State
Refresh each 5 seconds.
if match is ready, redirect to the watch/matchID

# runner 
(NEXT STEP) when getting start.match, Call API.runMatch/SERVER_ID, if success continue to process the match

- Add check for the minimum number of participants in the team 
to start the match
- Add Friendly fire support based on the game definition

# api runMatch
(NEXT STEP)
