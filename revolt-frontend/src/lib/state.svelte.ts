import { Client } from "./Client";
import { initialState, type State } from "./types";

/**
 * Utility class allowing the entire app to react to state update messages - the `global` object
 * should be used to access its members.
 * 
 * Components accessing `global.state` will rerender when the server sends a `ClientStateBroadcast`.
 */
class GlobalStore {
    state: State = $state(initialState);
    // The callback here is used to trigger rerenders when the game state updates.
    client = $state(new Client((update) => (this.state = update)));
}

export const global = new GlobalStore();