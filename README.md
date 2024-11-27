# revolt

A card game.

## Rules

Each player has two cards, which grant permission to perform certain actions. However, players can lie and perform any action.

### Default Actions

- Income: gain 1 coin
- Foreign Aid: gain 2 coins
- Revolt (7 coins): force another player to lose a card. Cannot be blocked.

### Card Actions

|Card |Action |Blocks|
|-----------|------------|--------------------|
|Duke       |Tax         |Block Foreign Aid   |
|Assassin   |Assassinate |n/a                 |
|Ambassador |Exchange    |Block Stealing      |
|Captain    |Steal       |Block Stealing      |
|Contessa   |n/a         |Block Assassination |

Turn:

- Action (lying/true) (can have target)
  - 10+ coins = automatic coup
- Challenge

- If lying, leader loses a card
- If true, leader gets a new card

OR. Block

- If action has target, target can block
- Block can be challenged

- Effect
  - If no challenge or block, action succeeds

## State Machine

There are 8 possible states in the Revolt state machine, detailed below. Each state transition function can only be called from the correct initial state.

```txt
                                            [Default]
                                                |
                                          AttemptAction()
                                                |
                    +----------------- [ActionPending] ------------------+
                    |                                |                   |
                Challenge()                     AttemptBlock()           |  
                    |                                |                   |  
                    |                       +---[BlockPending]----+      |
                    |                       |                     |      |  
                    |                    Challenge()              |      |  
                    |                       |                     |      | 
            [LeaderLostChallenge OR PlayerLostChallenge]          |      |   
                                  |                               |      | 
                           ResolveDeath()                       CommitTurn() 
                             |       |                         |     |     |
                    [ActionPending]  |                [PlayerKilled] | [ExchangePending]
                                     |                         |     |     |
                                     |                ResolveDeath() | ResolveExchange()
                                     |                          |    |     |  
                                     +--------------+-----------+----+-----+
                                                    |
                                                [Finished]
```

For example, a simple turn with no challenges or block would look like this:

`[Default]->AttemptAction()->[ActionPending]->CommitTurn()->[Finished]`

## WebSocket API

```json
// Client sends
{"type": "createGame"}
// Server responds.
{"type": "createGame", "payload":{"id": "aaa-bbb-ccc"}}
```
