import { writable } from 'svelte/store';

export type GameState = 'LOBBY' | 'QUESTION' | 'RESULT' | 'FINISHED';

export interface Question {
    text: string;
    options: string[];
    correct_answer: number;
    image_url?: string;
    time_limit?: number;
}

export const socketStore = writable<WebSocket | null>(null);
export const playersStore = writable<string[]>([]);
export const currentUser = writable<string>("");
export const gameStatus = writable<GameState>('LOBBY');
export const currentQuestion = writable<Question | null>(null);
export const answersStore = writable<{ username: string, answer_index: number } | null>(null);

let socket: WebSocket | null = null;

export function connect(roomID: string, username: string) {
    currentUser.set(username);
    socket = new WebSocket('ws://127.0.0.1:8081/ws');

    socket.onopen = () => {
        if (socket) {
            socket.send(JSON.stringify({
                type: 'join',
                payload: { room_id: roomID, username: username }
            }));
            socketStore.set(socket);

            // Pre-populate with self if not HOST
            if (username !== 'HOST') {
                playersStore.set([username]);
            }
        }
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        switch (msg.type) {
            case 'sync_players':
                playersStore.set(msg.payload.players.filter((p: string) => p !== 'HOST'));
                break;
            case 'player_joined':
                if (msg.payload.username === 'HOST') break;
                playersStore.update(p => {
                    if (p.includes(msg.payload.username)) return p;
                    return [...p, msg.payload.username];
                });
                break;
            case 'start_game':
            case 'next_question':
                gameStatus.set('QUESTION');
                currentQuestion.set(msg.payload.question);
                break;
            case 'show_results':
                gameStatus.set('RESULT');
                break;
            case 'finish_game':
                gameStatus.set('FINISHED');
                break;
            case 'player_left':
                playersStore.update(players => players.filter(p => p !== msg.payload.username));
                break;
            case 'submit_answer':
                answersStore.set(msg.payload);
                break;
            case 'cancel_game':
                gameStatus.set('LOBBY'); // Or handle redirect in component
                window.location.href = '/';
                break;
            case 'error':
                alert(msg.payload.message);
                break;
        }
    };

    socket.onclose = () => {
        socketStore.set(null);
        socket = null;
    };
}

export function sendMessage(type: string, payload: any) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type, payload }));
    }
}
