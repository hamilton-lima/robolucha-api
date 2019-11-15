select matches.id, time_end, time_start, duration, 
UNIX_TIMESTAMP(time_start) + duration, UNIX_TIMESTAMP(NOW())

from matches, game_definitions
where game_definition_id = game_definitions.id
and time_end < time_start 
and UNIX_TIMESTAMP(time_start) + (duration/1000) > UNIX_TIMESTAMP(NOW());
