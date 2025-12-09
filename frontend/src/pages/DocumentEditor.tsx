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
	ShareDialog,
} from "../components";
import { ROUTES } from "../constants";
import { useAuth } from "../contexts/AuthContextBase";
import { useLanguage } from "../contexts/LanguageContext";
import { useCollaborativeEditor, useDebounce, useDocumentSync } from "../hooks";
import { generateShareLink } from "../services";
import { getGroupMembers } from "../services/members";
import type { BaseComponentProps, ShareLinkResponse } from "../types";

type DocumentEditorProps = BaseComponentProps;
const DEFAULT_SHARE_EXPIRATION_DAYS = 7;

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
	const [saveError, setSaveError] = useState<string | null>(null);
	const [userRole, setUserRole] = useState<string | null>(null);

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
			role: userRole ?? (user ? "editor" : "viewer"),
		}),
		[user?.email, user?.login, user?.uuid, getAnonymousId, userRole, user],
	);

	const {
		document: documentData,
		loading,
		error,
		refetch,
		saveDocument,
	} = useDocumentSync(documentId);

	const [shareDialogOpen, setShareDialogOpen] = useState(false);
	const [shareLink, setShareLink] = useState<ShareLinkResponse | null>(null);
	const [shareError, setShareError] = useState<string | null>(null);
	const [shareLoading, setShareLoading] = useState(false);
	const [docName, setDocName] = useState(documentData?.name ?? "");
	const lastSavedRef = useRef<{ name: string; content: string } | null>(null);
	const initialContentHydratedRef = useRef(false);
	const saveRequestIdRef = useRef(0);

	const {
		content,
		setContent,
		isConnected,
		remoteUsers,
		updateCursorPosition,
		yDoc,
	} = useCollaborativeEditor({
		documentId,
		user: collaborativeUser,
	});

	useEffect(() => {
		if (documentData?.name) {
			setDocName(documentData.name);
		}
	}, [documentData?.name]);

	useEffect(() => {
		if (!documentData) {
			return;
		}
		lastSavedRef.current = {
			name: documentData.name,
			content: documentData.content ?? "",
		};
	}, [documentData?.content, documentData?.name]);

	useEffect(() => {
		let cancelled = false;
		const fetchRole = async () => {
			if (!documentData?.group_uuid || !user?.uuid) {
				return;
			}
			try {
				const members = await getGroupMembers(documentData.group_uuid);
				if (cancelled) {
					return;
				}
				const current = members.members?.find(
					(member) => member.user_uuid === user.uuid,
				);
				setUserRole(current?.role ?? null);
			} catch {
				if (!cancelled) {
					setUserRole(null);
				}
			}
		};

		fetchRole();

		return () => {
			cancelled = true;
		};
	}, [documentData?.group_uuid, user?.uuid]);

	useEffect(() => {
		if (initialContentHydratedRef.current) {
			return;
		}
		if (!documentData || !yDoc) {
			return;
		}

		const existingContent = documentData.content ?? "";
		const hasCollaborativeContent = yDoc.getText("content").length > 0;

		if (!hasCollaborativeContent && existingContent) {
			setContent(existingContent);
		}

		initialContentHydratedRef.current = true;
	}, [documentData, setContent, yDoc]);

	const debouncedContent = useDebounce(content, 300);
	const debouncedSaveState = useDebounce(
		{
			name: docName,
			content,
		},
		800,
	);

	useEffect(() => {
		if (!documentData) {
			return;
		}

		const payloadName = debouncedSaveState.name.trim();
		const payloadContent = debouncedSaveState.content ?? "";

		if (!payloadName) {
			return;
		}

		const lastSaved = lastSavedRef.current ?? {
			name: documentData.name,
			content: documentData.content ?? "",
		};

		if (
			lastSaved.name === payloadName &&
			lastSaved.content === payloadContent
		) {
			return;
		}

		const requestId = saveRequestIdRef.current + 1;
		saveRequestIdRef.current = requestId;
		setSaveError(null);

		saveDocument({
			name: payloadName,
			content: payloadContent,
		})
			.then((saved) => {
				if (requestId !== saveRequestIdRef.current) {
					return;
				}
				lastSavedRef.current = {
					name: saved.name,
					content: saved.content ?? payloadContent,
				};
				setSaveError(null);
			})
			.catch((saveErr) => {
				if (requestId !== saveRequestIdRef.current) {
					return;
				}
				const message =
					saveErr instanceof Error
						? saveErr.message
						: t("documentEditor.saveError");
				setSaveError(message || t("documentEditor.saveError"));
			});
	}, [debouncedSaveState, documentData, saveDocument, t]);

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

			const downloadFile = (value: string, type: string) => {
				const blob = new Blob([value], { type });
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

	const handleShareOpen = useCallback(() => {
		setShareDialogOpen(true);
		setShareLink(null);
		setShareError(null);
	}, []);

	const handleShareClose = useCallback(() => {
		if (shareLoading) {
			return;
		}
		setShareDialogOpen(false);
		setShareError(null);
	}, [shareLoading]);

	const handleGenerateShare = useCallback(
		async (expirationDays: number) => {
			if (!documentId) {
				return;
			}
			setShareLoading(true);
			setShareError(null);
			setShareLink(null);

			try {
				const result = await generateShareLink(documentId, expirationDays);
				setShareLink(result);
			} catch (shareErr) {
				const message =
					shareErr instanceof Error ? shareErr.message : t("shareDialog.error");
				setShareError(message || t("shareDialog.error"));
			} finally {
				setShareLoading(false);
			}
		},
		[documentId, t],
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

			{saveError && !loading && (
				<div className={styles.errorWrapper}>
					<ErrorAlert message={saveError} />
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
						onShare={documentData ? handleShareOpen : undefined}
						shareDisabled={shareLoading || !documentData}
					/>

					<div className={styles.editorShell}>
						<EditorLayout resizable={false} className={styles.editorLayout}>
							<Editor
								value={content}
								onChange={setContent}
								onCursorChange={updateCursorPosition}
								isConnected={isConnected}
								remoteUsers={remoteUsers.map((remoteUser) => ({
									...remoteUser,
									name: remoteUser.name ?? remoteUser.id,
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

			<ShareDialog
				open={shareDialogOpen}
				onClose={handleShareClose}
				onGenerate={handleGenerateShare}
				loading={shareLoading}
				link={shareLink}
				error={shareError}
				defaultExpirationDays={DEFAULT_SHARE_EXPIRATION_DAYS}
			/>
		</div>
	);
};

export default DocumentEditor;
