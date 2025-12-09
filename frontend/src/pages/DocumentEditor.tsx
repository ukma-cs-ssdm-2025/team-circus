import { renderToStaticMarkup } from "react-dom/server";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "./DocumentEditor.module.css";
import {
	Editor,
	EditorHeader,
	EditorLayout,
	ErrorAlert,
	LoadingSpinner,
	MarkdownPreview,
} from "../components";
import { ROUTES } from "../constants";
import { useAuth } from "../contexts/AuthContextBase";
import { useLanguage } from "../contexts/LanguageContext";
import { useCollaborativeEditor, useDebounce, useDocumentSync } from "../hooks";
import type { BaseComponentProps } from "../types";

type DocumentEditorProps = BaseComponentProps;

const escapeHtml = (value: string): string =>
	value
		.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#39;");

const DocumentEditor = ({ className = "" }: DocumentEditorProps) => {
	const { t } = useLanguage();
	const navigate = useNavigate();
	const { uuid } = useParams<{ uuid: string }>();
	const documentId = uuid ?? "";
	const { user } = useAuth();

	// Generate a stable anonymous session ID
	const anonymousIdRef = useRef<string | null>(null);
	const getAnonymousId = useCallback(() => {
		if (anonymousIdRef.current) {
			return anonymousIdRef.current;
		}

		// Try to get from sessionStorage first
		const stored = sessionStorage.getItem("anonymous-session-id");
		if (stored) {
			anonymousIdRef.current = stored;
			return stored;
		}

		// Generate new unique ID
		const newId = crypto.randomUUID();
		sessionStorage.setItem("anonymous-session-id", newId);
		anonymousIdRef.current = newId;
		return newId;
	}, []);

	const collaborativeUser = useMemo(
		() => ({
			id: user?.uuid ?? getAnonymousId(),
			name:
				user?.login ??
				user?.email ??
				`Anonymous (${getAnonymousId().slice(0, 8)})`,
		}),
		[user?.email, user?.login, user?.uuid, getAnonymousId],
	);

	const {
		document: documentData,
		loading,
		error,
		refetch,
	} = useDocumentSync(documentId);

	const [docName, setDocName] = useState(documentData?.name ?? "");

	const {
		content,
		setContent,
		isConnected,
		remoteUsers,
		updateCursorPosition,
	} = useCollaborativeEditor({
		documentId,
		user: collaborativeUser,
	});

	useEffect(() => {
		if (documentData?.name) {
			setDocName(documentData.name);
		}
	}, [documentData?.name]);

	const debouncedContent = useDebounce(content, 300);
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

	const buildFileName = useCallback(
		(extension: string) => {
			const safeTitle = docName.trim() || t("documentEditor.fallbackTitle");
			const sanitized =
				safeTitle.replace(/[\\/:*?"<>|]+/g, "").trim() || "document";
			return `${sanitized}.${extension}`;
		},
		[docName, t],
	);

	const handleExport = useCallback(
		(format: "md" | "html" | "pdf") => {
			const fileName = buildFileName(format);
			const markdown = content || "";
			const rawTitle = docName || t("documentEditor.fallbackTitle");
			const escapedTitle = escapeHtml(rawTitle);

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
					<h1>{docName || t("documentEditor.fallbackTitle")}</h1>
					<ReactMarkdown remarkPlugins={[remarkGfm]}>{markdown}</ReactMarkdown>
				</article>,
			);
			const htmlDocument = `<!doctype html><html><head><meta charset="utf-8" /><title>${escapedTitle}</title></head><body>${htmlBody}</body></html>`;

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
		[buildFileName, content, docName, t],
	);

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
				<div className={styles.emptyState}>{t("documentEditor.notFound")}</div>
			)}

			{!loading && documentData && (
				<div className={styles.content}>
					<EditorHeader
						docName={docName}
						onNameChange={setDocName}
						onExport={handleExport}
						isConnected={isConnected}
						wordCount={wordCount}
						readingTime={readingTime}
					/>

					<div className={styles.editorShell}>
						<EditorLayout resizable={false} className={styles.editorLayout}>
							<Editor
								value={content}
								onChange={setContent}
								onCursorChange={updateCursorPosition}
								isConnected={isConnected}
								remoteUsers={remoteUsers.map((user) => ({
									...user,
									name: user.name ?? user.id,
								}))}
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
