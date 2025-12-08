import { Box, Container } from "@mui/material";
import type { BaseComponentProps } from "../../types";

interface CenteredContentProps extends BaseComponentProps {
	children: React.ReactNode;
	maxWidth?: "xs" | "sm" | "md" | "lg" | "xl";
	padding?: number;
	minHeight?: string | number;
}

const CenteredContent = ({
	children,
	maxWidth = "md",
	padding = 8,
	minHeight = "auto",
	className = "",
}: CenteredContentProps) => {
	return (
		<Box className={className}>
			<Container maxWidth={maxWidth} sx={{ py: padding, minHeight }}>
				{children}
			</Container>
		</Box>
	);
};

export default CenteredContent;
