import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readFileSync } from 'node:fs';

const inputboxCss = readFileSync('./src/inputbox.css', 'utf8');

const terminalInputBoxSubmitMock = vi.fn(() => Promise.resolve());
const terminalTextInputMock = vi.fn(() => Promise.resolve());
const terminalKeyPressMock = vi.fn(() => Promise.resolve());
const eventsOnMock = vi.fn();
const showLocalMenuMock = vi.fn();

vi.mock('../wailsjs/go/main/WApp', () => ({
    TerminalInputBoxSubmit: terminalInputBoxSubmitMock,
    TerminalTextInput: terminalTextInputMock,
    TerminalKeyPress: terminalKeyPressMock,
}));

vi.mock('../wailsjs/runtime/runtime', () => ({
    EventsOn: eventsOnMock,
}));

vi.mock('./popup_menu', () => ({
    showLocalMenu: showLocalMenuMock,
}));

vi.mock('./vim-mode', () => ({
    attachVimMode: vi.fn(() => ({ detach: vi.fn(), getMode: vi.fn(() => 'insert') })),
}));

vi.mock('./spellcheck', () => ({
    attachSpellCheck: vi.fn(() => ({ detach: vi.fn(), setExclusions: vi.fn(), check: vi.fn() })),
}));

describe('inputbox', () => {
    beforeEach(() => {
        document.body.innerHTML = '<div id="terminal-app"></div>';
        vi.resetModules();
        eventsOnMock.mockReset();
        showLocalMenuMock.mockReset();
        terminalInputBoxSubmitMock.mockReset();
        terminalInputBoxSubmitMock.mockResolvedValue();
        terminalTextInputMock.mockClear();
        terminalKeyPressMock.mockClear();
    });

    it('shows a history popup button and writes the selected item into the input', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        expect(eventsOnMock).toHaveBeenCalledWith('terminalInputBox', expect.any(Function));

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 7,
            title: 'Rename file',
            defaultValue: 'draft.txt',
            placeholder: 'Enter a file name',
            multiline: false,
            history: ['notes.txt', 'archive.txt'],
        });

        expect(document.querySelector('#inputbox-title h3')?.textContent).toBe('Rename file');
        expect(document.getElementById('inputbox-title').classList.contains('markdown-body')).toBe(true);
        const input = document.querySelector('.inputbox-input');
        const historyButton = document.getElementById('inputbox-history-btn');

        expect(input).not.toBeNull();
        expect(historyButton).not.toBeNull();
        expect(historyButton.style.display).toBe('inline-flex');

        historyButton.click();

        expect(showLocalMenuMock).toHaveBeenCalledWith(expect.objectContaining({
            title: 'History',
            options: ['notes.txt', 'archive.txt'],
            onSelect: expect.any(Function),
        }));

        const { onSelect } = showLocalMenuMock.mock.calls[0][0];
        onSelect(1);

        expect(input.value).toBe('archive.txt');
    });

    it('does not close when text selection drag ends on backdrop', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 9,
            title: 'Search',
            defaultValue: 'selected text',
            placeholder: 'Type...',
            multiline: false,
            history: [],
        });

        const overlay = document.getElementById('terminal-inputbox');
        const input = document.querySelector('.inputbox-input');

        expect(overlay.style.display).toBe('flex');

        // Simulate text selection that starts inside input and ends outside dialog.
        input.dispatchEvent(new MouseEvent('pointerdown', { bubbles: true }));
        overlay.dispatchEvent(new MouseEvent('pointerup', { bubbles: true }));

        expect(overlay.style.display).toBe('flex');
        expect(terminalInputBoxSubmitMock).not.toHaveBeenCalled();
    });

    it('closes when pointer down/up both occur on backdrop', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 10,
            title: 'Search',
            defaultValue: '',
            placeholder: 'Type...',
            multiline: false,
            history: [],
        });

        const overlay = document.getElementById('terminal-inputbox');

        overlay.dispatchEvent(new MouseEvent('pointerdown', { bubbles: true }));
        overlay.dispatchEvent(new MouseEvent('pointerup', { bubbles: true }));

        expect(overlay.style.display).toBe('none');
        expect(terminalInputBoxSubmitMock).toHaveBeenCalledWith(10, {
            value: '',
            variables: {},
        }, false);
    });

    it('uses themed Kill and Cancel buttons and lets Escape cancel a confirmation', async () => {
        const { initInputBox } = await import('./inputbox.js');
        initInputBox(document.createElement('canvas'));
        const openHandler = eventsOnMock.mock.calls.find(([name]) => name === 'terminalInputBox')[1];
        openHandler({
            id: 12,
            title: 'Confirm pane close',
            message: 'Are you sure you wish to kill the active pane?\n\nID: %7\nName: shell',
            confirmation: true,
            confirmLabel: 'Kill',
            cancelLabel: 'Cancel',
            confirmIsDanger: true,
        });

        expect(document.querySelector('#inputbox-title h3')?.textContent).toBe('Confirm pane close');
        expect(document.querySelector('.inputbox-confirmation-message')?.textContent).toContain('ID: %7\nName: shell');
        expect(document.querySelector('.inputbox-input')).toBeNull();
        expect(document.getElementById('inputbox-ok').textContent).toBe('Kill');
        expect(document.getElementById('inputbox-ok').classList.contains('inputbox-danger')).toBe(true);
        expect(document.getElementById('inputbox-cancel').textContent).toBe('Cancel');
        expect(document.getElementById('inputbox-escape-hint').textContent).toBe('Escape to cancel');
        const visibleButtons = Array.from(document.querySelectorAll('.inputbox-buttons button'))
            .filter((button) => button.style.display !== 'none');
        expect(visibleButtons.map((button) => button.id)).toEqual(['inputbox-ok', 'inputbox-cancel']);

        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }));
        expect(terminalInputBoxSubmitMock).toHaveBeenCalledWith(12, { value: '', variables: {} }, false);
    });

    it('uses the shared yellow Jura heading style for dialog titles', () => {
        const style = document.createElement('style');
        style.textContent = inputboxCss;
        document.head.appendChild(style);
        try {
            const title = Array.from(style.sheet.cssRules).find((rule) => rule.selectorText === '.inputbox-title').style;
            expect(title.getPropertyValue('color')).toContain('--yellow');
            expect(title.getPropertyValue('font-family')).toContain('Jura');
            expect(title.getPropertyValue('font-size')).toContain('--notes-font-size');
        } finally {
            style.remove();
        }
    });

    it('keeps the Kill background transparent at rest and red-tinted on hover', () => {
        const style = document.createElement('style');
        style.textContent = inputboxCss;
        document.head.appendChild(style);
        try {
            const rules = Array.from(style.sheet.cssRules);
            const idle = rules.find((rule) => rule.selectorText === '.inputbox-danger').style;
            const hover = rules.find((rule) => rule.selectorText === '.inputbox-danger:hover').style;
            const cancel = rules.find((rule) => rule.selectorText === '.inputbox-cancel').style;
            expect(idle.getPropertyValue('background')).toBe('transparent');
            expect(idle.getPropertyValue('border-color')).toBe(cancel.getPropertyValue('border-color'));
            expect(idle.getPropertyValue('color')).toBe(cancel.getPropertyValue('color'));
            expect(hover.getPropertyValue('border-color')).toContain('--red');
            expect(hover.getPropertyValue('background')).toContain('--terminal-red-20');
        } finally {
            style.remove();
        }
    });

    it('submits a confirmation only when Kill is activated', async () => {
        const { initInputBox } = await import('./inputbox.js');
        initInputBox(document.createElement('canvas'));
        const openHandler = eventsOnMock.mock.calls.find(([name]) => name === 'terminalInputBox')[1];
        openHandler({ id: 13, title: 'Confirm pane close', confirmation: true, confirmLabel: 'Kill' });
        document.getElementById('inputbox-ok').click();
        expect(terminalInputBoxSubmitMock).toHaveBeenCalledWith(13, { value: '', variables: {} }, true);
    });

    it('blocks terminal input during confirmation and cycles only Kill and Cancel with Tab', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const { wireKeyboardEvents } = await import('./events.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);
        wireKeyboardEvents(canvas);
        initInputBox(canvas);
        const openHandler = eventsOnMock.mock.calls.find(([name]) => name === 'terminalInputBox')[1];
        openHandler({ id: 14, title: 'Confirm pane close', confirmation: true, confirmLabel: 'Kill' });

        const kill = document.getElementById('inputbox-ok');
        const cancel = document.getElementById('inputbox-cancel');
        cancel.focus();
        window.terminalFocusedState = true;
        for (const key of ['x', 'Enter', ' ', 'ArrowUp']) {
            for (const type of ['keydown', 'keypress', 'keyup']) {
                const event = new KeyboardEvent(type, { key, bubbles: true, cancelable: true });
                cancel.dispatchEvent(event);
                expect(event.defaultPrevented).toBe(true);
            }
        }
        cancel.dispatchEvent(new KeyboardEvent('keydown', { key: 'v', ctrlKey: true, bubbles: true, cancelable: true }));
        expect(terminalTextInputMock).not.toHaveBeenCalled();
        expect(terminalKeyPressMock).not.toHaveBeenCalled();
        expect(terminalInputBoxSubmitMock).not.toHaveBeenCalled();

        cancel.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(kill);
        kill.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(cancel);
        cancel.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true, bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(kill);

        kill.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true }));
        expect(terminalInputBoxSubmitMock).toHaveBeenCalledWith(14, { value: '', variables: {} }, false);
        expect(terminalKeyPressMock).not.toHaveBeenCalled();
    });

    it('maps Home and End to current line boundaries in multiline input', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 11,
            title: 'Multiline',
            defaultValue: 'one line\ntwo line\nthree line',
            placeholder: '',
            multiline: true,
            history: [],
        });

        const input = document.querySelector('.inputbox-input');
        const secondLineStart = input.value.indexOf('two line');
        const secondLineEnd = secondLineStart + 'two line'.length;

        input.selectionStart = secondLineStart + 4;
        input.selectionEnd = secondLineStart + 4;
        input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Home', bubbles: true, cancelable: true }));
        expect(input.selectionStart).toBe(secondLineStart);
        expect(input.selectionEnd).toBe(secondLineStart);

        input.selectionStart = secondLineStart + 1;
        input.selectionEnd = secondLineStart + 1;
        input.dispatchEvent(new KeyboardEvent('keydown', { key: 'End', bubbles: true, cancelable: true }));
        expect(input.selectionStart).toBe(secondLineEnd);
        expect(input.selectionEnd).toBe(secondLineEnd);
    });

    it('submits value and typed variable map when dynamic fields are present', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 22,
            title: 'Create thing',
            defaultValue: 'base-value',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [
                {
                    name: 'textVar',
                    label: 'Text value',
                    type: 'string',
                    default: 'hello',
                    description: 'string description',
                },
                {
                    name: 'count',
                    label: 'Count',
                    type: 'int',
                    default: '7',
                },
                {
                    name: 'enabled',
                    label: 'Enabled',
                    type: 'bool',
                    default: 'false',
                },
                {
                    name: 'mode',
                    label: 'Mode',
                    type: 'list',
                    default: 'alpha, beta, gamma',
                    description: 'Pick one',
                },
            ],
        });

        const inputs = Array.from(document.querySelectorAll('.inputbox-input'));
        const mainInput = inputs[0];
        const textVarInput = inputs[1];
        const intInput = document.querySelector('input[type="number"]');
        const checkbox = document.querySelector('.inputbox-variable-checkbox');
        const listButton = document.querySelector('.inputbox-variable-list-btn');

        expect(mainInput.value).toBe('base-value');
        expect(textVarInput.value).toBe('hello');
        expect(intInput.value).toBe('7');
        expect(checkbox.checked).toBe(false);
        expect(listButton.textContent).toBe('alpha');

        listButton.click();
        expect(showLocalMenuMock).toHaveBeenLastCalledWith(expect.objectContaining({
            title: 'Mode',
            options: ['alpha', 'beta', 'gamma'],
            onSelect: expect.any(Function),
        }));
        const { onSelect } = showLocalMenuMock.mock.calls[0][0];
        onSelect(2);
        expect(listButton.textContent).toBe('gamma');

        mainInput.value = 'updated-base';
        textVarInput.value = 'updated-string';
        intInput.value = '42';
        checkbox.checked = true;

        document.getElementById('inputbox-ok').click();

        expect(terminalInputBoxSubmitMock).toHaveBeenCalledWith(22, {
            value: 'updated-base',
            variables: {
                textVar: 'updated-string',
                count: 42,
                enabled: true,
                mode: 'gamma',
            },
        }, true);
    });

    it('shows list default from options array when provided', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 23,
            title: 'Select model',
            defaultValue: 'prompt',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [
                {
                    name: 'model',
                    label: 'Model',
                    type: 'list',
                    default: 'Claude: claude-3-5-sonnet',
                    options: [
                        'ChatGPT: gpt-4o',
                        'Claude: claude-3-5-sonnet',
                        'Claude: claude-3-haiku',
                    ],
                },
            ],
        });

        const listButton = document.querySelector('.inputbox-variable-list-btn');
        expect(listButton).toBeTruthy();
        expect(listButton.textContent).toBe('Claude: claude-3-5-sonnet');

        listButton.click();
        expect(showLocalMenuMock).toHaveBeenLastCalledWith(expect.objectContaining({
            title: 'Model',
            options: [
                'ChatGPT: gpt-4o',
                'Claude: claude-3-5-sonnet',
                'Claude: claude-3-haiku',
            ],
            onSelect: expect.any(Function),
        }));
    });

    it('keeps Tab navigation inside the inputbox modal while open', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        const outsideButton = document.createElement('button');
        outsideButton.type = 'button';
        outsideButton.id = 'outside-focus-target';
        document.body.appendChild(outsideButton);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 31,
            title: 'Tab trap',
            defaultValue: 'first',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [
                {
                    name: 'second',
                    label: 'Second',
                    type: 'string',
                    default: 'value',
                },
            ],
        });

        const inputs = Array.from(document.querySelectorAll('.inputbox-input'));
        const firstField = inputs[0];
        const secondField = inputs[1];

        await new Promise((resolve) => setTimeout(resolve, 0));

        firstField.focus();
        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(secondField);

        // Tab from second field wraps back to first (buttons excluded from tab cycle)
        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(firstField);

        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true, bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(secondField);

        outsideButton.focus();
        expect(document.activeElement).toBe(firstField);
        document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true, cancelable: true }));
        expect(document.activeElement).toBe(secondField);
    });

    it('prevents focus from escaping to outside elements while inputbox is open', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        const outsideInput = document.createElement('input');
        outsideInput.type = 'text';
        outsideInput.id = 'outside-input-target';
        document.body.appendChild(outsideInput);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 32,
            title: 'Focus lock',
            defaultValue: 'inside',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [],
        });

        const insideInput = document.querySelector('.inputbox-input');
        await new Promise((resolve) => setTimeout(resolve, 0));
        expect(document.activeElement).toBe(insideInput);

        outsideInput.focus();
        expect(document.activeElement).toBe(insideInput);
    });

    it('uses shared terminal focus state function while inputbox is open', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        const setFocusMock = vi.fn();
        window.ttyphoonSetTerminalFocusState = setFocusMock;
        window.terminalFocusedState = true;

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 33,
            title: 'Shared focus',
            defaultValue: 'abc',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [],
        });

        expect(setFocusMock).toHaveBeenCalledWith(false, { focusVisible: false, force: true });

        document.getElementById('inputbox-cancel').click();

        expect(setFocusMock).toHaveBeenCalledWith(true, { focusVisible: false, force: true });
        delete window.ttyphoonSetTerminalFocusState;
    });

    it('consumes Ctrl+Tab while inputbox is open', async () => {
        const { initInputBox } = await import('./inputbox.js');
        const canvas = document.createElement('canvas');
        document.body.appendChild(canvas);

        initInputBox(canvas);

        const openHandler = eventsOnMock.mock.calls[0][1];
        openHandler({
            id: 34,
            title: 'Ctrl+Tab guard',
            defaultValue: 'inside',
            placeholder: '',
            multiline: false,
            history: [],
            variables: [],
        });

        await new Promise((resolve) => setTimeout(resolve, 0));
        const insideInput = document.querySelector('.inputbox-input');
        expect(document.activeElement).toBe(insideInput);

        const evt = new KeyboardEvent('keydown', {
            key: 'Tab',
            ctrlKey: true,
            bubbles: true,
            cancelable: true,
        });

        const dispatchResult = document.dispatchEvent(evt);
        expect(dispatchResult).toBe(false);
        expect(document.activeElement).toBe(insideInput);
    });
});