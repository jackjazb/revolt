import { describe, expect, test } from "vitest";
import { Client } from "./client";

describe('Client', () => {
    test('should be a client', () => {
        const client = new Client;
        client.connect();
        expect(client).toBeDefined();
    });
});