<script lang="ts">
    import { global } from "../../state.svelte";
    import {
        ActionType,
        Card,
        GameStatus,
        TurnState,
        type State,
    } from "../../types";
    import Game from "./Game.svelte";

    const state: State = {
        timestamp: "2024-12-28T12:44:13.917749381Z",
        gameId: "6336bae3",
        ownerId: "1",
        self: {
            name: "Jack",
            id: "1",
            cards: [
                {
                    card: Card.Duke,
                    alive: true,
                },
                {
                    card: Card.Assassin,
                    alive: true,
                },
            ],
            credits: 7,
            leading: false,
            allowedActions: [
                ActionType.Income,
                ActionType.ForeignAid,
                ActionType.Revolt,
                ActionType.Assassinate,
                ActionType.Tax,
            ],
        },
        peers: [
            {
                name: "Alice",
                id: "a",
                cards: [],
                credits: 7,
                leading: false,
                allowedActions: undefined,
            },
            {
                name: "Bob",
                id: "b",
                cards: [
                    {
                        card: Card.Duke,
                        alive: false,
                    },
                ],
                credits: 10,
                leading: false,
                allowedActions: undefined,
            },
            {
                name: "Charlene",
                id: "c",
                cards: [
                    {
                        card: Card.Ambassador,
                        alive: false,
                    },
                ],
                credits: 2,
                leading: false,
                allowedActions: undefined,
            },
            {
                name: "Daniel",
                id: "d",
                cards: [
                    {
                        card: Card.Contessa,
                        alive: false,
                    },
                    {
                        card: Card.Duke,
                        alive: false,
                    },
                ],
                credits: 2,
                leading: false,
                allowedActions: undefined,
            },
            {
                name: "Emily",
                id: "e",
                cards: [
                    {
                        card: Card.Assassin,
                        alive: false,
                    },
                ],
                credits: 2,
                leading: false,
                allowedActions: undefined,
            },
        ],
        status: GameStatus.InProgress,
        nextDeath: "e",
        winner: "",
        pendingAction: {
            type: ActionType.Assassinate as ActionType,
            target: "b",
        },
        pendingBlock: {
            card: Card.Contessa as Card,
            initiator: "1",
        },
        pendingChallenge: {
            initiator: "",
        },
        turnState: TurnState.LeaderLostChallenge,
    };

    // Shortcut for setting the leader's ID and other current action initiators.
    let leader = "a";
    let target = "b";
    let blocker = "c";
    let nextDeath = "1";
    let turnState = TurnState.ActionPending;

    if (state.self.id === leader) {
        state.self.leading = true;
    }
    for (const p of state.peers) {
        if (p.id === leader) {
            console.log("hi");
            p.leading = true;
        }
    }
    state.pendingBlock.initiator = blocker;
    state.pendingAction.target = target;
    state.nextDeath = nextDeath;
    state.turnState = turnState;

    global.state = state;
</script>

<Game />
