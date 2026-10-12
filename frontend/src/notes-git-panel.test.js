import { describe, expect, it, vi } from 'vitest';

import { initNotesGitPanel } from './notes-git-panel.js';

function createElements() {
    const gitTab = document.createElement('button');
    gitTab.style.display = 'none';
    const gitPane = document.createElement('div');
    const gitWorkspace = document.createElement('button');
    const gitCompare = document.createElement('button');
    const gitRemote = document.createElement('button');
    const gitBranch = document.createElement('span');
    const gitAheadBehind = document.createElement('span');
    const gitRefresh = document.createElement('button');
    const gitFileList = document.createElement('div');
    gitPane.append(gitWorkspace, gitCompare, gitRemote, gitBranch, gitAheadBehind, gitRefresh, gitFileList);
    return {
        gitTab,
        gitPane,
        gitWorkspace,
        gitCompare,
        gitRemote,
        gitBranch,
        gitAheadBehind,
        gitRefresh,
        gitFileList,
    };
}

async function flushPromises() {
    await new Promise((resolve) => setTimeout(resolve, 0));
    await new Promise((resolve) => setTimeout(resolve, 0));
}

function createApi(overrides = {}) {
    return {
        getCurrentProject: vi.fn(() => Promise.resolve('/work/repo')),
        getBrowserLink: vi.fn(() => Promise.resolve({
            provider: 'github',
            url: 'https://github.com/org/repo/tree/main',
        })),
        getBranchDiffFiles: vi.fn(() => Promise.resolve([
            { path: 'src/main.go', code: ' M', staged: false, unstaged: true, compareOnly: true },
        ])),
        openRemote: vi.fn(() => Promise.resolve()),
        getStatus: vi.fn(() => Promise.resolve({
            root: '/work/repo',
            branch: 'main',
            ahead: 2,
            behind: 1,
            clean: false,
            files: [
                { path: 'src/main.go', code: ' M', staged: false, unstaged: true },
                { path: 'README.md', code: 'A ', staged: true, unstaged: false },
                { path: 'new-file.txt', code: '??', staged: false, unstaged: false, untracked: true },
            ],
        })),
        getDiff: vi.fn(() => Promise.resolve({ path: 'src/main.go', original: 'old', working: 'new' })),
        getDiffAgainstBranch: vi.fn(() => Promise.resolve({ path: 'src/main.go', original: 'branch', working: 'new' })),
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
        const showTooltip = vi.fn();
        const hideTooltip = vi.fn();
        const onOpenDiff = vi.fn();
        const beforeFileAction = vi.fn(() => Promise.resolve(true));
        const panel = initNotesGitPanel(elements, api, {
            showWorkspaceMenu, showTooltip, hideTooltip, onOpenDiff, beforeFileAction,
        });
        await flushPromises();

        expect(elements.gitTab.style.display).toBe('');
        expect(elements.gitWorkspace.textContent).toBe('Worktrees');
        expect(elements.gitBranch.textContent).toBe('main');
        expect(elements.gitCompare.textContent).toBe('HEAD');
        expect(elements.gitAheadBehind.textContent).toBe('↑2 ↓1');
        expect(elements.gitRemote.hidden).toBe(false);
        expect(elements.gitRemote.textContent).toBe(String.fromCodePoint(0xf09b));
        expect(elements.gitRemote.getAttribute('aria-label')).toBe('Open GitHub repository');
        elements.gitRemote.click();
        expect(api.openRemote).toHaveBeenCalledWith('https://github.com/org/repo/tree/main');
        expect(elements.gitFileList.textContent).toContain('src/main.go');
        expect(elements.gitFileList.textContent).toContain('README.md');
        const statusCode = elements.gitFileList.querySelector('.notes-git-status-code');
        statusCode.dispatchEvent(new MouseEvent('mouseover', { bubbles: true }));
        expect(showTooltip).toHaveBeenCalledWith(expect.stringContaining('working tree (M) modified'));
        statusCode.dispatchEvent(new MouseEvent('mouseout', { bubbles: true, relatedTarget: document.body }));
        expect(hideTooltip).toHaveBeenCalled();

        elements.gitWorkspace.click();
        expect(showWorkspaceMenu).toHaveBeenCalledWith('/work/repo', elements.gitWorkspace);

        const stageButton = elements.gitFileList.querySelector('[data-git-action="stage"]');
        const unstageButton = elements.gitFileList.querySelector('[data-git-action="unstage"]');
        expect(stageButton.classList.contains('notes-git-action-stage')).toBe(true);
        expect(unstageButton.classList.contains('notes-git-action-unstage')).toBe(true);
        expect(stageButton.textContent).toBe('Unstaged');
        expect(stageButton.dataset.tooltip).toBe('Click to stage');
        expect(stageButton.hasAttribute('title')).toBe(false);
        expect(unstageButton.textContent).toBe('Staged');
        expect(unstageButton.dataset.tooltip).toBe('Click to unstage');
        expect(unstageButton.hasAttribute('title')).toBe(false);
        const untrackedButton = elements.gitFileList.querySelector('[data-git-action="stage"][data-path="new-file.txt"]');
        expect(untrackedButton.textContent).toBe('Untracked');
        expect(untrackedButton.classList.contains('notes-git-action-stage')).toBe(true);
        expect(untrackedButton.classList.contains('notes-git-action-untracked')).toBe(true);
        expect(untrackedButton.dataset.tooltip).toBe('Click to stage');
        expect(untrackedButton.getAttribute('aria-label')).toBe('Click to stage');
        stageButton.dispatchEvent(new MouseEvent('mouseover', { bubbles: true }));
        expect(showTooltip).toHaveBeenCalledWith('Click to stage');
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

    it('switches comparison branches without changing the status branch and opens branch-based diffs', async () => {
        const elements = createElements();
        const api = createApi();
        const showCompareMenu = vi.fn();
        const onOpenDiff = vi.fn();
        const panel = initNotesGitPanel(elements, api, { showCompareMenu, onOpenDiff });
        await flushPromises();

        elements.gitCompare.click();
        expect(showCompareMenu).toHaveBeenCalledWith(
            '/work/repo', 'main', '', elements.gitCompare, expect.any(Function),
        );
        await panel.setCompareBranch('release/2');
        await flushPromises();

        expect(elements.gitCompare.textContent).toBe('release/2');
        expect(elements.gitBranch.textContent).toBe('main');
        expect(api.getBranchDiffFiles).toHaveBeenCalledWith('/work/repo', 'release/2');
        expect(elements.gitFileList.querySelector('.notes-git-status-code').textContent).toBe('≠');
        expect(panel.getCachedStatus().files).toHaveLength(3);

        elements.gitFileList.querySelector('.notes-git-file').click();
        await flushPromises();
        expect(api.getDiffAgainstBranch).toHaveBeenCalledWith('/work/repo', 'src/main.go', '', 'release/2');
        expect(api.getDiff).not.toHaveBeenCalled();
        expect(onOpenDiff).toHaveBeenCalledWith(expect.objectContaining({
            diff: expect.objectContaining({ original: 'branch' }),
        }));
    });

    it('keeps the Git panel available when the selected comparison branch disappears', async () => {
        const elements = createElements();
        const api = createApi({
            getBranchDiffFiles: vi.fn(() => Promise.reject(new Error('comparison branch does not exist'))),
        });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();

        await panel.setCompareBranch('deleted-branch');

        expect(elements.gitTab.style.display).toBe('');
        expect(elements.gitFileList.textContent).toContain('Could not compare with deleted-branch');
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

    it('shows provider-specific brand glyphs and hides unsupported remotes', async () => {
        const elements = createElements();
        const api = createApi();
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();

        for (const [provider, glyph] of [
            ['azuredevops', 0xf3ca],
            ['gitlab', 0xf296],
            ['bitbucket', 0xf171],
        ]) {
            api.getBrowserLink.mockResolvedValueOnce({ provider, url: `https://${provider}.example/repo` });
            await panel.refresh(false);
            await flushPromises();
            expect(elements.gitRemote.hidden).toBe(false);
            expect(elements.gitRemote.textContent).toBe(String.fromCodePoint(glyph));
        }

        api.getBrowserLink.mockResolvedValueOnce({ provider: '', url: '' });
        await panel.refresh(false);
        await flushPromises();
        expect(elements.gitRemote.hidden).toBe(true);
    });

    it('rehydrates its root cache from watcher-published Git status', async () => {
        const elements = createElements();
        const api = createApi();
        const eventHandlers = new Map();
        api.eventsOn.mockImplementation((name, handler) => eventHandlers.set(name, handler));
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();

        eventHandlers.get('notesGitStatusChanged')({
            root: '/work/repo',
            status: {
                root: '/work/repo',
                branch: 'updated',
                ahead: 0,
                behind: 0,
                clean: false,
                files: [{ path: 'other.txt', code: '??', untracked: true }],
            },
        });

        expect(elements.gitBranch.textContent).toBe('updated');
        expect(panel.getCachedStatus('/work/repo').files[0].path).toBe('other.txt');
    });

    it('starts the metadata watch only while the Git tab is active', async () => {
        const elements = createElements();
        const api = createApi({ setStatusWatch: vi.fn(() => Promise.resolve()) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();

        panel.setWatchActive(true);
        expect(api.setStatusWatch).toHaveBeenLastCalledWith('/work/repo', true);

        panel.setWatchActive(false);
        expect(api.setStatusWatch).toHaveBeenLastCalledWith('/work/repo', false);
    });
});
