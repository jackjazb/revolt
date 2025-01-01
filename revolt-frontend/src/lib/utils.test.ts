import { describe, expect, it } from 'vitest';
import { ActionType, Card } from "./types";
import { formatCard, getCurrentActionBlockers } from "./utils";

describe('formatCard', () => {
    it('should format a card in title case', () => {
        expect(formatCard(Card.Duke)).toBe('Duke');
    });
});

describe('getCurrentActionBlockers', () => {
    it('should return that Duke blocks foreign aid', () => {
        const expected = [Card.Duke];
        const result = getCurrentActionBlockers({
            pendingAction: {
                type: ActionType.ForeignAid
            }
        } as any);
        expect(result).toStrictEqual(expected);
    });

    it('should return blockers for targeted actions if the current peer is targeted', () => {
        const expected = [Card.Contessa];
        const result = getCurrentActionBlockers({
            self: {
                id: "1234"
            },
            pendingAction: {
                type: ActionType.Assassinate,
                target: "1234"
            }
        } as any);
        expect(result).toStrictEqual(expected);
    });

    it('should return an empty list for targeted actions if the current peer is not targeted', () => {
        const result = getCurrentActionBlockers({
            self: {
                id: "1234"
            },
            pendingAction: {
                type: ActionType.Assassinate,
                target: "5678"
            }
        } as any);
        expect(result).toStrictEqual([]);
    });
});