import { memo, useMemo } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";
import type { HTMLAttributes, ReactNode } from "react";
import styles from "./MarkdownPreview.module.css";
import { useLanguage } from "../../contexts/LanguageContext";

type MarkdownPreviewProps = {
	content: string;
	emptyText?: string;
};

type CodeProps = HTMLAttributes<HTMLElement> & {
	inline?: boolean;
	children?: ReactNode;
	className?: string;
};

const componentsMap: Components = {
	code: (props) => {
		const { inline, children, className } = props as CodeProps;
		const inferredInline =
			typeof inline === "boolean"
				? inline
				: !className || !className.includes("language-");

		return inferredInline ? (
			<code className={styles.inlineCode}>{children}</code>
		) : (
			<pre className={styles.codeBlock}>
				<code className={className}>{children}</code>
			</pre>
		);
	},
	a: ({ href, children }) => (
		<a
			className={styles.link}
			href={href}
			target="_blank"
			rel="noreferrer noopener"
		>
			{children}
		</a>
	),
	blockquote: ({ children }) => (
		<blockquote className={styles.blockquote}>{children}</blockquote>
	),
	table: ({ children }) => (
		<div className={styles.tableWrapper}>
			<table>{children}</table>
		</div>
	),
	th: ({ children }) => <th className={styles.tableHeader}>{children}</th>,
	td: ({ children }) => <td className={styles.tableCell}>{children}</td>,
};

export const MarkdownPreview = memo(
	({ content, emptyText }: MarkdownPreviewProps) => {
		const { t } = useLanguage();
		const trimmedContent = content.trim();

		const components = useMemo<Components>(() => componentsMap, []);

		const rendered = useMemo(() => {
			if (!trimmedContent) {
				return null;
			}

			return (
				<ReactMarkdown
					remarkPlugins={[remarkGfm]}
					className={styles.markdown}
					components={components}
				>
					{trimmedContent}
				</ReactMarkdown>
			);
		}, [components, trimmedContent]);

		if (!trimmedContent) {
			return (
				<div className={styles.empty}>
					{emptyText ?? t("documentEditor.previewEmpty")}
				</div>
			);
		}

		return <div className={styles.preview}>{rendered}</div>;
	},
);

MarkdownPreview.displayName = "MarkdownPreview";
