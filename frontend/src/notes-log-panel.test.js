import { describe, expect, it } from 'vitest';

import { initNotesLogPanel } from './notes-log-panel.js';

function createElements() {
    return {
        appRoot: document.createElement('div'),
        logOutput: document.createElement('div'),
        toolsLogWordwrap: document.createElement('button'),
        toolsLogTimestamp: document.createElement('button'),
        toolsLogMaximize: document.createElement('button'),
    };
}

function getLineTexts(logOutput) {
    return Array.from(logOutput.querySelectorAll('.notes-log-line')).map((line) => line.textContent);
}

describe('notes log panel', () => {
    it('caps rendered log output at max log lines', () => {
        const elements = createElements();
        const handlers = {};
        const eventsOn = (name, handler) => {
            handlers[name] = handler;
        };

        initNotesLogPanel(elements, eventsOn, 3);

        handlers.notesLog('one\ntwo\nthree\nfour\n');

        expect(elements.logOutput.childElementCount).toBe(3);
        expect(getLineTexts(elements.logOutput)).toEqual(['two', 'three', 'four']);
    });

    it('applies updated max line limits immediately', () => {
        const elements = createElements();
        const handlers = {};
        const eventsOn = (name, handler) => {
            handlers[name] = handler;
        };

        const panel = initNotesLogPanel(elements, eventsOn, 10);

        handlers.notesLog('alpha\nbeta\ngamma\ndelta\n');
        panel.setMaxLogLines(2);

        expect(elements.logOutput.childElementCount).toBe(2);
        expect(getLineTexts(elements.logOutput)).toEqual(['gamma', 'delta']);
    });

    it('keeps the Font Awesome maximize glyph while toggling the log panel', () => {
        const elements = createElements();
        const handlers = {};
        const eventsOn = (name, handler) => {
            handlers[name] = handler;
        };

        initNotesLogPanel(elements, eventsOn, 10);
        expect(elements.toolsLogMaximize.textContent).toBe(String.fromCodePoint(0xf065));

        elements.toolsLogMaximize.click();
        expect(elements.toolsLogMaximize.dataset.enabled).toBe('true');
        expect(elements.toolsLogMaximize.textContent).toBe(String.fromCodePoint(0xf065));

        elements.toolsLogMaximize.click();
        expect(elements.toolsLogMaximize.dataset.enabled).toBe('false');
        expect(elements.toolsLogMaximize.textContent).toBe(String.fromCodePoint(0xf065));
    });
});
