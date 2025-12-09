import { memo, useCallback, useRef } from "react";
import styles from "./MarkdownEditor.module.css";
import { useLanguage } from "../../contexts/LanguageContext";

type MarkdownEditorProps = {
	value: string;
	onChange: (content: string) => void;
	placeholder?: string;
	readOnly?: boolean;
	onSave?: () => void;
};

export const MarkdownEditor = memo(
	({
		value,
		onChange,
		placeholder,
		readOnly = false,
		onSave,
	}: MarkdownEditorProps) => {
		const { t } = useLanguage();
		const textareaRef = useRef<HTMLTextAreaElement | null>(null);

		const updateValue = useCallback(
			(nextValue: string, selection?: { start: number; end: number }) => {
				if (readOnly) {
					return;
				}

				onChange(nextValue);

				if (selection && textareaRef.current) {
					requestAnimationFrame(() => {
						textareaRef.current?.setSelectionRange(
							selection.start,
							selection.end,
						);
					});
				}
			},
			[onChange, readOnly],
		);

		const handleChange = useCallback(
			(event: React.ChangeEvent<HTMLTextAreaElement>) => {
				onChange(event.target.value);
			},
			[onChange],
		);

		const insertIndent = useCallback(() => {
			const textarea = textareaRef.current;
			if (!textarea) {
				return;
			}

			const { selectionStart, selectionEnd } = textarea;
			const indent = "\t";
			const prefix = value.slice(0, selectionStart);
			const selected = value.slice(selectionStart, selectionEnd);
			const suffix = value.slice(selectionEnd);

			const nextValue = `${prefix}${indent}${selected}${suffix}`;
			const nextStart = selectionStart + indent.length;
			const nextEnd = nextStart + selected.length;

			updateValue(nextValue, { start: nextStart, end: nextEnd });
		}, [updateValue, value]);

		const wrapSelection = useCallback(
			(wrapper: string) => {
				const textarea = textareaRef.current;
				if (!textarea) {
					return;
				}

				const { selectionStart, selectionEnd } = textarea;
				const selected = value.slice(selectionStart, selectionEnd);
				const prefix = value.slice(0, selectionStart);
				const suffix = value.slice(selectionEnd);

				const wrapped = `${wrapper}${selected || ""}${wrapper}`;
				const nextValue = `${prefix}${wrapped}${suffix}`;
				const caretStart = selectionStart + wrapper.length;
				const caretEnd = caretStart + (selected.length || 0);

				updateValue(nextValue, { start: caretStart, end: caretEnd });
			},
			[updateValue, value],
		);

		const handleKeyDown = useCallback(
			(event: React.KeyboardEvent<HTMLTextAreaElement>) => {
				const isMeta = event.metaKey || event.ctrlKey;
				const key = event.key.toLowerCase();

				if (isMeta && key === "s") {
					event.preventDefault();
					onSave?.();
					return;
				}

				if (readOnly) {
					return;
				}

				if (event.key === "Tab") {
					event.preventDefault();
					insertIndent();
					return;
				}

				if (isMeta && key === "b") {
					event.preventDefault();
					wrapSelection("**");
					return;
				}

				if (isMeta && key === "i") {
					event.preventDefault();
					wrapSelection("*");
				}
			},
			[insertIndent, onSave, readOnly, wrapSelection],
		);

		return (
			<div className={styles.editor}>
				<div className={styles.labelRow}>
					<label className={styles.label} htmlFor="markdown-editor">
						{t("documentEditor.contentLabel")}
					</label>
					<span className={styles.shortcut}>
						{t("documentEditor.shortcutSave")}
					</span>
				</div>
				<textarea
					ref={textareaRef}
					id="markdown-editor"
					className={styles.textarea}
					value={value}
					onChange={handleChange}
					onKeyDown={handleKeyDown}
					placeholder={placeholder}
					aria-label={t("documentEditor.contentLabel")}
					readOnly={readOnly}
				/>
			</div>
		);
	},
);

MarkdownEditor.displayName = "MarkdownEditor";
