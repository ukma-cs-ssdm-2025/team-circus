import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { Compartment, EditorState } from "@codemirror/state";
import {
	EditorView,
	drawSelection,
	dropCursor,
	highlightActiveLine,
	highlightActiveLineGutter,
	keymap,
	lineNumbers,
	rectangularSelection,
} from "@codemirror/view";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTheme } from "../../contexts/ThemeContext";
import type { EditorProps } from "./types";
import "./editor.css";

type Palette = {
	background: string;
	foreground: string;
	gutterBackground: string;
	gutterText: string;
	gutterActive: string;
	border: string;
	selection: string;
	cursor: string;
	activeLine: string;
};

const createPalette = (mode: "light" | "dark"): Palette =>
	mode === "light"
		? {
				background: "#ffffff",
				foreground: "#18181b",
				gutterBackground: "#f9fafb",
				gutterText: "#a1a1a6",
				gutterActive: "#3b82f6",
				border: "#e4e4e7",
				selection: "rgba(59, 130, 246, 0.28)",
				cursor: "#3b82f6",
				activeLine: "rgba(59, 130, 246, 0.12)",
			}
		: {
				background: "#0f172a",
				foreground: "#e2e8f0",
				gutterBackground: "#0a0a0a",
				gutterText: "#64748b",
				gutterActive: "#60a5fa",
				border: "#334155",
				selection: "rgba(96, 165, 250, 0.4)",
				cursor: "#60a5fa",
				activeLine: "rgba(96, 165, 250, 0.12)",
			};

