
export function randomName(): string {
    const adjectives = ['Cool', 'Angry', 'Purple', 'Red', 'Awful', 'Cool'];
    const nouns = ['Buffalo', 'Bean', 'Giraffe', 'Carrot', 'Loaf'];
    const adjective = adjectives[Math.floor(Math.random() * (adjectives.length - 1))];
    const noun = nouns[Math.floor(Math.random() * (nouns.length - 1))];
    return `${adjective}${noun}`;
}