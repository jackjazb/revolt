import { mount } from 'svelte';
import App from './App.svelte';
import './main.css';

const app = document.getElementById('app');

if (!app) {
    throw new Error('missing "app" element');
}

export default mount(App, {
    target: app,
});;