const buildTheme = (mode: "light" | "dark") => {
	const palette = createPalette(mode);
	return EditorView.theme(
		{
			"&": {
				backgroundColor: palette.background,
				color: palette.foreground,
			},
			"&.cm-editor": {
				border: `1px solid ${palette.border}`,
				borderRadius: "12px",
			},
			"&.cm-editor.cm-focused": {
				outline: `1px solid ${palette.gutterActive}`,
				outlineOffset: "2px",
			},
			".cm-scroller": {
				fontFamily:
					'"Monaco", "Courier New", "SFMono-Regular", ui-monospace, monospace',
				lineHeight: "1.6",
				scrollbarColor: `${palette.border} transparent`,
				tabSize: "4",
			},
			".cm-content": {
				caretColor: palette.cursor,
				minHeight: "320px",
			},
			".cm-line": {
				position: "relative",
			},
			".cm-selectionBackground, .cm-content ::selection": {
				backgroundColor: palette.selection,
			},
			".cm-activeLine": {
				backgroundColor: palette.activeLine,
			},
			".cm-cursor": {
				borderLeftColor: palette.cursor,
			},
			".cm-gutters": {
				backgroundColor: palette.gutterBackground,
				color: palette.gutterText,
				borderRight: `1px solid ${palette.border}`,
			},
			".cm-activeLineGutter": {
				color: palette.gutterActive,
				fontWeight: "700",
			},
		},
		{ dark: mode === "dark" },
	);
};

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
	const scheme = colorScheme ?? theme ?? "dark";

	const [presenceVisible, setPresenceVisible] = useState(showPresence);
	const hostRef = useRef<HTMLDivElement | null>(null);
	const viewRef = useRef<EditorView | null>(null);
	const themeCompartmentRef = useRef(new Compartment());
	const updatingFromProps = useRef(false);
	const lastCursorPosition = useRef<number | null>(null);
	const initialDocRef = useRef(value ?? defaultValue ?? "");

	useEffect(() => {
		setPresenceVisible(showPresence);
	}, [showPresence]);

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

	useEffect(() => {
		if (!hostRef.current) {
			return undefined;
		}

		const themeCompartment = themeCompartmentRef.current;
		const state = EditorState.create({
			doc: initialDocRef.current,
			extensions: [
				lineNumbers(),
				highlightActiveLine(),
				highlightActiveLineGutter(),
				drawSelection(),
				dropCursor(),
				rectangularSelection(),
				history(),
				keymap.of([...defaultKeymap, ...historyKeymap]),
				themeCompartment.of(buildTheme(scheme)),
				EditorView.contentAttributes.of({
					"aria-label": ariaLabel,
					spellcheck: "false",
					"data-testid": "team-circus-editor",
				}),
				EditorView.updateListener.of((update) => {
					if (update.docChanged && onChange && !updatingFromProps.current) {
						onChange(update.state.doc.toString());
					}

					if (!onCursorChange) {
						return;
					}
					if (update.selectionSet || update.docChanged) {
						const position = update.state.selection.main.head;
						if (lastCursorPosition.current !== position) {
							lastCursorPosition.current = position;
							onCursorChange(position);
						}
					}
				}),
			],
		});

		viewRef.current = new EditorView({
			state,
			parent: hostRef.current,
		});

		return () => {
			viewRef.current?.destroy();
			viewRef.current = null;
		};
	}, [ariaLabel, onChange, onCursorChange, scheme]);

	useEffect(() => {
		if (!viewRef.current) {
			return;
		}
		const nextDoc = value ?? defaultValue ?? "";
		const currentDoc = viewRef.current.state.doc.toString();
		if (nextDoc === currentDoc) {
			return;
		}
		const currentSelection = viewRef.current.state.selection;
		const scrollElement = viewRef.current.scrollDOM;
		const previousScrollTop = scrollElement.scrollTop;
		const previousScrollLeft = scrollElement.scrollLeft;
		updatingFromProps.current = true;
		let start = 0;
		while (
			start < currentDoc.length &&
			start < nextDoc.length &&
			currentDoc[start] === nextDoc[start]
		) {
			start += 1;
		}

		let endCurrent = currentDoc.length;
		let endNext = nextDoc.length;
		while (
			endCurrent > start &&
			endNext > start &&
			currentDoc[endCurrent - 1] === nextDoc[endNext - 1]
		) {
			endCurrent -= 1;
			endNext -= 1;
		}

		if (start === endCurrent && start === endNext) {
			updatingFromProps.current = false;
			return;
		}

		const changeSpec = {
			from: start,
			to: endCurrent,
			insert: nextDoc.slice(start, endNext),
		};

		const changeDesc = viewRef.current.state.changes(changeSpec);
		const mappedSelection = currentSelection.map(changeDesc, 1);

		viewRef.current.dispatch({
			changes: changeSpec,
			selection: mappedSelection,
		});
		requestAnimationFrame(() => {
			scrollElement.scrollTop = previousScrollTop;
			scrollElement.scrollLeft = previousScrollLeft;
		});
		updatingFromProps.current = false;
	}, [defaultValue, value]);

	useEffect(() => {
		if (!viewRef.current) {
			return;
		}
		viewRef.current.dispatch({
			effects: themeCompartmentRef.current.reconfigure(buildTheme(scheme)),
		});
	}, [scheme]);

	const editorClass = ["editor-frame", className].filter(Boolean).join(" ");

	return (
		<section className={editorClass} data-color-scheme={scheme}>
			<header className="editor-toolbar">
				<div className="editor-status-badge">
					<span className="editor-status-dot" />
					{isConnected ? "Connected" : "Connecting"}
				</div>
				<button
					type="button"
					className="remote-panel-toggle"
					onClick={() => setPresenceVisible((prev) => !prev)}
					aria-label={
						presenceVisible
							? "Hide collaborators panel"
							: "Show collaborators panel"
					}
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
			</header>
			<div className="editor-body">
				{presenceVisible && (
					<div className="remote-preview-panel">
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
									</li>
								);
							})}
						</ul>
					</div>
				)}

				<div className="editor-pane">
					<div className="editor-host" ref={hostRef} role="presentation" />
				</div>
			</div>
		</section>
	);
};
