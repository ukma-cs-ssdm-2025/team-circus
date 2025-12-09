import { renderToStaticMarkup } from "react-dom/server";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useCallback } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "./DocumentEditor.module.css";
import {
	EditorHeader,
	EditorLayout,
	ErrorAlert,
	LoadingSpinner,
	MarkdownEditor,
	MarkdownPreview,
} from "../components";
import { ROUTES } from "../constants";
import { useLanguage } from "../contexts/LanguageContext";
import { useDebounce, useDocumentSync, useEditorState } from "../hooks";
import type { BaseComponentProps } from "../types";

type DocumentEditorProps = BaseComponentProps;

const DocumentEditor = ({ className = "" }: DocumentEditorProps) => {
	const { t } = useLanguage();
	const navigate = useNavigate();
	const { uuid } = useParams<{ uuid: string }>();
	const documentId = uuid ?? "";

	const {
		document: documentData,
		loading,
		error,
		saveDocument,
		refetch,
	} = useDocumentSync(documentId);

	const { state, handlers, computed } = useEditorState(documentData);
	const debouncedContent = useDebounce(state.content, 300);

	const buildFileName = useCallback(
		(extension: string) => {
			const safeTitle =
				state.name.trim() || t("documentEditor.fallbackTitle");
			const sanitized = safeTitle.replace(/[\\/:*?"<>|]+/g, "").trim() || "document";
			return `${sanitized}.${extension}`;
		},
		[state.name, t],
	);

	const handleExport = useCallback(
		(format: "md" | "html" | "pdf") => {
			const fileName = buildFileName(format);
			const markdown = state.content || "";

			const downloadFile = (content: string, type: string) => {
				const blob = new Blob([content], { type });
				const url = URL.createObjectURL(blob);
				const link = window.document.createElement("a");
				link.href = url;
				link.download = fileName;
				link.click();
				URL.revokeObjectURL(url);
			};

			if (format === "md") {
				downloadFile(markdown, "text/markdown");
				return;
			}

			const htmlBody = renderToStaticMarkup(
				<article>
					<h1>{state.name || t("documentEditor.fallbackTitle")}</h1>
					<ReactMarkdown remarkPlugins={[remarkGfm]}>
						{markdown}
					</ReactMarkdown>
				</article>,
			);
			const htmlDocument = `<!doctype html><html><head><meta charset="utf-8" /><title>${state.name}</title></head><body>${htmlBody}</body></html>`;

			if (format === "html") {
				downloadFile(htmlDocument, "text/html");
				return;
			}

			const printable = window.open("", "_blank");
			if (!printable) {
				return;
			}
			printable.document.write(htmlDocument);
			printable.document.close();
			printable.focus();
			printable.print();
		},
		[buildFileName, state.content, state.name, t],
	);

	const handleSave = useCallback(async () => {
		if (computed.isSaveDisabled || !documentId) {
			return;
		}

		handlers.setSaveStatus("saving");
		try {
			const saved = await saveDocument({
				name: state.name.trim(),
				content: state.content,
			});

			handlers.markSaved({
				name: saved.name ?? state.name.trim(),
				content: saved.content ?? state.content,
			});
			handlers.setSaveStatus("success");
		} catch (saveError) {
			console.error("Failed to save document", saveError);
			handlers.setSaveStatus("error");
		}
	}, [
		computed.isSaveDisabled,
		documentId,
		handlers,
		saveDocument,
		state.content,
		state.name,
	]);

	if (!documentId) {
		return (
			<div className={`${styles.page} ${className ?? ""}`}>
				<div className={styles.backRow}>
					<button
						type="button"
						className={styles.backButton}
						onClick={() => navigate(ROUTES.DOCUMENTS)}
					>
						{t("documentEditor.backToList")}
					</button>
				</div>
				<div className={styles.errorWrapper}>
					<ErrorAlert message={t("documentEditor.notFound")} />
				</div>
			</div>
		);
	}

	return (
		<div className={`${styles.page} ${className ?? ""}`}>
			{error && (
				<div className={styles.errorWrapper}>
					<ErrorAlert
						message={t("documentEditor.loadError")}
						onRetry={refetch}
						retryText={t("documents.refresh")}
					/>
				</div>
			)}

			{loading && <LoadingSpinner py={4} />}

			{!loading && !error && !documentData && (
				<div className={styles.emptyState}>
					{t("documentEditor.notFound")}
				</div>
			)}

			{!loading && documentData && (
				<div className={styles.content}>
					<EditorHeader
						docName={state.name}
						onNameChange={handlers.setName}
						onSave={handleSave}
						onExport={handleExport}
						isSaving={state.saveStatus === "saving"}
						isDirty={state.isDirty}
						wordCount={state.wordCount}
						readingTime={state.readingTime}
						saveStatus={state.saveStatus}
					/>

					<div className={styles.editorShell}>
						<EditorLayout resizable={false} className={styles.editorLayout}>
							<MarkdownEditor
								value={state.content}
								onChange={handlers.setContent}
								placeholder={t("documentEditor.contentPlaceholder")}
								onSave={handleSave}
							/>
							<MarkdownPreview
								content={debouncedContent}
								emptyText={t("documentEditor.previewEmpty")}
							/>
						</EditorLayout>
					</div>
				</div>
			)}
		</div>
	);
};

export default DocumentEditor;
