import { Children, useEffect, useMemo, useRef, useState } from "react";
import type { CSSProperties, ReactNode } from "react";
import styles from "./EditorLayout.module.css";

type EditorLayoutProps = {
	children: ReactNode[];
	resizable?: boolean;
	className?: string;
};

const MIN_WIDTH_PERCENT = 25;
const MAX_WIDTH_PERCENT = 75;

export const EditorLayout = ({
	children,
	resizable = true,
	className,
}: EditorLayoutProps) => {
	const containerRef = useRef<HTMLDivElement | null>(null);
	const [leftWidth, setLeftWidth] = useState(50);
	const [isDragging, setIsDragging] = useState(false);

	const slots = useMemo(() => {
		const array = Children.toArray(children);
		if (array.length !== 2) {
			console.warn("EditorLayout expects exactly 2 children.");
		}
		return [array[0] ?? null, array[1] ?? null];
	}, [children]);

	useEffect(() => {
		if (!isDragging || !resizable) {
			return;
		}

		const handlePointerMove = (event: PointerEvent) => {
			const rect = containerRef.current?.getBoundingClientRect();
			if (!rect) {
				return;
			}

			const offsetX = event.clientX - rect.left;
			const nextWidth = (offsetX / rect.width) * 100;
			const clamped = Math.min(
				MAX_WIDTH_PERCENT,
				Math.max(MIN_WIDTH_PERCENT, nextWidth),
			);

			setLeftWidth(clamped);
		};

		const stopDragging = () => setIsDragging(false);

		window.addEventListener("pointermove", handlePointerMove);
		window.addEventListener("pointerup", stopDragging);

		return () => {
			window.removeEventListener("pointermove", handlePointerMove);
			window.removeEventListener("pointerup", stopDragging);
		};
	}, [isDragging, resizable]);

	const gridStyle = useMemo(() => {
		if (!resizable) {
			return { gridTemplateColumns: "1fr 1fr" } as CSSProperties;
		}

		return {
			"--editor-left-width": `${leftWidth}%`,
			"--editor-right-width": `${100 - leftWidth}%`,
		} as CSSProperties;
	}, [leftWidth, resizable]);

	const containerClassName = className
		? `${styles.container} ${className}`
		: styles.container;

	const fixedClass = !resizable ? styles.fixed : "";
	const mergedClassName = `${containerClassName} ${fixedClass}`.trim();

	return (
		<div ref={containerRef} className={mergedClassName} style={gridStyle}>
			<div className={styles.pane}>{slots[0]}</div>
			{resizable ? (
				<div
					className={`${styles.divider} ${
						isDragging ? styles.dividerActive : ""
					}`}
					onPointerDown={() => setIsDragging(true)}
					role="separator"
					aria-orientation="vertical"
					aria-label="Resize editor panes"
				/>
			) : null}
			<div className={styles.pane}>{slots[1]}</div>
		</div>
	);
};
