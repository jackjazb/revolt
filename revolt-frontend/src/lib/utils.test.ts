import { describe, expect, it } from 'vitest';
import { Card } from "./types";
import { formatCard } from "./utils";

describe('formatCard', () => {
    it('should format a card in title case', () => {
        expect(formatCard(Card.Duke)).toBe('Duke');
    });
});