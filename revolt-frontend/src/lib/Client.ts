import { type Action, type State } from "./types";

export enum MessageType {
    CreateGame = 'create_game',
    JoinGame = 'join_game',
    StartGame = 'start_game',
    AttemptAction = 'attempt_action',
    CommitTurn = 'commit_turn',
    EndTurn = 'end_turn'
}

export interface Message {
    type: MessageType,
    payload?: Record<string, any>;
}

export class Client {

    private socket?: WebSocket;
    private onStateUpdate: (state: State) => void;
    state?: State;

    constructor(onStateUpdate: (state: State) => void) {
        this.onStateUpdate = onStateUpdate;
    }

    async connect(uri: string) {
        const socket = new WebSocket(uri);
        return new Promise<void>(resolve => {
            socket.onopen = () => {
                socket.onmessage = e => this.handleMessage(e);
                this.socket = socket;
                resolve();
            };
        });
    }

    createGame(playerName: string) {
        this.sendMessage({
            type: MessageType.CreateGame,
            payload: {
                playerName
            }
        });
    }

    joinGame(gameId: string, playerName: string) {
        this.sendMessage({
            type: MessageType.JoinGame,
            payload: {
                gameId,
                playerName
            }
        });
    }

    startGame() {
        this.sendMessage({
            type: MessageType.StartGame,
        });
    }

    attemptAction(action: Action) {
        this.sendMessage({
            type: MessageType.AttemptAction,
            payload: {
                action
            }
        });
    }

    commitTurn() {
        this.sendMessage({
            type: MessageType.CommitTurn,
        });
    }

    endTurn() {
        this.sendMessage({
            type: MessageType.EndTurn
        });
    }

    private sendMessage(message: Message) {
        if (!this.socket) {
            return;
        }
        this.socket.send(JSON.stringify(message));
    }

    private handleMessage(event: MessageEvent) {
        if (!event.data) {
            return;
        }
        const message = JSON.parse(event.data) as State;
        this.state = message;
        this.onStateUpdate(this.state);
        console.log('received state update:', JSON.stringify(this.state, undefined, 2));
    }
}