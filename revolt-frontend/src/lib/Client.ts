import { MessageType, type Message, type StateUpdate } from "./types";

export class Client {

    private socket?: WebSocket;
    state?: StateUpdate;

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

    createGame() {
        this.sendMessage({
            type: MessageType.CreateGame
        });
    }

    joinGame(gameId: string) {
        this.sendMessage({
            type: MessageType.JoinGame,
            payload: {
                gameId
            }
        });
    }

    startGame() {
        this.sendMessage({
            type: MessageType.StartGame,
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
        const message = JSON.parse(event.data) as StateUpdate;
        this.state = message;

        console.log('my id:', this.state.clientId);
        console.log(JSON.stringify(this.state, undefined, 2));

    }
}