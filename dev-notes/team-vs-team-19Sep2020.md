# Team vs Team implementation plan 

# Notes from July2020
## API
- Preload teamdefinition
- Add endpoint to join match with teaminformation

## Runner
- verify min number of team participants to start match 
- verify max number of team participants to join match 

## Game 
- add route to /lobby/:availableMatchID

Logic
- try to check status (what status?)
- start match if not started 
- if has team information and more than one ask to choose team
- then show button that redirects to /watch 

# Lets go!

- Preload teamdefinition - OK

Optimization of play to allow multiple API running
- idx_available_match_state created - OK 

- when finishing the match update state with null
- Try to create match, If the creation fails, get active Match by GameDefinitionID
- If the creation succeed, return the created Match
- send message: match.start : runner is listening
- remove the mutex control 

join match with team information
