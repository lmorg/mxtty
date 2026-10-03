export function showFullscreenImageOverlay(options = {}) {
    const {
        dataURL,
        svgElement,
        sourceWidth = 0,
        sourceHeight = 0,
        onOpen,
        onClose,
    } = options;

    if (!dataURL && !svgElement) {
        return;
    }

    const existingOverlay = document.getElementById('fullscreen-image-overlay');
    if (existingOverlay) {
        existingOverlay.remove();
    }

    const overlay = document.createElement('div');
    overlay.id = 'fullscreen-image-overlay';
    overlay.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.95);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 999999;
        overflow: auto;
        padding: 20px;
        box-sizing: border-box;
    `;

    const container = document.createElement('div');
    container.style.cssText = `
        display: flex;
        align-items: center;
        justify-content: center;
        max-width: 100%;
        max-height: 100%;
    `;

    let contentElement;
    if (svgElement) {
        contentElement = svgElement.cloneNode(true);
        const viewBox = String(contentElement.getAttribute('viewBox') || '')
            .trim()
            .split(/[\s,]+/)
            .map(Number);
        const width = sourceWidth || viewBox[2] || Number.parseFloat(contentElement.getAttribute('width')) || 800;
        const height = sourceHeight || viewBox[3] || Number.parseFloat(contentElement.getAttribute('height')) || 600;
        if (!contentElement.hasAttribute('viewBox')) {
            contentElement.setAttribute('viewBox', `0 0 ${width} ${height}`);
        }
        contentElement.setAttribute('width', String(width));
        contentElement.setAttribute('height', String(height));
        contentElement.style.cssText = `
            display: block;
            width: ${width}px;
            height: ${height}px;
            max-width: calc(100vw - 40px);
            max-height: calc(100vh - 40px);
            object-fit: contain;
            box-shadow: 0 0 30px rgba(255, 255, 255, 0.3);
            border-radius: 8px;
            background: transparent;
        `;
        // Ensure SVG is visible
        contentElement.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    } else {
        // Use image element for data URLs
        contentElement = document.createElement('img');
        contentElement.src = dataURL;
        contentElement.style.cssText = `
            max-width: 100%;
            max-height: 100%;
            object-fit: contain;
            box-shadow: 0 0 30px rgba(255, 255, 255, 0.3);
            border-radius: 8px;
        `;
    }

    const info = document.createElement('div');
    info.style.cssText = `
        position: absolute;
        bottom: 20px;
        right: 20px;
        color: rgba(255, 255, 255, 0.7);
        font-size: 12px;
        font-family: var(--terminal-menu-font);
        background: rgba(0, 0, 0, 0.8);
        padding: 8px 12px;
        border-radius: 4px;
    `;
    info.textContent = `${sourceWidth}×${sourceHeight} | Press ESC to close`;

    container.appendChild(contentElement);
    overlay.appendChild(container);
    overlay.appendChild(info);
    document.body.appendChild(overlay);

    if (typeof onOpen === 'function') {
        onOpen();
    }

    const closeOverlay = () => {
        document.removeEventListener('keydown', handleKey);
        overlay.removeEventListener('click', handleClick);
        overlay.remove();
        if (typeof onClose === 'function') {
            onClose();
        }
    };

    const handleKey = (e) => {
        if (e.key === 'Escape') {
            e.stopPropagation();
            e.preventDefault();
            closeOverlay();
        }
    };

    const handleClick = () => {
        closeOverlay();
    };

    document.addEventListener('keydown', handleKey);
    overlay.addEventListener('click', handleClick);
}