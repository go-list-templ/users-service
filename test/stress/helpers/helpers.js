export function generateUsersData(count) {
    const jsonPath = import.meta.resolve('../fixtures/create.json');
    const base = JSON.parse(open(jsonPath));

    return Array.from({length: count}, (_, i) => ({
        ...base,
        email: base.email.replace('@', `+${i}@`)
    }));
}


export function getRandomItem(arr) {
    return arr[Math.floor(Math.random() * arr.length)];
}