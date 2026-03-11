import { writable } from 'svelte/store';
import { alertStore } from './alertStore';
import { WS_URL } from './api';

export type GameState = 'IDLE' | 'LOBBY' | 'QUESTION' | 'RESULT' | 'FINISHED';

export interface Question {
    text: string;
    options: string[];
    correct_answer: number;
    image_url?: string;
    time_limit?: number;
    points?: number;
}

export const socketStore = writable<WebSocket | null>(null);
export const playersStore = writable<string[]>([]);
export const currentUser = writable<string>("");
export const gameStatus = writable<GameState>('IDLE');
export const currentQuestion = writable<Question | null>(null);
export const answersStore = writable<{ username: string, answer_index: number } | null>(null);
export const allAnswersStore = writable<Record<string, number>>({});
export const questionProgressStore = writable<{ current: number, total: number }>({ current: 0, total: 0 });
export const leaderboardStore = writable<{username: string, score: number}[]>([]);

let socket: WebSocket | null = null;
let connectedUser: string = "";

export function connect(roomID: string, username: string) {
    // Reset all stores to initial state
    playersStore.set([]);
    gameStatus.set('IDLE');
    currentQuestion.set(null);
    answersStore.set(null);
    allAnswersStore.set({});
    questionProgressStore.set({ current: 0, total: 0 });
    leaderboardStore.set([]);

    currentUser.set(username);
    connectedUser = username;
    socket = new WebSocket(WS_URL);

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
                gameStatus.set('LOBBY');
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
                if (msg.payload.progress) {
                    questionProgressStore.set(msg.payload.progress);
                }
                allAnswersStore.set({}); // Reset answers for the new round
                answersStore.set(null);
                break;
            case 'show_results':
                gameStatus.set('RESULT');
                break;
            case 'finish_game':
                if (msg.payload.leaderboard) {
                    leaderboardStore.set(msg.payload.leaderboard);
                }
                gameStatus.set('FINISHED');
                break;
            case 'player_left':
                playersStore.update(players => players.filter(p => p !== msg.payload.username));
                break;
            case 'submit_answer':
                answersStore.set(msg.payload);
                allAnswersStore.update(answers => {
                    answers[msg.payload.username] = msg.payload.answer_index;
                    return answers;
                });
                break;
            case 'cancel_game':
                gameStatus.set('LOBBY');
                // Only redirect players to home, HOST handles their own navigation
                if (connectedUser !== 'HOST') {
                    window.location.href = '/';
                }
                break;
            case 'error':
                console.log('Backend Error Received:', msg.payload.message);
                alertStore.showAlert(msg.payload.message);
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

export function disconnect() {
    if (socket) {
        socket.close();
        socket = null;
        socketStore.set(null);
        playersStore.set([]);
        gameStatus.set('LOBBY');
        currentQuestion.set(null);
        answersStore.set(null);
        allAnswersStore.set({});
        questionProgressStore.set({ current: 0, total: 0 });
        connectedUser = "";
    }
}
