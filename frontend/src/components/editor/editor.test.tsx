import { normalizeLines, toCursorLocation } from "./Editor";
import type { RemoteUserPresence } from "./types";

type TestCase = {
	name: string;
	run: () => void;
};

const assertEqual = (received: unknown, expected: unknown) => {
	if (received !== expected) {
		throw new Error(`Expected ${expected as string}, got ${received as string}`);
	}
};

const assertDeepEqual = (received: unknown, expected: unknown) => {
	const receivedJson = JSON.stringify(received);
	const expectedJson = JSON.stringify(expected);
	if (receivedJson !== expectedJson) {
		throw new Error(`Expected ${expectedJson}, got ${receivedJson}`);
	}
};

const tests: TestCase[] = [
	{
		name: "normalizeLines returns a single empty line for empty input",
		run: () => {
			const lines = normalizeLines("");
			assertEqual(lines.length, 1);
			assertEqual(lines[0], "");
		},
	},
	{
		name: "normalizeLines splits unix and windows newlines consistently",
		run: () => {
			const unixLines = normalizeLines("alpha\nbeta");
			const windowsLines = normalizeLines("alpha\r\nbeta");
			assertDeepEqual(unixLines, ["alpha", "beta"]);
			assertDeepEqual(windowsLines, ["alpha", "beta"]);
		},
	},
	{
		name: "normalizeLines keeps unicode characters intact",
		run: () => {
			const unicodeText = "привіт\n😀 emoji";
			const lines = normalizeLines(unicodeText);
			assertDeepEqual(lines, ["привіт", "😀 emoji"]);
		},
	},
	{
		name: "toCursorLocation prefers explicit cursorLocation",
		run: () => {
			const user: RemoteUserPresence = {
				id: "1",
				name: "Pat",
				cursorLocation: { line: 2, column: 4 },
			};
			assertDeepEqual(toCursorLocation(["hello", "world", "third"], user), {
				line: 2,
				column: 4,
			});
		},
	},
	{
		name: "toCursorLocation translates absolute offsets into line/column",
		run: () => {
			const user: RemoteUserPresence = {
				id: "2",
				name: "Alex",
				cursorPosition: 7, // hel|lo\nw|orld
			};
			assertDeepEqual(toCursorLocation(["hello", "world", "third"], user), {
				line: 1,
				column: 1,
			});
		},
	},
	{
		name: "toCursorLocation caps positions to final line when offset exceeds content",
		run: () => {
			const user: RemoteUserPresence = {
				id: "3",
				name: "Taylor",
				cursorPosition: 200,
			};
			assertDeepEqual(toCursorLocation(["short"], user), { line: 0, column: 5 });
		},
	},
];

// Execute lightweight tests when this module is loaded (type-time safety only).
tests.forEach((testCase) => {
	try {
		testCase.run();
		// eslint-disable-next-line no-console
		console.info(`[editor.test] ✅ ${testCase.name}`);
	} catch (error) {
		// eslint-disable-next-line no-console
		console.error(`[editor.test] ❌ ${testCase.name}`, error);
		throw error;
	}
});
