# Request

After implementing the single-process WebSocket close policy primitive, stop at a new explicit maintainer confirmation point.

The next work must not implicitly add protocol logout routes, reconnect/epoch behavior, protocol session carriers, concrete transport close handoff, operations/admin disconnect, direct Nakama/Pitaya API compatibility, or broader game backend behavior.

The maintainer then clarified that vibit should become a Nakama/Pitaya-class product and cover their common feature surface. The selected direction is `expand_core_game_backend_modules_after_nakama_pitaya_review`, scoped first to a checkable product parity roadmap before broad module implementation.
