import { isXMLDocument } from './swagger-utils.js';

function pathKey(path) {
    return JSON.stringify(path);
}

function decodeXmlEntities(value) {
    return String(value).replace(/&(#x[\da-f]+|#\d+|amp|lt|gt|quot|apos);/gi, (entity, token) => {
        const lower = token.toLowerCase();
        if (lower === 'amp') return '&';
        if (lower === 'lt') return '<';
        if (lower === 'gt') return '>';
        if (lower === 'quot') return '"';
        if (lower === 'apos') return "'";
        const codePoint = lower.startsWith('#x')
            ? Number.parseInt(lower.slice(2), 16)
            : Number.parseInt(lower.slice(1), 10);
        try {
            return String.fromCodePoint(codePoint);
        } catch {
            return entity;
        }
    });
}

function escapeXmlText(value) {
    return String(value ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
}

function escapeXmlAttribute(value, quote) {
    const escaped = escapeXmlText(value);
    return quote === "'" ? escaped.replace(/'/g, '&apos;') : escaped.replace(/"/g, '&quot;');
}

function findTagEnd(source, start) {
    let quote = '';
    for (let index = start + 1; index < source.length; index += 1) {
        const char = source[index];
        if (quote) {
            if (char === quote) quote = '';
        } else if (char === '"' || char === "'") {
            quote = char;
        } else if (char === '>') {
            return index;
        }
    }
    throw new Error('XML tag is not closed.');
}

function findDeclarationEnd(source, start) {
    let quote = '';
    let subsetDepth = 0;
    for (let index = start + 2; index < source.length; index += 1) {
        const char = source[index];
        if (quote) {
            if (char === quote) quote = '';
            continue;
        }
        if (char === '"' || char === "'") {
            quote = char;
        } else if (char === '[') {
            subsetDepth += 1;
        } else if (char === ']') {
            subsetDepth = Math.max(0, subsetDepth - 1);
        } else if (char === '>' && subsetDepth === 0) {
            return index + 1;
        }
    }
    throw new Error('XML declaration is not closed.');
}

function parseAttributes(source, tagStart, nameEnd, tagEnd) {
    const attributes = [];
    const zone = source.slice(nameEnd, tagEnd);
    const pattern = /(\s+)([^\s=/>]+)\s*=\s*(["'])([\s\S]*?)\3/g;
    let match;
    while ((match = pattern.exec(zone)) !== null) {
        const attrStart = nameEnd + match.index;
        const nameStart = attrStart + match[1].length;
        const nameEndOffset = nameStart + match[2].length;
        const quoteOffset = match[0].indexOf(match[3], match[1].length + match[2].length);
        const valueStart = attrStart + quoteOffset + 1;
        attributes.push({
            name: match[2],
            nameStart,
            nameEnd: nameEndOffset,
            valueStart,
            valueEnd: valueStart + match[4].length,
            fullStart: attrStart,
            fullEnd: nameEnd + match.index + match[0].length,
            quote: match[3],
            value: decodeXmlEntities(match[4]),
        });
    }
    return attributes;
}

function parseXmlTokens(source) {
    const documentNode = { type: 'document', children: [] };
    const stack = [documentNode];
    let index = 0;

    while (index < source.length) {
        if (source[index] !== '<') {
            const end = source.indexOf('<', index);
            const textEnd = end < 0 ? source.length : end;
            if (textEnd > index) {
                stack[stack.length - 1].children.push({
                    type: 'text',
                    start: index,
                    end: textEnd,
                    valueStart: index,
                    valueEnd: textEnd,
                    value: decodeXmlEntities(source.slice(index, textEnd)),
                    cdata: false,
                });
            }
            index = textEnd;
            continue;
        }

        if (source.startsWith('<!--', index)) {
            const end = source.indexOf('-->', index + 4);
            if (end < 0) throw new Error('XML comment is not closed.');
            index = end + 3;
            continue;
        }

        if (source.startsWith('<![CDATA[', index)) {
            const end = source.indexOf(']]>', index + 9);
            if (end < 0) throw new Error('XML CDATA section is not closed.');
            stack[stack.length - 1].children.push({
                type: 'text',
                start: index,
                end: end + 3,
                valueStart: index + 9,
                valueEnd: end,
                value: source.slice(index + 9, end),
                cdata: true,
            });
            index = end + 3;
            continue;
        }

        if (source.startsWith('<?', index)) {
            const end = source.indexOf('?>', index + 2);
            if (end < 0) throw new Error('XML processing instruction is not closed.');
            index = end + 2;
            continue;
        }

        if (source.startsWith('<!', index)) {
            index = findDeclarationEnd(source, index);
            continue;
        }

        const tagEnd = findTagEnd(source, index);
        if (source.startsWith('</', index)) {
            const closing = source.slice(index, tagEnd + 1).match(/^<\/\s*([^\s>]+)\s*>$/);
            if (!closing || stack.length === 1) throw new Error('Unexpected XML closing tag.');
            const current = stack.pop();
            if (current.name !== closing[1]) {
                throw new Error(`Mismatched XML closing tag: expected </${current.name}>.`);
            }
            const nameOffset = source.slice(index, tagEnd + 1).indexOf(closing[1]);
            current.closeNameStart = index + nameOffset;
            current.closeNameEnd = current.closeNameStart + closing[1].length;
            current.end = tagEnd + 1;
            current.endTagStart = index;
            index = tagEnd + 1;
            continue;
        }

        const opening = source.slice(index, tagEnd + 1).match(/^<\s*([^\s/>]+)/);
        if (!opening) throw new Error('Invalid XML opening tag.');
        const name = opening[1];
        const nameOffset = source.slice(index, tagEnd + 1).indexOf(name);
        const nameStart = index + nameOffset;
        const openEnd = tagEnd + 1;
        const selfClosing = /\/\s*>$/.test(source.slice(index, openEnd));
        const node = {
            type: 'element',
            name,
            start: index,
            openEnd,
            end: selfClosing ? openEnd : null,
            endTagStart: selfClosing ? openEnd : null,
            nameStart,
            nameEnd: nameStart + name.length,
            selfClosing,
            attributes: parseAttributes(source, index, nameStart + name.length, tagEnd),
            children: [],
        };
        stack[stack.length - 1].children.push(node);
        if (!selfClosing) stack.push(node);
        index = openEnd;
    }

    if (stack.length !== 1) {
        throw new Error(`XML element <${stack[stack.length - 1].name}> is not closed.`);
    }

    const roots = documentNode.children.filter((child) => child.type === 'element');
    if (roots.length !== 1) throw new Error('XML must contain exactly one root element.');
    return { documentNode, root: roots[0] };
}

function validXmlName(name) {
    return /^[A-Za-z_][A-Za-z\d_.:-]*$/.test(String(name));
}

function buildXmlValue(node, path, index) {
    index.elements.set(pathKey(path), node);
    const elements = node.children.filter((child) => child.type === 'element');
    const textNodes = node.children.filter((child) => child.type === 'text');
    const visibleTextNodes = elements.length > 0
        ? textNodes.filter((child) => child.value.trim() !== '')
        : textNodes;
    const hasAttributes = node.attributes.length > 0;

    if (!hasAttributes && elements.length === 0) {
        if (visibleTextNodes.length === 0) {
            index.nodes.set(pathKey(path), { type: 'empty-element', node });
            return '';
        }
        const value = visibleTextNodes.map((child) => child.value).join('');
        index.nodes.set(pathKey(path), visibleTextNodes.length === 1
            ? { type: 'text', node: visibleTextNodes[0] }
            : { type: 'mixed-text', nodes: visibleTextNodes });
        return value;
    }

    const result = {};
    index.nodes.set(pathKey(path), { type: 'element', node });

    for (const attribute of node.attributes) {
        const key = `@${attribute.name}`;
        result[key] = attribute.value;
        index.nodes.set(pathKey([...path, key]), { type: 'attribute', node: attribute });
    }

    const groups = new Map();
    for (const child of elements) {
        if (!groups.has(child.name)) groups.set(child.name, []);
        groups.get(child.name).push(child);
    }
    for (const [name, children] of groups) {
        const childPath = [...path, name];
        if (children.length === 1) {
            result[name] = buildXmlValue(children[0], childPath, index);
            continue;
        }
        const array = [];
        index.arrays.set(pathKey(childPath), { parent: node, name, nodes: children });
        children.forEach((child, childIndex) => {
            array.push(buildXmlValue(child, [...childPath, childIndex], index));
        });
        result[name] = array;
    }

    if (visibleTextNodes.length > 0) {
        const textPath = [...path, '#text'];
        result['#text'] = visibleTextNodes.map((child) => child.value).join('');
        index.nodes.set(pathKey(textPath), visibleTextNodes.length === 1
            ? { type: 'text', node: visibleTextNodes[0] }
            : { type: 'mixed-text', nodes: visibleTextNodes });
    }

    return result;
}

function parseXml(source) {
    const text = String(source ?? '');
    if (typeof DOMParser !== 'undefined') {
        const parsed = new DOMParser().parseFromString(text, 'application/xml');
        if (parsed.querySelector('parsererror')) throw new Error('Unable to parse XML document.');
    }
    const parsed = parseXmlTokens(text);
    const index = { nodes: new Map(), elements: new Map(), arrays: new Map(), root: parsed.root };
    const value = { [parsed.root.name]: buildXmlValue(parsed.root, [parsed.root.name], index) };
    return { value, index };
}

function replaceRanges(source, replacements) {
    let result = source;
    for (const replacement of [...replacements].sort((a, b) => b.start - a.start)) {
        result = result.slice(0, replacement.start) + replacement.text + result.slice(replacement.end);
    }
    return result;
}

function getPathNode(parsed, path) {
    return parsed.index.nodes.get(pathKey(path));
}

function setElementText(source, element, value) {
    const text = escapeXmlText(value);
    if (element.selfClosing) {
        const opening = source.slice(element.start, element.openEnd).replace(/\/\s*>$/, '>');
        return `${source.slice(0, element.start)}${opening}${text}</${element.name}>${source.slice(element.openEnd)}`;
    }
    return source.slice(0, element.openEnd) + text + source.slice(element.openEnd);
}

function insertXmlChild(source, parent, name, value) {
    if (!validXmlName(name)) throw new Error(`Invalid XML element name: ${name}`);
    const childXml = `<${name}>${escapeXmlText(value)}</${name}>`;
    if (parent.selfClosing) {
        const opening = source.slice(parent.start, parent.openEnd).replace(/\/\s*>$/, '>');
        return `${source.slice(0, parent.start)}${opening}${childXml}</${parent.name}>${source.slice(parent.openEnd)}`;
    }
    return `${source.slice(0, parent.endTagStart)}${childXml}${source.slice(parent.endTagStart)}`;
}

function setXmlText(source, node, value) {
    const text = String(value ?? '');
    if (node.cdata) {
        if (text.includes(']]>')) throw new Error('XML text cannot contain ]]> inside a CDATA section.');
        return source.slice(0, node.valueStart) + text + source.slice(node.valueEnd);
    }
    return source.slice(0, node.valueStart) + escapeXmlText(text) + source.slice(node.valueEnd);
}

function insertXmlAttribute(source, parent, name, value) {
    if (!validXmlName(name)) throw new Error(`Invalid XML attribute name: ${name}`);
    if (parent.attributes.some((attribute) => attribute.name === name)) throw new Error(`XML attribute "${name}" already exists.`);
    const at = parent.selfClosing ? source.slice(parent.start, parent.openEnd).search(/\/\s*>$/) + parent.start : parent.openEnd - 1;
    return source.slice(0, at) + ` ${name}="${escapeXmlAttribute(value, '"')}"` + source.slice(at);
}

function makeXmlEditor() {
    return {
        id: 'xml',
        match(fileName, source) {
            return /\.xml$/i.test(fileName || '') || isXMLDocument(source);
        },
        parse(source) {
            return parseXml(source).value;
        },
        setValue(source, path, value) {
            const target = getPathNode(parseXml(source), path);
            if (!target) throw new Error('Unable to locate XML value.');
            if (target.type === 'attribute') {
                const attr = target.node;
                return source.slice(0, attr.valueStart) + escapeXmlAttribute(value, attr.quote) + source.slice(attr.valueEnd);
            }
            if (target.type === 'text') return setXmlText(source, target.node, value);
            if (target.type === 'mixed-text') throw new Error('Cannot edit mixed XML text while preserving child elements.');
            if (target.type === 'empty-element') {
                const opening = source.slice(target.node.start, target.node.openEnd).replace(/\/\s*>$/, '>');
                const replacement = `${opening}${escapeXmlText(value)}</${target.node.name}>`;
                return source.slice(0, target.node.start) + replacement + source.slice(target.node.openEnd);
            }
            throw new Error('Only XML text and attributes can be edited as values.');
        },
        renameKey(source, path, nextKey) {
            const parsed = parseXml(source);
            const target = getPathNode(parsed, path);
            if (target?.type === 'attribute') {
                const attrName = String(nextKey).replace(/^@/, '');
                if (!validXmlName(attrName)) throw new Error(`Invalid XML attribute name: ${nextKey}`);
                return replaceRanges(source, [{ start: target.node.nameStart, end: target.node.nameEnd, text: attrName }]);
            }
            const node = parsed.index.elements.get(pathKey(path));
            if (node) {
                if (String(nextKey).startsWith('@')) {
                    const attrName = String(nextKey).slice(1);
                    if (!validXmlName(attrName)) throw new Error(`Invalid XML attribute name: ${nextKey}`);
                    const parentPath = path.slice(0, -1);
                    const parent = parsed.index.elements.get(pathKey(parentPath));
                    if (!parent) throw new Error('Cannot add an attribute outside an XML element.');
                    const content = node.selfClosing ? '' : source.slice(node.openEnd, node.endTagStart);
                    if (content.includes('<')) throw new Error('Only text-only elements can be converted to XML attributes.');
                    const sourceWithoutElement = source.slice(0, node.start) + source.slice(node.end);
                    const parentWithoutElement = parseXml(sourceWithoutElement).index.elements.get(pathKey(parentPath));
                    return insertXmlAttribute(sourceWithoutElement, parentWithoutElement, attrName, decodeXmlEntities(content));
                }
                if (!validXmlName(nextKey)) throw new Error(`Invalid XML element name: ${nextKey}`);
                const replacements = [{ start: node.nameStart, end: node.nameEnd, text: nextKey }];
                if (!node.selfClosing) replacements.push({ start: node.closeNameStart, end: node.closeNameEnd, text: nextKey });
                return replaceRanges(source, replacements);
            }
            throw new Error('Only XML elements and attributes can be renamed.');
        },
        addKey(source, parentPath, key, value) {
            const parsed = parseXml(source);
            const parent = parsed.index.elements.get(pathKey(parentPath));
            if (!parent || parent.type !== 'element') throw new Error('Cannot add a sibling XML root element.');
            const name = String(key);
            if (name.startsWith('@')) {
                return insertXmlAttribute(source, parent, name.slice(1), value);
            }
            if (name === '#text') {
                if (parent.children.some((child) => child.type === 'text' && child.value.trim() !== '')) {
                    throw new Error('XML text already exists at this element.');
                }
                return setElementText(source, parent, value);
            }
            return insertXmlChild(source, parent, name, value);
        },
        addItem(source, arrayPath, value) {
            const parsed = parseXml(source);
            const array = parsed.index.arrays.get(pathKey(arrayPath));
            if (!array) throw new Error('XML repeated-element list not found.');
            return insertXmlChild(source, array.parent, array.name, value);
        },
        deleteNode(source, path) {
            const parsed = parseXml(source);
            const target = getPathNode(parsed, path);
            if (target?.type === 'attribute') {
                const attr = target.node;
                return source.slice(0, attr.fullStart) + source.slice(attr.fullEnd);
            }
            const element = parsed.index.elements.get(pathKey(path));
            if (element && path[path.length - 1] !== '#text') {
                if (element === parsed.index.root) throw new Error('Cannot delete the XML root element.');
                return source.slice(0, element.start) + source.slice(element.end);
            }
            if (target?.type === 'text') {
                return source.slice(0, target.node.start) + source.slice(target.node.end);
            }
            if (target?.type === 'mixed-text') {
                return replaceRanges(source, target.nodes.map((node) => ({ start: node.start, end: node.end, text: '' })));
            }
            const array = parsed.index.arrays.get(pathKey(path));
            if (array) {
                return replaceRanges(source, array.nodes.map((node) => ({ start: node.start, end: node.end, text: '' })));
            }
            throw new Error('Unable to locate XML node to delete.');
        },
    };
}

export const xmlEditor = makeXmlEditor();
