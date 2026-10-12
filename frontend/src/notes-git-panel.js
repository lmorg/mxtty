function statusCode(file) {
    return String(file?.code || '??');
}

function statusTooltip(code) {
    if (code === '??') {
        return '??: untracked file; Git has not added this path to the index';
    }
    const meaning = {
        ' ': 'no change',
        A: 'added',
        M: 'modified',
        D: 'deleted',
        R: 'renamed',
        C: 'copied',
        U: 'unmerged conflict',
        T: 'file type changed',
        '!': 'ignored',
        '?': 'untracked',
    };
    const [index, worktree] = String(code || '  ').padEnd(2, ' ');
    return `${code}: index (${index}) ${meaning[index] || 'changed'}; working tree (${worktree}) ${meaning[worktree] || 'changed'}`;
}

function rootKey(root) {
    return String(root || '').replaceAll('\\', '/').replace(/\/+$/, '');
}

const gitRemoteIcons = {
    github: [0xf09b, 'GitHub'],
    azuredevops: [0xf3ca, 'Azure DevOps'],
    gitlab: [0xf296, 'GitLab'],
    bitbucket: [0xf171, 'Bitbucket'],
};

export function initNotesGitPanel(elements, api, options = {}) {
    if (!elements?.gitTab || !elements?.gitPane || !elements?.gitFileList || !elements?.gitRemote || !elements?.gitCompare) {
        return null;
    }

    let projectRoot = '';
    let requestGeneration = 0;
    let currentStatus = null;
    let compareBranch = '';
    let watchActive = false;
    let remoteRequest = 0;
    const statusCache = new Map();

    const syncStatusWatch = () => {
        if (typeof api.setStatusWatch !== 'function') {
            return;
        }
        Promise.resolve(api.setStatusWatch(projectRoot, watchActive && Boolean(currentStatus))).catch(() => {});
    };

    const setAvailable = (available) => {
        elements.gitTab.style.display = available ? '' : 'none';
        if (!available) {
            elements.gitPane.dataset.active = 'false';
            elements.gitTab.setAttribute('aria-selected', 'false');
        }
        options.onAvailabilityChange?.(available);
    };

    const showMessage = (message, kind = 'info') => {
        elements.gitFileList.replaceChildren();
        const item = document.createElement('div');
        item.className = `notes-git-message notes-git-message-${kind}`;
        item.textContent = message;
        elements.gitFileList.appendChild(item);
    };

    const hideRemoteLink = () => {
        elements.gitRemote.hidden = true;
        delete elements.gitRemote.dataset.url;
    };

    const updateRemoteLink = async (status, root) => {
        const request = ++remoteRequest;
        if (!status || typeof api.getBrowserLink !== 'function' || typeof api.openRemote !== 'function') {
            hideRemoteLink();
            return;
        }
        try {
            const link = await api.getBrowserLink(root, status.branch || '');
            if (request !== remoteRequest || root !== projectRoot) {
                return;
            }
            const remote = gitRemoteIcons[link?.provider];
            if (!remote || !link?.url) {
                hideRemoteLink();
                return;
            }
            const [glyph, name] = remote;
            elements.gitRemote.textContent = String.fromCodePoint(glyph);
            elements.gitRemote.dataset.url = link.url;
            elements.gitRemote.dataset.tooltip = `Open ${name} repository`;
            elements.gitRemote.setAttribute('aria-label', `Open ${name} repository`);
            elements.gitRemote.hidden = false;
        } catch {
            if (request === remoteRequest) {
                hideRemoteLink();
            }
        }
    };

    const renderStatus = (status, cachedStatus = status) => {
        currentStatus = status;
        syncStatusWatch();
        void updateRemoteLink(status, projectRoot);
        elements.gitCompare.textContent = compareBranch || 'HEAD';
        statusCache.set(rootKey(projectRoot), { status: cachedStatus });
        elements.gitWorkspace.textContent = 'Worktrees';
        elements.gitWorkspace.dataset.tooltip = status.root || projectRoot;
        elements.gitBranch.textContent = status.branch || 'Detached HEAD';
        elements.gitAheadBehind.textContent = [
            status.ahead > 0 ? `↑${status.ahead}` : '',
            status.behind > 0 ? `↓${status.behind}` : '',
        ].filter(Boolean).join(' ');
        elements.gitFileList.replaceChildren();
        options.onStatusChange?.(status, projectRoot);

        if (status.clean || !status.files?.length) {
            showMessage('Working tree clean');
            return;
        }

        const fragment = document.createDocumentFragment();
        for (const file of status.files) {
            const row = document.createElement('div');
            row.className = 'notes-git-file-row';
            row.dataset.path = file.path;
            row.setAttribute('role', 'listitem');

            const open = document.createElement('button');
            open.type = 'button';
            open.className = 'notes-git-file';
            open.dataset.path = file.path;
            open.dataset.tooltip = file.originalPath
                ? `Open diff: ${file.originalPath} → ${file.path}`
                : `Open diff: ${file.path}`;
            open.setAttribute('aria-label', open.dataset.tooltip);
            open.textContent = file.path;

            const code = document.createElement('span');
            code.className = 'notes-git-status-code';
            code.textContent = file.compareOnly ? '≠' : statusCode(file);
            code.dataset.statusCode = statusCode(file);
            code.dataset.tooltip = file.compareOnly
                ? `Differs from ${compareBranch}`
                : statusTooltip(code.dataset.statusCode);
            code.setAttribute('aria-label', `Git status ${code.dataset.tooltip}`);

            const actions = document.createElement('span');
            actions.className = 'notes-git-file-actions';
            if (file.staged) {
                actions.appendChild(createAction('Staged', 'unstage', file.path, 'Click to unstage'));
            }
            if (file.unstaged || file.untracked) {
                const label = file.untracked ? 'Untracked' : 'Unstaged';
                const stageButton = createAction(label, 'stage', file.path, 'Click to stage');
                if (file.untracked) {
                    stageButton.classList.add('notes-git-action-untracked');
                }
                actions.appendChild(stageButton);
            }

            row.append(open, code, actions);
            fragment.appendChild(row);
        }
        elements.gitFileList.appendChild(fragment);
    };

    const refresh = async (force = true) => {
        const generation = ++requestGeneration;
        if (!projectRoot) {
            currentStatus = null;
            setAvailable(false);
            return false;
        }

        elements.gitRefresh.disabled = true;
        try {
            const status = force && api.refreshStatus
                ? await api.refreshStatus(projectRoot)
                : await api.getStatus(projectRoot);
            if (generation !== requestGeneration) {
                return false;
            }
            let displayStatus = status;
            if (compareBranch) {
                const branch = compareBranch;
                let files;
                try {
                    if (typeof api.getBranchDiffFiles !== 'function') {
                        throw new Error('Branch comparison is unavailable');
                    }
                    files = await api.getBranchDiffFiles(projectRoot, branch);
                } catch (error) {
                    if (generation !== requestGeneration || branch !== compareBranch) {
                        return false;
                    }
                    showMessage(`Could not compare with ${branch}: ${error?.message || error}`, 'error');
                    setAvailable(true);
                    return false;
                }
                if (generation !== requestGeneration || branch !== compareBranch) {
                    return false;
                }
                displayStatus = { ...status, files, clean: files.length === 0 };
            }
            renderStatus(displayStatus, status);
            setAvailable(true);
            return true;
        } catch (error) {
            if (generation !== requestGeneration) {
                return false;
            }
            currentStatus = null;
            void updateRemoteLink(null, projectRoot);
            statusCache.set(rootKey(projectRoot), { error });
            showMessage(String(error?.message || error || 'Git is unavailable'), 'error');
            setAvailable(false);
            options.onStatusChange?.(null, projectRoot);
            return false;
        } finally {
            if (generation === requestGeneration) {
                elements.gitRefresh.disabled = false;
            }
        }
    };

    const setProjectRoot = async (root) => {
        const nextRoot = String(root || '').trim();
        if (nextRoot === projectRoot && currentStatus) {
            return true;
        }
        projectRoot = nextRoot;
        compareBranch = '';
        elements.gitCompare.textContent = 'HEAD';
        currentStatus = null;
        syncStatusWatch();
        void updateRemoteLink(null, nextRoot);
        elements.gitFileList.replaceChildren();
        const cached = statusCache.get(rootKey(projectRoot));
        if (cached?.status) {
            renderStatus(cached.status);
            setAvailable(true);
        } else if (cached?.error) {
            currentStatus = null;
            showMessage(String(cached.error?.message || cached.error), 'error');
            setAvailable(false);
        }
        return refresh(false);
    };

    elements.gitWorkspace.addEventListener('click', () => {
        if (currentStatus) {
            options.showWorkspaceMenu?.(currentStatus.root || projectRoot, elements.gitWorkspace);
        }
    });

    const setCompareBranch = async (branch) => {
        const nextBranch = String(branch || '').trim();
        if (nextBranch === compareBranch) {
            return true;
        }
        compareBranch = nextBranch;
        elements.gitCompare.textContent = compareBranch || 'HEAD';
        return refresh(true);
    };

    elements.gitCompare.addEventListener('click', () => {
        if (currentStatus) {
            options.showCompareMenu?.(
                currentStatus.root || projectRoot,
                currentStatus.branch || '',
                compareBranch,
                elements.gitCompare,
                setCompareBranch,
            );
        }
    });

    elements.gitRefresh.addEventListener('click', () => {
        void refresh();
    });

    elements.gitRemote.addEventListener('click', () => {
        const url = elements.gitRemote.dataset.url;
        if (url) {
            void api.openRemote(url);
        }
    });

    elements.gitPane.addEventListener('mouseover', (event) => {
        const target = event.target.closest('[data-tooltip]');
        if (target && elements.gitPane.contains(target)) {
            options.showTooltip?.(target.dataset.tooltip);
        }
    });

    elements.gitPane.addEventListener('mouseout', (event) => {
        const target = event.target.closest('[data-tooltip]');
        if (target && elements.gitPane.contains(target) && !target.contains(event.relatedTarget)) {
            options.hideTooltip?.();
        }
    });

    const openFileDiff = async (file) => {
        if (!file) {
            return false;
        }
        const root = projectRoot;
        const generation = requestGeneration;
        try {
            const diff = compareBranch
                ? await api.getDiffAgainstBranch(root, file.path, '', compareBranch)
                : await api.getDiff(root, file.path, file.originalPath || '');
            if (generation !== requestGeneration || root !== projectRoot) {
                return false;
            }
            await options.onOpenDiff?.({ root, file, diff });
            return true;
        } catch (error) {
            showMessage(String(error?.message || error || `Could not read ${file.path}`), 'error');
            return false;
        }
    };

    elements.gitFileList.addEventListener('click', async (event) => {
        const action = event.target.closest('[data-git-action]');
        if (action && elements.gitFileList.contains(action)) {
            event.stopPropagation();
            const path = action.dataset.path;
            const operationName = action.dataset.gitAction;
            const operation = operationName === 'stage' ? api.stageFile : api.unstageFile;
            action.disabled = true;
            try {
                if (options.beforeFileAction
                    && !await options.beforeFileAction({ root: projectRoot, path, action: operationName })) {
                    return;
                }
                await operation(projectRoot, path);
                await refresh();
            } catch (error) {
                showMessage(String(error?.message || error || 'Git operation failed'), 'error');
            } finally {
                action.disabled = false;
            }
            return;
        }

        const fileButton = event.target.closest('.notes-git-file');
        if (!fileButton || !elements.gitFileList.contains(fileButton)) {
            return;
        }
        const file = currentStatus?.files?.find((item) => item.path === fileButton.dataset.path);
        if (!file) {
            return;
        }
        fileButton.disabled = true;
        try {
            await openFileDiff(file);
        } finally {
            fileButton.disabled = false;
        }
    });

    if (typeof api.eventsOn === 'function') {
        api.eventsOn('notesContextChanged', (context) => {
            void setProjectRoot(context?.projectRoot || '');
        });
        api.eventsOn('notesGitStatusChanged', (event) => {
            const root = rootKey(event?.root);
            if (!root) {
                return;
            }
            if (event?.error) {
                statusCache.set(root, { error: event.error });
                if (root === rootKey(projectRoot)) {
                    currentStatus = null;
                    void updateRemoteLink(null, projectRoot);
                    showMessage(String(event.error), 'error');
                    setAvailable(false);
                    options.onStatusChange?.(null, projectRoot);
                }
                return;
            }
            if (!event?.status) {
                return;
            }
            statusCache.set(root, { status: event.status });
            if (root === rootKey(projectRoot)) {
                if (compareBranch) {
                    void refresh(false);
                } else {
                    renderStatus(event.status);
                    setAvailable(true);
                }
            }
        });
    }

    void api.getCurrentProject().then(setProjectRoot).catch(() => setProjectRoot(''));

    return {
        refresh,
        setProjectRoot,
        setCompareBranch,
        getCompareBranch: () => compareBranch,
        setWatchActive(active) {
            watchActive = active === true;
            syncStatusWatch();
        },
        getCurrentStatus: () => currentStatus,
        getCurrentRoot: () => projectRoot,
        getCachedStatus(root = projectRoot) {
            return statusCache.get(rootKey(root))?.status || null;
        },
        getFileStatus(filePath, root = projectRoot) {
            const status = statusCache.get(rootKey(root))?.status;
            return status?.files?.find((file) => file.path === filePath) || null;
        },
        openFileDiff(filePath, root = projectRoot) {
            const status = statusCache.get(rootKey(root))?.status;
            const file = status?.files?.find((item) => item.path === filePath);
            return openFileDiff(file);
        },
    };
}

function createAction(label, action, path, tooltip) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = `notes-git-action notes-git-action-${action}`;
    button.dataset.gitAction = action;
    button.dataset.path = path;
    button.textContent = label;
    button.dataset.tooltip = tooltip;
    button.setAttribute('aria-label', tooltip);
    return button;
}
