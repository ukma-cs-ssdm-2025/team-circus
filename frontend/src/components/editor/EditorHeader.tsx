import { useMemo, useState } from "react";
import styles from "./EditorHeader.module.css";
import { useLanguage } from "../../contexts/LanguageContext";
import DownloadIcon from "../icons/DownloadIcon";

type ExportFormat = "md" | "html" | "pdf";

type EditorHeaderProps = {
	docName: string;
	onNameChange: (name: string) => void;
	onExport: (format: ExportFormat) => void;
	wordCount: number;
	readingTime: number;
	isConnected: boolean;
};

export const EditorHeader = ({
	docName,
	onNameChange,
	onExport,
	wordCount,
	readingTime,
	isConnected,
}: EditorHeaderProps) => {
	const { t } = useLanguage();
	const [isExportOpen, setIsExportOpen] = useState(false);

	const formattedWordCount = useMemo(() => {
		return new Intl.NumberFormat().format(wordCount);
	}, [wordCount]);

	const statusLabel = useMemo(
		() =>
			isConnected
				? t("documentEditor.liveStatusConnected")
				: t("documentEditor.liveStatusConnecting"),
		[isConnected, t],
	);

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
						isConnected ? styles.success : styles.neutral
					}`}
					role="status"
					aria-live="polite"
				>
					{statusLabel}
				</div>
			</div>
		</header>
	);
};
