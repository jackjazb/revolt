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
        ownerId: "dc32ad16",
        self: {
            name: "Jack",
            id: "dc32ad16",
            cards: [
                {
                    card: Card.Contessa,
                    alive: true,
                },
                {
                    card: Card.Assassin,
                    alive: true,
                },
            ],
            credits: 7,
            leading: true,
            allowedActions: [
                ActionType.Income,
                ActionType.ForeignAid,
                ActionType.Revolt,
                ActionType.Assassinate,
            ],
        },
        peers: [
            {
                name: "Alice",
                id: "badbe5f2",
                cards: [],
                credits: 7,
                leading: false,
                allowedActions: undefined,
            },
            {
                name: "Bob",
                id: "5717f23c",
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
                id: "6717f23c",
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
                id: "6717f23c",
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
                id: "6717f23c",
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
        turnState: TurnState.ActionPending,
        nextDeath: "dc32ad16",
        winner: "",
        pendingAction: {
            type: ActionType.Assassinate as ActionType,
            target: "6717f23c",
        },
        pendingBlock: {
            card: Card.Contessa as Card,
            initiator: "6717f23c",
        },
        pendingChallenge: {
            initiator: "",
        },
    };

    global.state = state;
</script>

<Game />
