import { useMemo, useState } from "react";
import styles from "./EditorHeader.module.css";
import { useLanguage } from "../../contexts/LanguageContext";
import DownloadIcon from "../icons/DownloadIcon";

type ExportFormat = "md" | "html" | "pdf";

type EditorHeaderProps = {
	docName: string;
	onNameChange: (name: string) => void;
	onSave: () => void | Promise<void>;
	onExport: (format: ExportFormat) => void;
	isSaving: boolean;
	isDirty: boolean;
	wordCount: number;
	readingTime: number;
	saveStatus: "idle" | "saving" | "success" | "error";
};

export const EditorHeader = ({
	docName,
	onNameChange,
	onSave,
	onExport,
	isSaving,
	isDirty,
	wordCount,
	readingTime,
	saveStatus,
}: EditorHeaderProps) => {
	const { t } = useLanguage();
	const [isExportOpen, setIsExportOpen] = useState(false);

	const formattedWordCount = useMemo(() => {
		return new Intl.NumberFormat().format(wordCount);
	}, [wordCount]);

	const statusLabel = useMemo(() => {
		if (saveStatus === "saving") {
			return t("documentEditor.savingButton");
		}
		if (saveStatus === "success") {
			return t("documentEditor.saveStatusSaved");
		}
		if (saveStatus === "error") {
			return t("documentEditor.saveStatusError");
		}
		if (isDirty) {
			return t("documentEditor.unsavedChanges");
		}
		return t("documentEditor.saveStatusIdle");
	}, [isDirty, saveStatus, t]);

	const isSaveDisabled = isSaving || !isDirty || docName.trim().length === 0;

	const exportOptions: { value: ExportFormat; label: string }[] = useMemo(
		() => [
			{ value: "md", label: t("documentEditor.exportMd") },
			{ value: "html", label: t("documentEditor.exportHtml") },
			{ value: "pdf", label: t("documentEditor.exportPdf") },
		],
		[t],
	);

	const handleExportSelect = (format: ExportFormat) => {
		onExport(format);
		setIsExportOpen(false);
	};

	return (
		<header className={styles.header}>
			<div className={styles.left}>
				<input
					id="document-name"
					type="text"
					value={docName}
					onChange={(event) => onNameChange(event.target.value)}
					placeholder={t("documentEditor.namePlaceholder")}
					className={styles.nameInput}
					aria-label={t("documentEditor.nameLabel")}
					aria-invalid={docName.trim().length === 0}
				/>

				<div
					className={styles.exportWrapper}
					onBlur={(event) => {
						if (!event.currentTarget.contains(event.relatedTarget)) {
							setIsExportOpen(false);
						}
					}}
				>
					<button
						type="button"
						className={styles.exportButton}
						onClick={() => setIsExportOpen((open) => !open)}
						aria-haspopup="menu"
						aria-expanded={isExportOpen}
						title={t("documentEditor.exportLabel")}
					>
						<DownloadIcon className={styles.exportIcon} />
					</button>
					{isExportOpen && (
						<div className={styles.exportMenu} role="menu">
							{exportOptions.map((option) => (
								<button
									key={option.value}
									type="button"
									className={styles.exportItem}
									onClick={() => handleExportSelect(option.value)}
									role="menuitem"
								>
									{option.label}
								</button>
							))}
						</div>
					)}
				</div>

				<button
					type="button"
					className={styles.saveButton}
					disabled={isSaveDisabled}
					onClick={() => void onSave()}
				>
					{isSaving
						? t("documentEditor.savingButton")
						: t("documentEditor.saveButton")}
				</button>
			</div>

			<div className={styles.right}>
				<div className={styles.stats}>
					<span>
						{t("documentEditor.wordCountLabel")}: {formattedWordCount}
					</span>
					<span className={styles.divider}>|</span>
					<span>
						{t("documentEditor.readingTimeLabel")}: {readingTime || 0}{" "}
						{t("documentEditor.readingTimeUnit")}
					</span>
				</div>
				<div
					className={`${styles.status} ${
						saveStatus === "success"
							? styles.success
							: saveStatus === "error"
								? styles.error
								: styles.neutral
					}`}
					role="status"
					aria-live="polite"
				>
					{saveStatus === "success" && (
						<span className={styles.dot} aria-hidden />
					)}
					{saveStatus === "error" && (
						<span className={styles.dot} aria-hidden />
					)}
					{statusLabel}
				</div>
			</div>
		</header>
	);
};
