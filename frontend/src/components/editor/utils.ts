import type { CursorLocation, RemoteUserPresence } from "./types";

export const normalizeLines = (value: string | string[]): string[] => {
	if (Array.isArray(value)) {
		return value.length > 0 ? value : [""];
	}

	const safeValue = value ?? "";
	const rows = safeValue.replace(/\r\n/g, "\n").split("\n");
	return rows.length > 0 ? rows : [""];
};

export const toCursorLocation = (
	lines: string[],
	user: RemoteUserPresence,
): CursorLocation | undefined => {
	if (user.cursorLocation) {
		return user.cursorLocation;
	}
	if (typeof user.cursorPosition !== "number") {
		return undefined;
	}

	if (lines.length === 0) {
		return { line: 0, column: 0 };
	}

	const totalChars =
		lines.reduce((sum, line) => sum + line.length, 0) +
		Math.max(0, lines.length - 1);

	const clampedPosition = Math.max(
		0,
		Math.min(user.cursorPosition, totalChars),
	);

	let remaining = clampedPosition;
	for (let lineIndex = 0; lineIndex < lines.length; lineIndex += 1) {
		const length = lines[lineIndex]?.length ?? 0;
		if (remaining <= length) {
			return { line: lineIndex, column: remaining };
		}
		remaining -= length + 1;
	}

	const lastLineIndex = lines.length - 1;
	const lastLineLength = lines[lastLineIndex]?.length ?? 0;
	return {
		line: lastLineIndex,
		column: Math.min(remaining, lastLineLength),
	};
};
