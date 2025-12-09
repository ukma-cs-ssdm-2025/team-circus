import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { DocumentItem } from "../types";

type SaveStatus = "idle" | "saving" | "success" | "error";

type SavedState = {
	name: string;
	content: string;
};

interface EditorStateReturn {
	state: {
		name: string;
		content: string;
		lastSaved: SavedState | null;
		isDirty: boolean;
		saveStatus: SaveStatus;
		wordCount: number;
		readingTime: number;
	};
	handlers: {
		setName: (name: string) => void;
		setContent: (content: string) => void;
		markSaved: (savedState: SavedState) => void;
		setSaveStatus: (status: SaveStatus) => void;
		reset: () => void;
	};
	computed: {
		isNameValid: boolean;
		isSaveDisabled: boolean;
	};
}

export function useEditorState(
	initialDocument?: Partial<DocumentItem> | null,
): EditorStateReturn {
	const initialName = initialDocument?.name ?? "";
	const initialContent = initialDocument?.content ?? "";

	const initialStateRef = useRef<SavedState>({
		name: initialName,
		content: initialContent,
	});

	const [name, setName] = useState(initialName);
	const [content, setContent] = useState(initialContent);
	const [lastSaved, setLastSaved] = useState<SavedState | null>(
		initialDocument ? { name: initialName, content: initialContent } : null,
	);
	const [saveStatus, setSaveStatus] = useState<SaveStatus>("idle");

	useEffect(() => {
		const nextName = initialDocument?.name ?? "";
		const nextContent = initialDocument?.content ?? "";

		setName(nextName);
		setContent(nextContent);
		setLastSaved({ name: nextName, content: nextContent });
		initialStateRef.current = { name: nextName, content: nextContent };
		setSaveStatus("idle");
	}, [initialDocument?.name, initialDocument?.content]);

	useEffect(() => {
		if (saveStatus !== "success" && saveStatus !== "error") {
			return;
		}

		const timer = window.setTimeout(() => {
			setSaveStatus("idle");
		}, 3000);

		return () => {
			window.clearTimeout(timer);
		};
	}, [saveStatus]);

	const wordCount = useMemo(() => {
		const trimmed = content.trim();
		if (!trimmed) {
			return 0;
		}

		return trimmed.split(/\s+/).filter(Boolean).length;
	}, [content]);

	const readingTime = useMemo(() => {
		if (wordCount === 0) {
			return 0;
		}

		return Math.ceil(wordCount / 200);
	}, [wordCount]);

	const isNameValid = name.trim().length > 0;
	const isDirty =
		lastSaved !== null &&
		(lastSaved.name !== name || lastSaved.content !== content);
	const isSaveDisabled = !isNameValid || !isDirty || saveStatus === "saving";

	const markSaved = useCallback((savedState: SavedState) => {
		setLastSaved(savedState);
	}, []);

	const reset = useCallback(() => {
		const fallback = lastSaved ?? initialStateRef.current;
		setName(fallback?.name ?? "");
		setContent(fallback?.content ?? "");
		setSaveStatus("idle");
	}, [lastSaved]);

	return {
		state: {
			name,
			content,
			lastSaved,
			isDirty,
			saveStatus,
			wordCount,
			readingTime,
		},
		handlers: {
			setName,
			setContent,
			markSaved,
			setSaveStatus,
			reset,
		},
		computed: {
			isNameValid,
			isSaveDisabled,
		},
	};
}
