import {
	type ClipboardEvent,
	type FormEvent,
	type KeyboardEvent,
	type UIEvent,
	useCallback,
	useEffect,
	useMemo,
	useRef,
	useState,
} from "react";
import { useTheme } from "../../contexts/ThemeContext";
import type {
	CursorLocation,
	EditorProps,
	RemoteUserPresence,
} from "./types";
import "./editor.css";

const ZERO_WIDTH_SPACE = "\u200b";
const computeAbsoluteOffset = (
	lineIndex: number,
	caretOffset: number,
	currentLines: string[],
) => {
	let offset = caretOffset;
	for (let i = 0; i < lineIndex; i += 1) {
		offset += (currentLines[i]?.length ?? 0) + 1;
	}
	return offset;
};

export const normalizeLines = (value: string | string[]): string[] => {
	if (Array.isArray(value)) {
		return value.length > 0 ? value : [""];
	}

	const safeValue = value ?? "";
	const rows = safeValue.replace(/\r\n/g, "\n").split("\n");
	return rows.length > 0 ? rows : [""];
};

const getActiveLineIndex = (container: HTMLElement): number | null => {
	const selection = window.getSelection();
	if (!selection?.anchorNode) {
		return null;
	}

	const anchorElement =
		selection.anchorNode instanceof HTMLElement
			? selection.anchorNode
			: selection.anchorNode?.parentElement;
	if (!anchorElement) {
		return null;
	}

	const lineElement = anchorElement.closest<HTMLElement>(".editor-line");
	if (!lineElement || !container.contains(lineElement)) {
		return null;
	}

	const index = Number.parseInt(
		lineElement.dataset.lineIndex ?? lineElement.getAttribute("data-line-index") ?? "",
		10,
	);
	return Number.isNaN(index) ? null : index;
};

const caretOffsetWithin = (element: HTMLElement): number => {
	const selection = window.getSelection();
	if (!selection || selection.rangeCount === 0) {
		return element.innerText.replaceAll(ZERO_WIDTH_SPACE, "").length;
	}

	const range = selection.getRangeAt(0);
	const preSelectionRange = range.cloneRange();
	preSelectionRange.selectNodeContents(element);
	preSelectionRange.setEnd(range.endContainer, range.endOffset);
	return preSelectionRange.toString().replaceAll(ZERO_WIDTH_SPACE, "").length;
};

const moveCaretTo = (element: HTMLElement, offset: number) => {
	const selection = window.getSelection();
	if (!selection) {
		return;
	}

	const targetNode = element.firstChild ?? element;
	const textLength = element.textContent?.length ?? 0;
	const clampedOffset = Math.max(0, Math.min(offset, textLength));

	const range = document.createRange();
	try {
		range.setStart(targetNode, clampedOffset);
	} catch {
		range.setStart(element, element.childNodes.length);
	}
	range.collapse(true);
	selection.removeAllRanges();
	selection.addRange(range);
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

	let remaining = user.cursorPosition;
	for (let lineIndex = 0; lineIndex < lines.length; lineIndex += 1) {
		const length = lines[lineIndex]?.length ?? 0;
		if (remaining <= length) {
			return { line: lineIndex, column: remaining };
		}
		remaining -= length + 1; // +1 for newline
	}

	return {
		line: lines.length - 1,
		column: lines[lines.length - 1]?.length ?? 0,
	};
};

const lineKey = (index: number) => index;

