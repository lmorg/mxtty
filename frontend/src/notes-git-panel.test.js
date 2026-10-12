import { describe, expect, it, vi } from 'vitest';

import { initNotesGitPanel } from './notes-git-panel.js';

function createElements() {
    const gitTab = document.createElement('button');
    gitTab.style.display = 'none';
    const gitPane = document.createElement('div');
    const gitWorkspace = document.createElement('button');
    const gitCompare = document.createElement('button');
    const gitCommits = document.createElement('button');
    const gitRemote = document.createElement('button');
    const gitBranch = document.createElement('span');
    const gitAheadBehind = document.createElement('span');
    const gitRefresh = document.createElement('button');
    const gitFileList = document.createElement('div');
    const gitCommitList = document.createElement('div');
    gitPane.append(gitWorkspace, gitCompare, gitRemote, gitBranch, gitAheadBehind, gitRefresh, gitFileList);
    return {
        gitTab,
        gitPane,
        gitWorkspace,
        gitCompare,
        gitCommits,
        gitRemote,
        gitBranch,
        gitAheadBehind,
        gitRefresh,
        gitFileList,
        gitCommitList,
    };
}

async function flushPromises() {
    await new Promise((resolve) => setTimeout(resolve, 0));
    await new Promise((resolve) => setTimeout(resolve, 0));
}

