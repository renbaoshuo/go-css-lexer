import fs from 'node:fs/promises';
import assert from "node:assert";

const source = `"a\\12\r\nb"`;
await fs.writeFile('tests/escaped-code-point/0016/source.css', source);
const writtenAs = await fs.readFile('tests/escaped-code-point/0016/source.css', 'utf8');

assert.strictEqual(writtenAs, source);
