import { Button, useTheme } from "@mui/material";
import type { ElementType, ReactNode } from "react";
import type { BaseComponentProps } from "../../types";

interface ActionButtonProps extends BaseComponentProps {
	children: ReactNode;
	onClick?: () => void;
	variant?: "contained" | "outlined" | "text";
	size?: "small" | "medium" | "large";
	startIcon?: ReactNode;
	endIcon?: ReactNode;
	fullWidth?: boolean;
	disabled?: boolean;
	component?: ElementType;
	to?: string;
}

const ActionButton = ({
	children,
	onClick,
	variant = "contained",
	size = "large",
	startIcon,
	endIcon,
	fullWidth = false,
	disabled = false,
	component,
	to,
	className = "",
}: ActionButtonProps) => {
	const theme = useTheme();

	const getButtonStyles = () => {
		if (variant === "contained") {
			return {
				backgroundColor: theme.palette.primary.main,
				"&:hover": {
					backgroundColor: theme.palette.primary.dark,
				},
			};
		}
		return {};
	};

	const linkProps =
		component !== undefined
			? {
					component,
					...(to !== undefined ? { to } : {}),
				}
			: {};

	return (
		<Button
			variant={variant}
			size={size}
			onClick={onClick}
			startIcon={startIcon}
			endIcon={endIcon}
			fullWidth={fullWidth}
			disabled={disabled}
			className={className}
			sx={{
				minWidth: fullWidth ? "auto" : 200,
				height: size === "large" ? 60 : size === "medium" ? 48 : 36,
				fontSize:
					size === "large" ? "1.1rem" : size === "medium" ? "1rem" : "0.875rem",
				fontWeight: 600,
				textTransform: "none",
				borderRadius: 2,
				...getButtonStyles(),
			}}
			{...linkProps}
		>
			{children}
		</Button>
	);
};

export default ActionButton;
