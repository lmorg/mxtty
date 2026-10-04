import { describe, expect, it } from 'vitest';
import { getStructuredEditor, xmlEditor } from './structured-editors.js';

const source = '<?xml version="1.0"?>\n<root id="old">\n  <!-- keep this comment -->\n  <name>Old &amp; valid</name>\n  <item>A</item>\n  <item>B</item>\n</root>\n';

describe('xmlEditor', () => {
    it('is selected for XML files and projects elements, attributes, text, and repeated elements', () => {
        expect(getStructuredEditor('settings.xml')).toBe(xmlEditor);
        expect(getStructuredEditor('SETTINGS.XML')).toBe(xmlEditor);
        const bomXml = '\uFEFF<?xml version="1.0"?><root/>';
        expect(getStructuredEditor('settings.txt', bomXml)).toBe(xmlEditor);
        expect(xmlEditor.parse(bomXml)).toEqual({ root: '' });
        expect(xmlEditor.parse(source)).toEqual({
            root: {
                '@id': 'old',
                name: 'Old & valid',
                item: ['A', 'B'],
            },
        });
    });

    it('edits only the targeted text or attribute value and preserves comments and layout', () => {
        const textUpdated = xmlEditor.setValue(source, ['root', 'name'], 'New <value>');
        expect(textUpdated).toContain('<name>New &lt;value&gt;</name>');
        expect(textUpdated).toContain('<!-- keep this comment -->');
        expect(textUpdated).toContain('<item>A</item>\n  <item>B</item>');

        const attributeUpdated = xmlEditor.setValue(source, ['root', '@id'], 'next&value');
        expect(attributeUpdated).toContain('id="next&amp;value"');
        expect(attributeUpdated).toContain('<name>Old &amp; valid</name>');
    });

    it('renames an element and its closing tag without reformatting siblings', () => {
        const result = xmlEditor.renameKey(source, ['root', 'name'], 'title');
        expect(result).toContain('<title>Old &amp; valid</title>');
        expect(result).toContain('<item>A</item>');
        expect(result).toContain('<!-- keep this comment -->');
    });

    it('adds an attribute, child element, and repeated child element', () => {
        const withAttribute = xmlEditor.addKey(source, ['root'], '@enabled', 'true');
        expect(withAttribute).toContain('id="old" enabled="true"');

        const placeholder = xmlEditor.addKey(source, ['root'], 'newKey', '');
        const renamedToAttribute = xmlEditor.renameKey(placeholder, ['root', 'newKey'], '@enabled');
        expect(xmlEditor.parse(renamedToAttribute).root['@enabled']).toBe('');
        expect(renamedToAttribute).not.toContain('<newKey>');
        expect(renamedToAttribute).toContain('<!-- keep this comment -->');

        const withChild = xmlEditor.addKey(source, ['root'], 'newChild', 'value');
        expect(withChild).toContain('<newChild>value</newChild>');
        expect(withChild).toContain('<!-- keep this comment -->');

        const withRepeatedChild = xmlEditor.addItem(source, ['root', 'item'], 'C');
        expect(xmlEditor.parse(withRepeatedChild).root.item).toEqual(['A', 'B', 'C']);
    });

    it('deletes an attribute, child, or one repeated child without touching siblings', () => {
        const withoutAttribute = xmlEditor.deleteNode(source, ['root', '@id']);
        expect(withoutAttribute).not.toContain('id="old"');

        const withoutChild = xmlEditor.deleteNode(source, ['root', 'name']);
        expect(withoutChild).not.toContain('<name>');
        expect(withoutChild).toContain('<item>A</item>');

        const withoutRepeatedChild = xmlEditor.deleteNode(source, ['root', 'item', 0]);
        expect(xmlEditor.parse(withoutRepeatedChild).root.item).toBe('B');
    });

    it('rejects malformed XML and unsafe edits', () => {
        expect(() => xmlEditor.parse('<root><item></root>')).toThrow();
        expect(() => xmlEditor.renameKey(source, ['root', 'name'], 'bad name')).toThrow();
        expect(() => xmlEditor.deleteNode(source, ['root'])).toThrow();
    });
});
