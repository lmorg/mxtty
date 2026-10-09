import { describe, expect, it, vi } from 'vitest';

import { initNotesGitPanel } from './notes-git-panel.js';

function createElements() {
    const gitTab = document.createElement('button');
    gitTab.style.display = 'none';
    return {
        gitTab,
        gitPane: document.createElement('div'),
        gitWorkspace: document.createElement('button'),
        gitBranch: document.createElement('span'),
        gitAheadBehind: document.createElement('span'),
        gitRefresh: document.createElement('button'),
        gitFileList: document.createElement('div'),
    };
}

async function flushPromises() {
    await new Promise((resolve) => setTimeout(resolve, 0));
    await new Promise((resolve) => setTimeout(resolve, 0));
}

function createApi(overrides = {}) {
    return {
        getCurrentProject: vi.fn(() => Promise.resolve('/work/repo')),
        getStatus: vi.fn(() => Promise.resolve({
            root: '/work/repo',
            branch: 'main',
            ahead: 2,
            behind: 1,
            clean: false,
            files: [
                { path: 'src/main.go', code: ' M', staged: false, unstaged: true },
                { path: 'README.md', code: 'A ', staged: true, unstaged: false },
            ],
        })),
        getDiff: vi.fn(() => Promise.resolve({ path: 'src/main.go', original: 'old', working: 'new' })),
        stageFile: vi.fn(() => Promise.resolve()),
        unstageFile: vi.fn(() => Promise.resolve()),
        eventsOn: vi.fn(),
        ...overrides,
    };
}

describe('Notes Git panel', () => {
    it('shows repository status and exposes workspace, stage, and diff actions', async () => {
        const elements = createElements();
        const api = createApi();
        const showWorkspaceMenu = vi.fn();
        const onOpenDiff = vi.fn();
        const beforeFileAction = vi.fn(() => Promise.resolve(true));
        const panel = initNotesGitPanel(elements, api, { showWorkspaceMenu, onOpenDiff, beforeFileAction });
        await flushPromises();

        expect(elements.gitTab.style.display).toBe('');
        expect(elements.gitBranch.textContent).toBe('main');
        expect(elements.gitAheadBehind.textContent).toBe('↑2 ↓1');
        expect(elements.gitFileList.textContent).toContain('src/main.go');
        expect(elements.gitFileList.textContent).toContain('README.md');

        elements.gitWorkspace.click();
        expect(showWorkspaceMenu).toHaveBeenCalledWith('/work/repo', elements.gitWorkspace);

        const stageButton = elements.gitFileList.querySelector('[data-git-action="stage"]');
        stageButton.click();
        await flushPromises();
        expect(beforeFileAction).toHaveBeenCalledWith({
            root: '/work/repo',
            path: 'src/main.go',
            action: 'stage',
        });
        expect(api.stageFile).toHaveBeenCalledWith('/work/repo', 'src/main.go');

        const fileButton = elements.gitFileList.querySelector('.notes-git-file');
        fileButton.click();
        await flushPromises();
        expect(api.getDiff).toHaveBeenCalledWith('/work/repo', 'src/main.go', '');
        expect(onOpenDiff).toHaveBeenCalledWith(expect.objectContaining({
            root: '/work/repo',
            file: expect.objectContaining({ path: 'src/main.go' }),
            diff: expect.objectContaining({ original: 'old', working: 'new' }),
        }));
        expect(panel.getCurrentRoot()).toBe('/work/repo');
    });

    it('hides the Git tab when the project root is not a usable Git root', async () => {
        const elements = createElements();
        const api = createApi({
            getStatus: vi.fn(() => Promise.reject(new Error('project root is not a Git working tree'))),
        });

        initNotesGitPanel(elements, api);
        await flushPromises();

        expect(elements.gitTab.style.display).toBe('none');
        expect(elements.gitFileList.textContent).toContain('project root is not a Git working tree');
    });
});