export const Editor = ({
	value,
	defaultValue = "",
	onChange,
	onCursorChange,
	remoteUsers = [],
	isConnected = false,
	showPresence = true,
	className = "",
	ariaLabel = "Collaborative editor",
	colorScheme,
}: EditorProps) => {
	const { theme } = useTheme();
	const scheme = colorScheme ?? theme ?? "light";

	const [lines, setLines] = useState<string[]>(() =>
		normalizeLines(value ?? defaultValue),
	);
	const [activeLine, setActiveLine] = useState(0);
	const [presenceVisible, setPresenceVisible] = useState(showPresence);

	const contentRef = useRef<HTMLDivElement | null>(null);
	const lineNumberRef = useRef<HTMLDivElement | null>(null);
	const lastEmittedValue = useRef(lines.join("\n"));
	const linesRef = useRef(lines);
	const lastCursorOffset = useRef<number | null>(null);

	useEffect(() => {
		setPresenceVisible(showPresence);
	}, [showPresence]);

	// keep external value in sync
	useEffect(() => {
		if (value === undefined) {
			return;
		}
		const nextLines = normalizeLines(value);
		const joined = nextLines.join("\n");
		if (joined !== lastEmittedValue.current) {
			setLines(nextLines);
			lastEmittedValue.current = joined;
		}
	}, [value]);

	useEffect(() => {
		linesRef.current = lines;
	}, [lines]);

	const syncSelectionState = useCallback(() => {
		if (!contentRef.current) {
			return;
		}
		const next = getActiveLineIndex(contentRef.current);
		if (next === null) {
			return;
		}
		const clamped = Math.min(next, Math.max(linesRef.current.length - 1, 0));
		setActiveLine(clamped);

		if (!onCursorChange) {
			return;
		}

		const lineElement = contentRef.current.querySelector<HTMLElement>(
			`.editor-line[data-line-index="${clamped}"]`,
		);
		const caret = lineElement ? caretOffsetWithin(lineElement) : 0;
		const absolute = computeAbsoluteOffset(clamped, caret, linesRef.current);
		if (lastCursorOffset.current !== absolute) {
			lastCursorOffset.current = absolute;
			onCursorChange(absolute);
		}
	}, [onCursorChange]);

	useEffect(() => {
		const lastIndex = Math.max(lines.length - 1, 0);
		if (activeLine > lastIndex) {
			setActiveLine(lastIndex);
		}
	}, [activeLine, lines.length]);

	useEffect(() => {
		document.addEventListener("selectionchange", syncSelectionState);
		return () => document.removeEventListener("selectionchange", syncSelectionState);
	}, [syncSelectionState]);

	useEffect(() => {
		if (!contentRef.current || !lineNumberRef.current) {
			return;
		}
		lineNumberRef.current.scrollTop = contentRef.current.scrollTop;
	}, [lines.length]);

	useEffect(() => {
		requestAnimationFrame(syncSelectionState);
	}, [lines.length, syncSelectionState]);

	const emitChange = useCallback(
		(nextLines: string[]) => {
			const joined = nextLines.join("\n");
			lastEmittedValue.current = joined;
			onChange?.(joined);
		},
		[onChange],
	);

	const updateLines = useCallback(
		(updater: (prev: string[]) => string[]) => {
			setLines((prev) => {
				const next = updater(prev);
				emitChange(next);
				return next;
			});
		},
		[emitChange],
	);

	const focusLine = useCallback((index: number, offset = 0) => {
		const target = contentRef.current?.querySelector<HTMLElement>(
			`.editor-line[data-line-index="${index}"]`,
		);
		if (!target) {
			return;
		}
		target.focus();
		requestAnimationFrame(() => moveCaretTo(target, offset));
	}, []);

	const handleInput = useCallback(
		(index: number, event: FormEvent<HTMLDivElement>) => {
			const caret = caretOffsetWithin(event.currentTarget);
			const text = (event.currentTarget.innerText ?? "")
				.replace(/\n/g, "")
				.replaceAll(ZERO_WIDTH_SPACE, "");
			updateLines((prev) => {
				const next = [...prev];
				next[index] = text;
				return normalizeLines(next);
			});
			setActiveLine(index);
			requestAnimationFrame(() => {
				moveCaretTo(event.currentTarget, caret);
				requestAnimationFrame(syncSelectionState);
			});
		},
		[syncSelectionState, updateLines],
	);

	const handleEnter = useCallback(
		(index: number, event: KeyboardEvent<HTMLDivElement>) => {
			event.preventDefault();
			const target = event.currentTarget;
			const text = (target.innerText ?? "").replace(/\n/g, "").replaceAll(ZERO_WIDTH_SPACE, "");
			const caret = caretOffsetWithin(target);

			updateLines((prev) => {
				const next = [...prev];
				const before = text.slice(0, caret);
				const after = text.slice(caret);
				next[index] = before;
				next.splice(index + 1, 0, after);
				return normalizeLines(next);
			});

			requestAnimationFrame(() => {
				focusLine(index + 1, 0);
				requestAnimationFrame(syncSelectionState);
			});
		},
		[focusLine, syncSelectionState, updateLines],
	);

	const handleBackspace = useCallback(
		(index: number, event: KeyboardEvent<HTMLDivElement>) => {
			const target = event.currentTarget;
			const caret = caretOffsetWithin(target);
			if (caret !== 0 || index === 0) {
				return;
			}

			event.preventDefault();
			let nextCaretOffset = 0;
			updateLines((prev) => {
				const next = [...prev];
				const previous = next[index - 1] ?? "";
				const current = next[index] ?? "";
				nextCaretOffset = previous.length;
				next[index - 1] = `${previous}${current}`;
				next.splice(index, 1);
				return normalizeLines(next);
			});
			requestAnimationFrame(() => {
				focusLine(index - 1, nextCaretOffset);
				requestAnimationFrame(syncSelectionState);
			});
		},
		[focusLine, syncSelectionState, updateLines],
	);

	const handlePaste = useCallback(
		(index: number, event: ClipboardEvent<HTMLDivElement>) => {
			event.preventDefault();
			const raw = event.clipboardData?.getData("text/plain") ?? "";
			const target = event.currentTarget;
			const caret = caretOffsetWithin(target);
			const clean = raw.replaceAll(ZERO_WIDTH_SPACE, "");
			const parts = normalizeLines(clean);
			const current = (target.innerText ?? "").replace(/\n/g, "").replaceAll(ZERO_WIDTH_SPACE, "");
			const before = current.slice(0, caret);
			const after = current.slice(caret);

			updateLines((prev) => {
				const next = [...prev];
				if (parts.length === 1) {
					next[index] = `${before}${parts[0]}${after}`;
					return normalizeLines(next);
				}

				const head = `${before}${parts[0]}`;
				const tail = `${parts[parts.length - 1]}${after}`;
				const middle = parts.slice(1, -1);

				next[index] = head;
				next.splice(index + 1, 0, ...middle, tail);
				return normalizeLines(next);
			});

			requestAnimationFrame(() => {
				focusLine(index + parts.length - 1, parts[parts.length - 1].length);
				requestAnimationFrame(syncSelectionState);
			});
		},
		[focusLine, syncSelectionState, updateLines],
	);

	const handleDelete = useCallback(
		(index: number, event: KeyboardEvent<HTMLDivElement>) => {
			const target = event.currentTarget;
			const text = (target.innerText ?? "").replace(/\n/g, "").replaceAll(ZERO_WIDTH_SPACE, "");
			const caret = caretOffsetWithin(target);
			const hasNextLine = Boolean(target.nextElementSibling);
			if (!hasNextLine || caret !== text.length) {
				return;
			}

			event.preventDefault();
			updateLines((prev) => {
				const next = [...prev];
				const current = next[index] ?? "";
				const following = next[index + 1] ?? "";
				next[index] = `${current}${following}`;
				next.splice(index + 1, 1);
				return normalizeLines(next);
			});

			requestAnimationFrame(() => {
				focusLine(index, caret);
				requestAnimationFrame(syncSelectionState);
			});
		},
		[focusLine, syncSelectionState, updateLines],
	);

	const handleKeyDown = useCallback(
		(index: number, event: KeyboardEvent<HTMLDivElement>) => {
			if (event.key === "Enter") {
				handleEnter(index, event);
				return;
			}
			if (event.key === "Backspace") {
				handleBackspace(index, event);
				return;
			}
			if (event.key === "Delete") {
				handleDelete(index, event);
			}
		},
		[handleBackspace, handleDelete, handleEnter],
	);

	const handleScroll = useCallback(
		(event: UIEvent<HTMLDivElement>) => {
			if (!lineNumberRef.current) {
				return;
			}
			lineNumberRef.current.scrollTop = event.currentTarget.scrollTop;
		},
		[],
	);

	const resolvedRemoteUsers = useMemo(
		() =>
			remoteUsers.map((user, idx) => ({
				...user,
				color:
					user.color ||
					[
						"var(--color-cursor-other-1)",
						"var(--color-cursor-other-2)",
						"var(--color-cursor-other-3)",
						"var(--color-cursor-other-4)",
					][idx % 4],
			})),
		[remoteUsers],
	);

	const remoteCursorLocations = useMemo(() => {
		return resolvedRemoteUsers.reduce<Record<string, CursorLocation>>((map, user) => {
			const location = toCursorLocation(lines, user);
			if (location) {
				map[user.id] = location;
			}
			return map;
		}, {});
	}, [lines, resolvedRemoteUsers]);

	const renderRemoteCursorIndicators = () => {
		const scrollTop = contentRef.current?.scrollTop ?? 0;
		const scrollLeft = contentRef.current?.scrollLeft ?? 0;
		const lineHeightPx = 1.6 * 14; // matches CSS line-height and font size
		const paddingY = 12;
		const paddingX = 16;
		const characterWidth = 8;

		return (
			<div className="editor-content-overlay" aria-hidden="true">
				{resolvedRemoteUsers.map((user) => {
					const location = remoteCursorLocations[user.id];
					if (!location) {
						return null;
					}
					const top = paddingY + location.line * lineHeightPx - scrollTop;
					const left = paddingX + location.column * characterWidth - scrollLeft;
					return (
						<div
							key={user.id}
							className="editor-cursor-indicator"
							style={{
								top,
								left,
								background: user.color,
							}}
							aria-label={`${user.name} cursor`}
						>
							<span className="editor-cursor-label" style={{ borderColor: user.color }}>
								{user.name}
							</span>
						</div>
					);
				})}
			</div>
		);
	};

	const editorClass = ["editor-frame", className].filter(Boolean).join(" ");

	return (
		<section className={editorClass} data-color-scheme={scheme}>
			<div className="sr-only" aria-live="polite">
				Line {activeLine + 1} of {lines.length}
			</div>
			<header className="editor-toolbar">
				<div className="editor-status-badge">
					<span className="editor-status-dot" />
					{isConnected ? "Connected" : "Connecting"}
				</div>
				<span className="editor-toolbar-subtle">
					Monospace · Inline selection · Live line numbers
				</span>
			</header>
			<div className="editor-container">
				<div
					className="editor-line-numbers"
					aria-hidden="true"
					ref={lineNumberRef}
					style={{ lineHeight: "var(--editor-line-height)" }}
					>
						{lines.map((_, index) => (
							<div
								key={lineKey(index)}
								className={`editor-line-number ${index === activeLine ? "active" : ""}`}
							>
								{index + 1}
							</div>
						))}
				</div>

				<div className="editor-content-column">
					<div
						ref={contentRef}
						className="editor-content-area"
						role="textbox"
						aria-multiline="true"
						aria-label={ariaLabel}
						onScroll={handleScroll}
					>
						{lines.map((line, index) => (
							<div
								key={lineKey(index)}
								className={`editor-line ${index === activeLine ? "active" : ""}`}
								data-line-index={index}
								contentEditable
								suppressContentEditableWarning
								onInput={(event) => handleInput(index, event)}
								onKeyDown={(event) => handleKeyDown(index, event)}
								onFocus={() => {
									setActiveLine(index);
									requestAnimationFrame(syncSelectionState);
								}}
								onClick={() => {
									setActiveLine(index);
									requestAnimationFrame(syncSelectionState);
								}}
								onPaste={(event) => handlePaste(index, event)}
								spellCheck={false}
							>
								{line || ZERO_WIDTH_SPACE}
							</div>
						))}
					</div>
					{renderRemoteCursorIndicators()}
				</div>

				{showPresence && presenceVisible && (
					<aside className="remote-preview-panel">
						<div className="remote-panel-header">
							<div className="remote-panel-labels">
								<span className={`remote-badge ${isConnected ? "connected" : ""}`}>
									{isConnected ? "Connected" : "Connecting"}
								</span>
								<span className={`remote-badge just-you ${remoteUsers.length === 0 ? "active" : ""}`}>
									{remoteUsers.length === 0 ? "Just you" : `${remoteUsers.length} collaborators`}
								</span>
							</div>
							<button
								type="button"
								className="remote-panel-toggle"
								onClick={() => setPresenceVisible(false)}
								aria-label="Hide collaborators panel"
							>
								<svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
									<path
										d="M18 6L6 18M6 6l12 12"
										stroke="currentColor"
										strokeWidth="2"
										strokeLinecap="round"
									/>
								</svg>
							</button>
						</div>
						<ul className="remote-users">
							{remoteUsers.length === 0 && (
								<li className="remote-user" aria-label="No remote users">
									<span className="user-status-dot offline" />
									<div className="user-meta">
										<span className="user-name">No remote collaborators</span>
										<span className="editor-toolbar-subtle">
											Invite teammates to see their cursors here.
										</span>
									</div>
								</li>
							)}
							{resolvedRemoteUsers.map((user) => {
								const cursor = remoteCursorLocations[user.id];
								const status = user.status ?? "online";
								return (
									<li key={user.id} className="remote-user">
										<span className={`user-status-dot ${status}`} />
										<div className="user-meta">
											<span className="user-name" title={user.name}>
												{user.name}
											</span>
											<span className="user-role">
												<span
													className="user-color"
													style={{ background: user.color }}
													aria-hidden="true"
												/>
												{user.role ?? "Editor"}
											</span>
										</div>
										{cursor && (
											<span className="user-cursor-chip">
												<span
													className="user-color"
													style={{ background: user.color }}
													aria-hidden="true"
												/>
												Line {cursor.line + 1}, Col {cursor.column + 1}
											</span>
										)}
									</li>
								);
							})}
						</ul>
					</aside>
				)}

				{showPresence && !presenceVisible && (
					<div className="remote-panel-collapsed">
						<button
							type="button"
							className="remote-panel-toggle"
							onClick={() => setPresenceVisible(true)}
							aria-label="Show collaborators panel"
						>
							<svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
								<path
									d="M4 7a4 4 0 118 0 4 4 0 01-8 0zm8.5 6c-3.5 0-6.5 1.7-6.5 3.8V19h13v-2.2C19 14.7 16 13 12.5 13zm5-2a2.5 2.5 0 100-5 2.5 2.5 0 000 5zm-1 1a4.3 4.3 0 014 2.3 4.2 4.2 0 01.5 2.1V18h-3"
									stroke="currentColor"
									strokeWidth="1.6"
									fill="none"
									strokeLinecap="round"
									strokeLinejoin="round"
								/>
							</svg>
						</button>
					</div>
				)}
			</div>
		</section>
	);
};