function createApi(overrides = {}) {
    return {
        getCurrentProject: vi.fn(() => Promise.resolve('/work/repo')),
        filterStrings: vi.fn((query, items) => Promise.resolve({
            List: items.filter((item) => item.toLowerCase().includes(query.toLowerCase())),
        })),
        listCommits: vi.fn(() => Promise.resolve([
            { hash: 'abcdef1', subject: 'Merge feature', date: '2026-10-12 12:30:00', author: 'Test User', graph: ['*   ', '|\\  '] },
            { hash: 'abcdef2', subject: 'Feature commit', date: '2026-10-12 12:00:00', author: 'Test User', graph: ['| * ', '|/  '] },
        ])),
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
    it('opens a commit context menu without changing comparison until requested', async () => {
        const elements = createElements();
        const api = createApi();
        const showCommitMenu = vi.fn();
        const panel = initNotesGitPanel(elements, api, { showCommitMenu });
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        const event = new MouseEvent('contextmenu', { bubbles: true, cancelable: true, clientX: 25, clientY: 40 });
        elements.gitCommitList.querySelector('.notes-git-commit').dispatchEvent(event);
        expect(event.defaultPrevented).toBe(true);
        expect(panel.getCompareBranch()).toBe('');
        expect(showCommitMenu).toHaveBeenCalledWith(expect.objectContaining({ root: '/work/repo', x: 25, y: 40, commit: expect.objectContaining({ hash: 'abcdef1' }) }));
        showCommitMenu.mock.calls[0][0].compare();
        await flushPromises();
        expect(panel.getCompareBranch()).toBe('abcdef1');
    });

    it('selects a commit through the same comparison refresh and diff path as the branch menu', async () => {
        const elements = createElements();
        const fullHash = 'abcdef1' + '0'.repeat(33);
        const api = createApi({
            refreshStatus: vi.fn(() => Promise.resolve({ root: '/work/repo', branch: 'main', files: [], clean: true })),
            listCommits: vi.fn(() => Promise.resolve([
                { hash: 'abcdef1', fullHash, subject: 'Base commit', date: '', author: '', graph: ['*'] },
            ])),
        });
        const onStatusChange = vi.fn();
        const panel = initNotesGitPanel(elements, api, { onStatusChange });
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        elements.gitCommitList.querySelector('.notes-git-commit').click();
        await flushPromises();
        expect(elements.gitCompare.textContent).toBe('abcdef1');
        expect(panel.getCompareBranch()).toBe(fullHash);
        expect(api.refreshStatus).toHaveBeenCalledWith('/work/repo');
        expect(api.getBranchDiffFiles).toHaveBeenCalledWith('/work/repo', fullHash);
        expect(onStatusChange).toHaveBeenLastCalledWith(expect.objectContaining({ files: [expect.objectContaining({ path: 'src/main.go' })] }), '/work/repo');
        elements.gitFileList.querySelector('.notes-git-file').click();
        await flushPromises();
        expect(api.getDiffAgainstBranch).toHaveBeenCalledWith('/work/repo', 'src/main.go', '', fullHash);
        await panel.setCompareBranch('');
        expect(elements.gitCompare.textContent).toBe('HEAD');
        for (const key of ['Enter', ' ']) {
            api.refreshStatus.mockClear();
            elements.gitCommitList.querySelector('.notes-git-commit').dispatchEvent(new KeyboardEvent('keydown', { key, bubbles: true }));
            await flushPromises();
            expect(elements.gitCompare.textContent).toBe('abcdef1');
            expect(api.refreshStatus).toHaveBeenCalledOnce();
            await panel.setCompareBranch('');
            await flushPromises();
        }
    });

    it('filters full hashes and full messages with the shared filter without reloading Git per query', async () => {
        const elements = createElements();
        const commits = [
            { hash: 'aaaaaaa', fullHash: 'a'.repeat(34) + 'ef0123', subject: 'First merge', message: 'First merge\n\nRollout ticket', date: '', author: 'Test User', graph: ['*'] },
            { hash: 'bbbbbbb', fullHash: 'b'.repeat(40), subject: 'Second commit', message: 'Second commit', date: '', author: 'Test User', graph: ['*'] },
        ];
        const api = createApi({ listCommits: vi.fn(() => Promise.resolve(commits)) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        await panel.setCommitFilter('ef0123');
        expect(api.listCommits).toHaveBeenLastCalledWith('/work/repo', 10000);
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(1);
        expect(elements.gitCommitList.querySelector('.notes-git-commit-hash').textContent).toBe('aaaaaaa');
        expect(elements.gitCommitList.querySelector('.notes-git-commit').classList.contains('notes-git-commit-filtered')).toBe(true);
        expect(api.filterStrings).toHaveBeenCalledWith('ef0123', commits.map((commit) => `${commit.fullHash}\n${commit.message}`));
        await panel.setCommitFilter('ROLLOUT');
        expect(elements.gitCommitList.querySelector('.notes-git-commit-subject').textContent).toBe('First merge');
        expect(api.listCommits).toHaveBeenCalledTimes(2);
        await panel.setCommitFilter('Test User');
        expect(elements.gitCommitList.textContent).toBe('No matching commits');
        await panel.setCommitFilter('');
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(2);
        expect(elements.gitCommitList.querySelector('.notes-git-commit-filtered')).toBeNull();
        expect(elements.gitCommitList.querySelector('.notes-git-commit-graph')).not.toBeNull();
        expect(api.listCommits).toHaveBeenCalledTimes(2);
    });

    it('uses the current shared query when re-enabling Commits after the filter was cleared', async () => {
        const elements = createElements();
        let query = 'Merge';
        const panel = initNotesGitPanel(elements, createApi(), { getCommitFilter: () => query });
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(1);
        panel.setCommitsActive(false);
        query = '';
        panel.setCommitsActive(true);
        await flushPromises();
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(2);
    });

    it('ignores a delayed filter response after a newer commit query', async () => {
        const elements = createElements();
        const api = createApi();
        let resolveFilter;
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        api.filterStrings.mockImplementationOnce(() => new Promise((resolve) => { resolveFilter = resolve; }));
        const oldQuery = panel.setCommitFilter('Merge');
        await flushPromises();
        await panel.setCommitFilter('Feature commit');
        resolveFilter({ List: ['abcdef1\nMerge feature'] });
        await oldQuery;
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(1);
        expect(elements.gitCommitList.querySelector('.notes-git-commit-subject').textContent).toBe('Feature commit');
    });

    it('continues a56202f0 collapsing slashes past adjacent vertical lanes', async () => {
        const elements = createElements();
        const api = createApi({ listCommits: vi.fn(() => Promise.resolve([
            {
                hash: 'a56202f0', subject: 'Merge', date: '', author: '',
                graph: String.raw`| | * | | | | | | | |
| | |\ \ \ \ \ \ \ \ \
| | |/ / / / / / / / /
| |/| | | | | | | | |
`.trimEnd().split('\n'),
            },
            { hash: 'f22599f3', subject: 'Merge', date: '', author: '', graph: ['| * | | | | | | | | |'] },
        ])) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        const path = elements.gitCommitList.querySelector('svg path').getAttribute('d');
        expect(path).toContain('M 44 45 L 28 63');
        expect(path).toContain('M 36 45 L 36 63');
        expect(path).not.toContain('M 44 45 L 36 63');
        expect(path).not.toContain('M 36 45 L 28 63');
    });

    it('updates overlap masking on resize and releases its observer on toggle-off', async () => {
        let resized;
        const disconnect = vi.fn();
        vi.stubGlobal('ResizeObserver', class {
            constructor(callback) { resized = callback; }
            observe() {}
            disconnect() { disconnect(); }
        });
        try {
            const elements = createElements();
            const panel = initNotesGitPanel(elements, createApi());
            await flushPromises();
            panel.setCommitsActive(true);
            await flushPromises();
            const row = elements.gitCommitList.querySelector('.notes-git-commit');
            const graph = row.querySelector('.notes-git-commit-graph');
            const content = row.querySelector('.notes-git-commit-content');
            vi.spyOn(graph, 'getBoundingClientRect').mockReturnValue({ width: 600, left: 10 });
            const bounds = vi.spyOn(content, 'getBoundingClientRect').mockReturnValue({ left: 130 });
            resized([{ target: row }]);
            expect(row.style.getPropertyValue('--git-overlap-start')).toBe('120px');
            bounds.mockReturnValue({ left: 82 });
            resized([{ target: row }]);
            expect(row.style.getPropertyValue('--git-overlap-start')).toBe('72px');
            panel.setCommitsActive(false);
            expect(disconnect).toHaveBeenCalledOnce();
        } finally {
            vi.unstubAllGlobals();
        }
    });

    it.each([
        ['a56202f0 direction reversal', String.raw`| | * | | | | | | | |
| | |\ \ \ \ \ \ \ \ \
| | |/ / / / / / / / /
| |/| | | | | | | | |
| * | | | | | | | | |
`],
        ['eebc4be8 underscore bridge', String.raw`* | | | | | | | | | | | | | |
| |_|_|_|_|_|_|_|_|_|_|_|_|/
|/| | | | | | | | | | | | |
* | | | | | | | | | | | | |
`],
        ['9a1dddf5 nested collapse', String.raw`| | | | * | | | | | | | | | | |
| | | |/ / / / / / / / / / / /
| | |/| | | | | | | | | | | |
| | * | | | | | | | | | | | |
`],
        ['f601effe overlapping collapse', String.raw`| | | | | | | | | | | | * | | | | | | | | | | | | | | | | | |
| | | | | | | | | | | |/| | | | | | | | | | | | | | | | | | |
| | | | | | | | | | |/|/ / / / / / / / / / / / / / / / / / /
| | | | | | | | | | | * | | | | | | | | | / / / / / / / / /
`],
    ])('preserves all incoming diagonal lanes in %s', async (_name, raw) => {
        const graph = raw.trimEnd().split('\n');
        const elements = createElements();
        const api = createApi({ listCommits: vi.fn(() => Promise.resolve([
            { hash: 'fixture', subject: 'Merge', date: '', author: '', graph },
        ])) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        const path = elements.gitCommitList.querySelector('svg path').getAttribute('d');
        const segments = [...path.matchAll(/M ([\d.-]+) ([\d.-]+) L ([\d.-]+) ([\d.-]+)/g)]
            .map((match) => match.slice(1).map(Number));
        const step = 72 / graph.length;
        for (let row = 1; row < graph.length; row++) {
            for (const [column, character] of [...graph[row]].entries()) {
                if (!'/\\'.includes(character) || '_-'.includes(graph[row][column - 1] || ' ')) continue;
                const matches = segments.filter((segment) => Math.abs(segment[1] - (row - 0.5) * step) < 0.001
                    && Math.abs(segment[3] - (row + 0.5) * step) < 0.001 && segment[2] === column * 8 + 4);
                expect(matches, `missing input for ${character} at row ${row}, column ${column}`).not.toHaveLength(0);
                if (graph[row - 1][column] === '\\' && character === '/') {
                    expect(matches.some((segment) => segment[0] === column * 8 + 4), `shifted reversal at row ${row}, column ${column}`).toBe(true);
                }
            }
        }
    });

    it('connects octopus merge elbows and collapsing lanes to their actual neighbours', async () => {
        const elements = createElements();
        const graphs = [
            ['*-. ', '|\\ \\'],
            ['| | * ', '| |/ '],
            ['| * ', '|/ '],
            ['* '],
        ];
        const api = createApi({ listCommits: vi.fn(() => Promise.resolve(graphs.map((graph, index) => ({
            hash: String(index), subject: 'Merge', date: '', author: '', graph,
        })))) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        const paths = [...elements.gitCommitList.querySelectorAll('svg path')].map((path) => path.getAttribute('d'));
        expect(paths[0]).toContain('M 20 18 L 28 54');
        expect(paths[0]).toContain('M 4 18 H 20');
        expect(paths[1]).toContain('M 36 18 L 28 54');
        expect(paths[1]).toContain('M 28 54 L 24 72');
        expect(paths[2]).toContain('M 24 0 L 20 18');
        expect(paths[2]).toContain('M 20 18 L 12 54');
        expect(paths[3]).toContain('M 8 0 L 4 36');
        expect(paths[3]).not.toContain('72');
    });

    it('joins expanding diagonal lanes across consecutive commit rows', async () => {
        const elements = createElements();
        const graphs = [
            ['| | * |   ', '| | |\\ \\'],
            ['| | * \\ \\   ', '| | |\\ \\ \\'],
            ['| | * \\ \\ \\   ', '| | |\\ \\ \\ \\'],
            ['| * | \\ \\ \\ \\   ', '| |\\ \\ \\ \\ \\ \\'],
        ];
        const api = createApi({ listCommits: vi.fn(() => Promise.resolve(graphs.map((graph, index) => ({
            hash: String(index), subject: 'Merge', date: '', author: '', graph,
        })))) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        await flushPromises();
        const paths = [...elements.gitCommitList.querySelectorAll('svg path')].map((path) => path.getAttribute('d'));
        const endpoints = (path, top) => {
            const positions = [];
            for (const match of path.matchAll(/M ([\d.-]+) ([\d.-]+) (?:L ([\d.-]+) ([\d.-]+)|V ([\d.-]+)|H ([\d.-]+))/g)) {
                const startX = Number(match[1]);
                const startY = Number(match[2]);
                const endX = match[3] !== undefined ? Number(match[3]) : match[6] !== undefined ? Number(match[6]) : startX;
                const endY = match[4] !== undefined ? Number(match[4]) : match[5] !== undefined ? Number(match[5]) : startY;
                if (top && startY === 0) positions.push(startX);
                if (!top && endY === 72) positions.push(endX);
            }
            return [...new Set(positions)].sort((left, right) => left - right);
        };
        for (let index = 0; index < paths.length - 1; index++) {
            expect(endpoints(paths[index], false)).toEqual(endpoints(paths[index + 1], true));
        }
        expect(endpoints(paths[0], false)).toEqual([4, 20, 36, 48, 64]);
    });

    it('discards delayed commit responses after switching back to files', async () => {
        const elements = createElements();
        let resolveHistory;
        const api = createApi({ listCommits: vi.fn(() => new Promise((resolve) => { resolveHistory = resolve; })) });
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        panel.setCommitsActive(true);
        panel.setCommitsActive(false);
        resolveHistory([{ hash: 'late', subject: 'Old history', date: '', author: '', graph: ['*'] }]);
        await flushPromises();
        expect(elements.gitCommitList.querySelector('.notes-git-commit')).toBeNull();
        expect(elements.gitFileList.querySelector('.notes-git-file')).not.toBeNull();
    });

    it('toggles a graph-backed commit list and restores files on workspace reset', async () => {
        const elements = createElements();
        const api = createApi();
        const panel = initNotesGitPanel(elements, api);
        await flushPromises();
        elements.gitCommits.click();
        await flushPromises();
        expect(elements.gitCommits.getAttribute('aria-pressed')).toBe('true');
        expect(api.listCommits).toHaveBeenCalledWith('/work/repo', 100);
        expect(elements.gitCommitList.querySelectorAll('.notes-git-commit')).toHaveLength(2);
        expect(elements.gitCommitList.querySelector('.notes-git-commit-hash').textContent).toBe('abcdef1');
        expect(elements.gitCommitList.querySelector('.notes-git-commit-subject').textContent).toBe('Merge feature');
        expect(Array.from(elements.gitCommitList.querySelector('.notes-git-commit').querySelectorAll('.notes-git-commit-meta'), (line) => line.textContent)).toEqual([
            '2026-10-12 12:30:00', '(Test User)',
        ]);
        expect(elements.gitCommitList.querySelector('svg path').getAttribute('d')).toContain('L');
        expect(elements.gitCommitList.querySelector('.notes-git-commit').style.getPropertyValue('--git-graph-width')).toBe('32px');
        expect(elements.gitFileList.querySelector('.notes-git-commit')).toBeNull();
        expect(elements.gitFileList.querySelector('.notes-git-file')).not.toBeNull();
        panel.setCommitsActive(false);
        expect(elements.gitCommits.getAttribute('aria-pressed')).toBe('false');
        expect(elements.gitFileList.querySelector('.notes-git-file')).not.toBeNull();
    });

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
