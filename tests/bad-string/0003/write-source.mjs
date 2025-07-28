import fs from 'node:fs/promises';
import assert from "node:assert";

const source = `"foo\r\n"\n`;
await fs.writeFile('tests/bad-string/0003/source.css', source);
const writtenAs = await fs.readFile('tests/bad-string/0003/source.css', 'utf8');

assert.strictEqual(writtenAs, source);
