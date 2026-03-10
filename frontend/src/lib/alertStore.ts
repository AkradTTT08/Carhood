import { writable } from 'svelte/store';

export interface AlertState {
    show: boolean;
    message: string;
    type: 'error' | 'success' | 'info';
}

const initialState: AlertState = {
    show: false,
    message: '',
    type: 'error'
};

const store = writable<AlertState>(initialState);

export const alertStore = {
    subscribe: store.subscribe,
    showAlert: (message: string, type: 'error' | 'success' | 'info' = 'error') => {
        store.set({ show: true, message, type });
    },
    hideAlert: () => {
        store.set(initialState);
    }
};
