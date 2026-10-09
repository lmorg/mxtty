function workspaceName(root) {
    const parts = String(root || '').replace(/[\\/]+$/, '').split(/[\\/]/);
    return parts.at(-1) || String(root || 'Git workspace');
}

function statusCode(file) {
    return String(file?.code || '??');
}

export function initNotesGitPanel(elements, api, options = {}) {
    if (!elements?.gitTab || !elements.gitPane || !elements.gitFileList) {
        return null;
    }

    let projectRoot = '';
    let requestGeneration = 0;
    let currentStatus = null;

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

    const renderStatus = (status) => {
        currentStatus = status;
        elements.gitWorkspace.textContent = workspaceName(status.root || projectRoot);
        elements.gitWorkspace.title = status.root || projectRoot;
        elements.gitBranch.textContent = status.branch || 'Detached HEAD';
        elements.gitAheadBehind.textContent = [
            status.ahead > 0 ? `↑${status.ahead}` : '',
            status.behind > 0 ? `↓${status.behind}` : '',
        ].filter(Boolean).join(' ');
        elements.gitFileList.replaceChildren();

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
            open.title = file.originalPath ? `${file.originalPath} → ${file.path}` : file.path;
            open.textContent = file.path;

            const code = document.createElement('span');
            code.className = 'notes-git-status-code';
            code.textContent = statusCode(file);
            code.setAttribute('aria-label', `Git status ${statusCode(file)}`);

            const actions = document.createElement('span');
            actions.className = 'notes-git-file-actions';
            if (file.staged) {
                actions.appendChild(createAction('Unstage', 'unstage', file.path));
            }
            if (file.unstaged || file.untracked) {
                actions.appendChild(createAction('Stage', 'stage', file.path));
            }

            row.append(open, code, actions);
            fragment.appendChild(row);
        }
        elements.gitFileList.appendChild(fragment);
    };

    const refresh = async () => {
        const generation = ++requestGeneration;
        if (!projectRoot) {
            currentStatus = null;
            setAvailable(false);
            return false;
        }

        elements.gitRefresh.disabled = true;
        try {
            const status = await api.getStatus(projectRoot);
            if (generation !== requestGeneration) {
                return false;
            }
            renderStatus(status);
            setAvailable(true);
            return true;
        } catch (error) {
            if (generation !== requestGeneration) {
                return false;
            }
            currentStatus = null;
            showMessage(String(error?.message || error || 'Git is unavailable'), 'error');
            setAvailable(false);
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
        currentStatus = null;
        elements.gitFileList.replaceChildren();
        return refresh();
    };

    elements.gitWorkspace.addEventListener('click', () => {
        if (currentStatus) {
            options.showWorkspaceMenu?.(currentStatus.root || projectRoot, elements.gitWorkspace);
        }
    });

    elements.gitRefresh.addEventListener('click', () => {
        void refresh();
    });

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
        const root = projectRoot;
        const generation = requestGeneration;
        fileButton.disabled = true;
        try {
            const diff = await api.getDiff(root, file.path, file.originalPath || '');
            if (generation !== requestGeneration || root !== projectRoot) {
                return;
            }
            await options.onOpenDiff?.({ root, file, diff });
        } catch (error) {
            showMessage(String(error?.message || error || `Could not read ${file.path}`), 'error');
        } finally {
            fileButton.disabled = false;
        }
    });

    if (typeof api.eventsOn === 'function') {
        api.eventsOn('notesContextChanged', (context) => {
            void setProjectRoot(context?.projectRoot || '');
        });
    }

    void api.getCurrentProject().then(setProjectRoot).catch(() => setProjectRoot(''));

    return {
        refresh,
        setProjectRoot,
        getCurrentStatus: () => currentStatus,
        getCurrentRoot: () => projectRoot,
    };
}

function createAction(label, action, path) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'notes-git-action';
    button.dataset.gitAction = action;
    button.dataset.path = path;
    button.textContent = label;
    button.title = `${label} ${path}`;
    return button;
}
