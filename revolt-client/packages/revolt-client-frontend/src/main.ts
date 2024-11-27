import { Client } from './client';
import './style.css';

const client = new Client;
await client.connect("ws://localhost:8080");

const app = document.querySelector('#app');
if (!app) {
    throw new Error('no app element');
}

const button = document.createElement('button');
button.onclick = () => client.createGame();
button.textContent = 'create game';

app.appendChild(button);


const input = document.createElement('input');
input.type = 'text';
app.appendChild(input);

const join = document.createElement('button');
join.onclick = () => client.joinGame(input.value);
join.textContent = 'join game';

document.querySelector('#app')?.appendChild(join);