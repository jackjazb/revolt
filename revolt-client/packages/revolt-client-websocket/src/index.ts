import { Client } from "./client";

export function test(): string {
    return 'Hello!';
}

console.log(test());

const client = new Client;
client.connect();